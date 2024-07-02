package utils

import (
	"fmt"
	"log"
	"reflect"
)

// AssertEqual checks if two values are deeply equal.
// If not, it logs a fatal error with details of the values.
func AssertEqual(expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		log.Fatalf("Assertion failed: Expected %+v, but got %+v", expected, actual)
	}
}
