package maths

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
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

func AbsUint(x uint) uint {
	if x < 0 {
		x = -x
	}
	return x
}
