package comparator

func InSlice[T comparable](val T, sliceVal []T) bool {
	for _, v := range sliceVal {
		if v == val {
			return true
		}
	}

	return false
}

func IdxSlice[T comparable](val T, sliceVal []T) int {
	for idx, v := range sliceVal {
		if v == val {
			return idx
		}
	}

	return -1
}
