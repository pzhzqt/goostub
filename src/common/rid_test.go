// Copyright (c) 2021 Qitian Zeng
// 
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package common

import (
	"testing"
)

func TestRID(t *testing.T) {
	rid := &RID{
		pageID:  10,
		slotNum: 14,
	}

	a := NewRID(rid.pageID, rid.slotNum)
	if a.GetPageId() != rid.pageID || a.GetSlotNum() != rid.slotNum {
		t.Error("Expected: ", rid.ToString(), "Actual: ", a.ToString())
	}

	var val int64 = 10<<32 | 14
	b := NewRIDFromInt64(val)
	if b.GetPageId() != PageID(val>>32) || b.GetSlotNum() != uint32(val&0xffffffff) {
		t.Error("Expected: PageID: ", val>>32, " SlotNum: ", val&0xffffffff, "\nActual: ", b.ToString())
	}

	if b.Get() != val {
		t.Error("Expected: ", val, " Actual: ", b.Get())
	}
}
