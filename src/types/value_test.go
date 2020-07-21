package types

import (
	"fmt"
	"testing"
)

func TestInt8(t *testing.T) {
	val := NewValue(BOOLEAN, 1)
	if val == nil {
		t.Error("Nil value")
	}
	fmt.Println("id: ", val.typeID)
	fmt.Println("size: ", val.isNull)
	fmt.Println("manage_data: ", val.manageData)
	fmt.Println("val: ", val.val)
}
