package coursebus

import (
	"fmt"
	"reflect"
)

// Helper function to generate a unique key for a struct
func structKey(v interface{}) string {
	return fmt.Sprintf("%#v", v)
}

// Generic function to compare slices of any struct type, ignoring order
func slicesEqualUnordered[T any](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	countA := make(map[string]int)
	countB := make(map[string]int)

	for _, item := range a {
		countA[structKey(item)]++
	}
	for _, item := range b {
		countB[structKey(item)]++
	}

	return reflect.DeepEqual(countA, countB)
}

// Helper function to check if a slice is nil or empty
func isNilOrEmpty[T any](s []T) bool {
	return len(s) == 0
}
