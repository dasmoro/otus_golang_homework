package hw09_struct_validator //nolint:golint,stylecheck

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	result := strings.Builder{}
	for _, vErr := range v {
		result.WriteString(fmt.Sprintf("%v has error: %v\n", vErr.Field, vErr.Err))
	}
	return result.String()
}

type Validator struct {
	fieldT reflect.StructField
	fieldV reflect.Value
}

func (v *Validator) isInt() *ValidationError {
	if v.fieldV.Kind() != reflect.Int {
		vErr := &ValidationError{v.fieldT.Name, errors.New("field must be an integer")}
		return vErr
	}
	return nil
}

func (v *Validator) isStr() *ValidationError {
	if v.fieldV.Kind() != reflect.String {
		vErr := &ValidationError{v.fieldT.Name, fmt.Errorf("field must be a string")}
		return vErr
	}
	return nil
}

func (v *Validator) required(rexp regexp.Regexp, tag string) *ValidationError {
	if rexp.MatchString(tag) && v.fieldV.IsZero() {
		vErr := &ValidationError{v.fieldT.Name, errors.New("field is required")}
		return vErr
	}
	return nil
}

func (v *Validator) min(rexp regexp.Regexp, tag string) *ValidationError {
	if !rexp.MatchString(tag) {
		return nil
	}
	err := v.isInt()
	if err != nil {
		return err
	}
	minS := rexp.FindStringSubmatch(tag)[1]
	if min, _ := strconv.Atoi(minS); int(v.fieldV.Int()) < min {
		vErr := &ValidationError{v.fieldT.Name, fmt.Errorf("field must be greater or equal than %v", min)}
		return vErr
	}
	return nil
}

func (v *Validator) max(rexp regexp.Regexp, tag string) *ValidationError {
	if !rexp.MatchString(tag) {
		return nil
	}
	err := v.isInt()
	if err != nil {
		return err
	}
	maxS := rexp.FindStringSubmatch(tag)[1]
	if max, _ := strconv.Atoi(maxS); int(v.fieldV.Int()) > max {
		vErr := &ValidationError{v.fieldT.Name, fmt.Errorf("field must be less or equal than %v", max)}
		return vErr
	}
	return nil
}

func (v *Validator) strSize(rexp regexp.Regexp, tag string) *ValidationError {
	if !rexp.MatchString(tag) {
		return nil
	}
	err := v.isStr()
	if err != nil {
		return err
	}
	sizeS := rexp.FindStringSubmatch(tag)[1]
	strV := v.fieldV.String()
	if size, _ := strconv.Atoi(sizeS); len(strV) != size {
		vErr := &ValidationError{v.fieldT.Name, fmt.Errorf("field must contain a %v characters", size)}
		return vErr
	}
	return nil
}

func (v *Validator) inSlice(rexp regexp.Regexp, tag string) *ValidationError {
	if !rexp.MatchString(tag) {
		return nil
	}
	inV := rexp.FindStringSubmatch(tag)[1]
	items := strings.Split(inV, ",")
	var vErr *ValidationError
	if v.fieldV.Kind() == reflect.String {
		for _, item := range items {
			if item == v.fieldV.String() {
				return nil
			}
		}
		vErr = &ValidationError{v.fieldT.Name, fmt.Errorf("field must be in set of %v", items)}
	}
	if v.fieldV.Kind() == reflect.Int {
		for _, item := range items {
			if item == strconv.Itoa(int(v.fieldV.Int())) {
				return nil
			}
		}
		vErr = &ValidationError{v.fieldT.Name, fmt.Errorf("field must be in set of %v", items)}
	}
	if vErr == nil {
		vErr = &ValidationError{v.fieldT.Name, fmt.Errorf("field must be integer or string")}
	}
	return vErr
}

func (v *Validator) rexp(rexp regexp.Regexp, tag string) *ValidationError {
	if !rexp.MatchString(tag) {
		return nil
	}

	typeErr := v.isStr()
	if typeErr != nil {
		return typeErr
	}
	pattern := rexp.FindStringSubmatch(tag)[1]

	pattern = strings.ReplaceAll(pattern, `\\`, `\`)
	rexpp, err := regexp.Compile(pattern)
	if err != nil {
		return &ValidationError{v.fieldT.Name, fmt.Errorf("%v is not valid pattern", pattern)}
	}
	if !rexpp.MatchString(v.fieldV.String()) {
		return &ValidationError{v.fieldT.Name, fmt.Errorf("field is not valid for pattern %v", pattern)}
	}
	return nil
}

func (v *Validator) Run() ValidationErrors {
	fieldTag := string(v.fieldT.Tag)
	if !strings.Contains(fieldTag, `validate:"`) {
		return nil
	}
	vRexp := regexp.MustCompile(`validate:"(.+)"`)
	tag := vRexp.FindStringSubmatch(fieldTag)[1]

	tags := strings.Split(tag, "|")
	rexps := map[string]*regexp.Regexp{
		"required": regexp.MustCompile(`^required`),
		"len":      regexp.MustCompile(`^len:(\d+)`),
		"in":       regexp.MustCompile(`^in:([^|$]+)`),
		"min":      regexp.MustCompile(`^min:(\d+)`),
		"max":      regexp.MustCompile(`^max:(\d+)`),
		"regexp":   regexp.MustCompile(`^regexp:(.+)`),
	}
	valFuncs := map[string]func(rexp regexp.Regexp, tag string) *ValidationError{
		"required": v.required,
		"max":      v.max,
		"min":      v.min,
		"in":       v.inSlice,
		"regexp":   v.rexp,
		"len":      v.strSize,
	}
	vErrs := make(ValidationErrors, 0)
	for _, tag := range tags {
		for key, valFunc := range valFuncs {
			err := valFunc(*rexps[key], tag)
			if err != nil {
				vErrs = append(vErrs, *err)
			}
		}
	}
	return vErrs
}

func checkField(fieldT reflect.StructField, fieldV reflect.Value) ValidationErrors {
	if fieldV.Kind() == reflect.Struct {
		err := Validate(Validate(fieldV.Interface()))
		if err != nil {
			return []ValidationError{{fieldT.Name, err}}
		}
		return nil
	}
	validator := Validator{fieldT, fieldV}
	return validator.Run()
}

var validateMtx = sync.Mutex{}

func validateField(refT reflect.Type, refV reflect.Value, i int, errBuf io.StringWriter) {
	fieldT := refT.Field(i)
	fieldV := refV.Field(i)
	if !fieldV.CanInterface() {
		return
	}
	if fieldV.Kind() == reflect.Slice {
		for j := 0; j < fieldV.Len(); j++ {
			errs := checkField(fieldT, fieldV.Index(j))
			if len(errs) > 0 {
				validateMtx.Lock()
				_, _ = errBuf.WriteString(errs.Error())
				validateMtx.Unlock()
			}
		}
	} else {
		errs := checkField(fieldT, fieldV)
		if len(errs) > 0 {
			validateMtx.Lock()
			_, _ = errBuf.WriteString(errs.Error())
			validateMtx.Unlock()
		}
	}
}

func Validate(v interface{}) error {
	errBuf := strings.Builder{}
	refV := reflect.ValueOf(v)
	refT := reflect.TypeOf(v)
	wg := sync.WaitGroup{}
	if refV.Kind() == reflect.Struct {
		for i := 0; i < refT.NumField(); i++ {
			wg.Add(1)
			go func(i int) {
				validateField(refT, refV, i, &errBuf)
				wg.Done()
			}(i)
		}
		wg.Wait()
	}
	if errBuf.Len() > 0 {
		return errors.New(errBuf.String())
	}
	return nil
}
