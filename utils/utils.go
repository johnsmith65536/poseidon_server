package utils

func Intersection(alice, bob []int64) map[int64]bool {
	aliceMap, ret := make(map[int64]bool), make(map[int64]bool)
	for _, value := range alice {
		aliceMap[value] = true
	}
	for _, value := range bob {
		if aliceMap[value] {
			ret[value] = true
		}
	}
	return ret
}
