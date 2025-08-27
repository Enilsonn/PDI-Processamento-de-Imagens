package utils

func Max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func Limiar(x, min, max int) uint8 {
	if x < min {
		return uint8(min)
	} else if x > max {
		return uint8(max)
	}
	return uint8(x)
}

func Abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
