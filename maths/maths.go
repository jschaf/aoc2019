package maths

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func ModWithSameSign(d, m int) int {
	var res = d % m
	if (res < 0 && m > 0) || (res > 0 && m < 0) {
		return res + m
	}
	return res
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func AbsInt(x int) int {
	if x < 0 {
		x = -x
	}
	return x
}

func MaxUint(a, b uint) uint {
	if a > b {
		return a
	}
	return b
}

func MinUint(a, b uint) uint {
	if a < b {
		return a
	}
	return b
}
