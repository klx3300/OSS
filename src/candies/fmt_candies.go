package candies

import "strconv"

// Ftoa is the faster path to format float. only 64 acceptable
func Ftoa(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

// Atof is the faster path too
func Atof(a string) (float64, error) {
	return strconv.ParseFloat(a, 64)
}
