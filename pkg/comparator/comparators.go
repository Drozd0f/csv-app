package comparator

func InSlice[T comparable](val T, sliceVal []T) bool {
	for _, v := range sliceVal {
		if v == val {
			return true
		}
	}

	return false
}
