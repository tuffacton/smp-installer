package util

func Contains[V comparable](data []V, toSearch V) bool {
	for _, v := range data {
		if v == toSearch {
			return true
		}
	}
	return false
}
