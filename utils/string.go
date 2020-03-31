package utils

import (
	"regexp"
)

func Substr(str string, start int, length int) string {
	var (
		strLen  int
		strRune []rune
	)

	if length == 0 {
		return ""
	}

	strRune = []rune(str)
	strLen = len(strRune)

	if start < 0 {
		start = strLen + start
	}
	if start > strLen {
		start = strLen
	}
	end := start + length
	if end > strLen {
		end = strLen
	}
	if length < 0 {
		end = strLen + length
	}
	if start > end {
		start, end = end, start
	}
	return string(strRune[start:end])
}

func SimpleSafeReplace(str string) string {
	var (
		re  *regexp.Regexp
		err error
	)

	if re, err = regexp.Compile(`(?i:script|onerror|expression|onmousemove|onload|onclick|onmouseover)`); err == nil {
		str = re.ReplaceAllString(str, " ")
	}

	return str
}

func CheckRuneRepeatNum(str string, rate int) bool {
	var (
		strRunes []rune
		counter  map[rune]int
	)

	counter = make(map[rune]int)
	strRunes = []rune(str)
	for _, runeTmp := range strRunes {
		if currentCount, exists := counter[runeTmp]; exists {
			if currentCount > rate {
				return true
			} else {
				counter[runeTmp]++
			}
		} else {
			counter[runeTmp] = 1
		}
	}

	return false
}
