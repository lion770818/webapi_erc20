package utils

import (
	"strconv"
	"strings"
)

func IsEmpty(str string) bool {
	str = strings.TrimSpace(str)
	return str == ""
}

func IsNotEmpty(str string) bool {
	return !IsEmpty(str)
}

func IsEmptySlice(str []string) bool {
	return len(str) == 0
}

func IsNotEmptySlice(str []string) bool {
	return !IsEmptySlice(str)
}

func Atoi(str string, defaultValue int) (int, error) {
	if IsEmpty(str) {
		return defaultValue, nil
	}

	result, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue, err
	}
	return result, nil
}

func AtoiAndDefault(str string) (int, error) {
	return Atoi(str, DefaultNotAssignInt)
}
