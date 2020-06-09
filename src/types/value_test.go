package types

import (
    "testing"
    "fmt"
)

func TestInt8(t *testing.T) {
    val := NewValue(BOOLEAN, 1)
    if val == nil {
        t.Error("Nil value")
    }
    fmt.Println("id: ", val.typeID)
    fmt.Println("size: ", val.size)
    fmt.Println("manage_data: ", val.manage_data)
    fmt.Println("val: ", val.val)
}
