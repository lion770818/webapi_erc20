package utils

func IsZero(i int) bool {
	return i == 0
}

func IsNotZero(i int) bool {
	return !IsZero(i)
}

func IsNegativeOne(i int) bool {
	return i == -1
}

func IsNotNegativeOne(i int) bool {
	return !IsNegativeOne(i)
}
