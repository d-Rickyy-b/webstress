package util

import "math"

type Number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

func NumOfChars[N Number](n N) N {
	if n < 0 {
		return NumOfChars(-n) + 1
	} else if n == 0 {
		return 1
	}
	return (N)(math.Log10(float64(n)) + 1)
}
