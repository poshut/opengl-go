package main

func clampf64(x, min, max float64) float64 {
	if x > max {
		return max
	} else if x < min {
		return min
	}
	return x
}
