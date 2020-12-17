package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"regexp"
	"sort"
	"strings"
)

type wordCounter struct {
	word  string
	count int
}

func Top10(s string) []string {
	s = strings.ToLower(s)

	freq := make(map[string]int)
	regexCompound := regexp.MustCompile(`[А-Яа-яA-Za-z]+(?:['-]+[А-Яа-яA-Za-z]+)+`)
	compounds := regexCompound.FindAllString(s, -1)
	s = regexCompound.ReplaceAllString(s, "")

	rexp := regexp.MustCompile(`[[:punct:]]`)
	s = rexp.ReplaceAllString(s, " ")

	sSlice := strings.Fields(s)
	for _, val := range sSlice {
		freq[val]++
	}
	for _, val := range compounds {
		freq[val]++
	}

	wcList := make([]wordCounter, len(freq))
	i := 0
	for key, val := range freq {
		wcList[i] = wordCounter{key, val}
		i++
	}

	sort.Slice(wcList, func(i, j int) bool {
		return wcList[i].count > wcList[j].count
	})

	result := make([]string, 0)
	for i, val := range wcList {
		result = append(result, val.word)
		if i > 10 {
			break
		}
	}
	return result
}
