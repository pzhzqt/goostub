package common

import (
	"fmt"
)

type RID struct {
	pageID  PageID
	slotNum uint32
}

func DefaultRID() *RID {
	return &RID{
		pageID:  InvalidPageID,
		slotNum: 0,
	}
}

func NewRIDFromInt64(rid int64) *RID {
	return &RID{
		pageID:  PageID(rid >> 32),
		slotNum: uint32(rid),
	}
}

func NewRID(pid PageID, snum uint32) *RID {
	return &RID{
		pageID:  pid,
		slotNum: snum,
	}
}

func (r *RID) Get() int64 {
	return int64(r.slotNum) | int64(r.pageID)<<32
}

func (r *RID) GetPageId() PageID {
	return r.pageID
}

func (r *RID) GetSlotNum() uint32 {
	return r.slotNum
}

func (r *RID) Set(pid PageID, snum uint32) {
	r.pageID = pid
	r.slotNum = snum
}

func (r RID) String() string {
	return fmt.Sprintln("PageID: ", r.pageID, " SlotNum: ", r.slotNum)
}

func (r *RID) IsEqual(other *RID) bool {
	return r.pageID == other.pageID && r.slotNum == other.slotNum
}
