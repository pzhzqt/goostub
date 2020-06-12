package common

import (
    "fmt"
    "testing"
)

func TestError(t *testing.T) {
    a := NewErrorf(OUT_OF_RANGE, "value %d out of range", 5)
    if !CheckErrorType(a, OUT_OF_RANGE) {
        t.Error("type not matched")
    }
    if a.Error() != fmt.Sprintf("value %d out of range", 5) {
        t.Error("msg not matched")
    }
}
