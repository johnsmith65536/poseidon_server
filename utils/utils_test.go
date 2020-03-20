package utils

import "testing"

func TestIntersection(t *testing.T) {
	alice := []int64{1,2,3}
	bob := []int64{1}
	t.Log(Intersection(alice,bob))
}
