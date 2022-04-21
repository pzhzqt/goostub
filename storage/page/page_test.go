// Copyright (c) 2021 Qitian Zeng
// 
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package page

import (
	"fmt"
	"testing"
)

func TestLSN(t *testing.T) {
	page := NewPage()

	page.SetLSN(238521)
	if page.GetLSN() != 238521 {
		t.Error("LSN not correctly set\n")
	}

	data := page.GetData()

	for i := 0; i < 8; i++ {
		fmt.Println(data[i])
	}
}
