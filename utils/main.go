package utils

func Clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

