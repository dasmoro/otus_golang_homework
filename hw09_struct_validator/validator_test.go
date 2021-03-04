package hw09StructValidator //nolint:golint,stylecheck

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	ResponseWithSecret struct {
		Code   int    `validate:"in:200,404,500"`
		Body   string `json:"omitempty"`
		secret string `validate:"in:secret,s"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				"qwertyqwertyqwertyqwertyqwertyqwerty",
				"User",
				18,
				"qwe@qwe.qw",
				"admin",
				[]string{"+1234567890", "+1231231231"},
				nil,
			},
			expectedErr: nil,
		},
		{
			in: Response{
				200,
				`{"json":"json"}`,
			},
			expectedErr: nil,
		},
		{
			in: User{
				"not valid id",
				"User",
				18,
				"qwe@qwe.qw",
				"admin",
				[]string{"+1234567890", "+1231231231"},
				nil,
			},
			expectedErr: errors.New("ID has error: field must contain a 36 characters\n"),
		},
		{
			in: User{
				"not valid id",
				"User",
				51,
				"qweqweqwe.qw",
				"admin",
				[]string{"+1234567890", "+1231231231"},
				nil,
			},
			expectedErr: errors.New("ID has error: field must contain a 36 characters\nEmail has error: field is not valid for pattern ^\\w+@\\w+\\.\\w+$\nAge has error: field must be less or equal than 50\n"),
		},
		{
			in: ResponseWithSecret{
				404,
				`{"text":"<html></html>}`,
				"ss",
			},
			expectedErr: nil,
		},
		{
			in: Token{
				[]byte{1, 2, 3, 4, 5},
				[]byte{1, 2, 3, 4, 5},
				[]byte{1, 2, 3, 4, 5},
			},
			expectedErr: nil,
		},
		// ...
		// Place your code here
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			// Place your code here
			vErr := Validate(tt.in)
			if vErr != nil {
				errorsReceived := strings.Split(vErr.Error(), "\n")
				errorsExpected := strings.Split(tt.expectedErr.Error(), "\n")
				if len(errorsExpected) != len(errorsReceived) {
					t.Error(vErr)
				} else {
					for _, err := range errorsExpected {
						if !strings.Contains(vErr.Error(), err) {
							t.Error(vErr)
							break
						}
					}
				}
			} else {
				if vErr != tt.expectedErr {
					t.Error(vErr)
				}
			}
		})
	}
}
