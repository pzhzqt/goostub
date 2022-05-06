// Copyright (c) 2021 Qitian Zeng
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package table

import (
	"goostub/common"
)

type TableHeap struct {
	// TODO
}

func (t *TableHeap) ApplyDelete(rid common.RID, txn common.Transaction) {
	// TODO
}

func (t *TableHeap) RollbackDelete(rid common.RID, txn common.Transaction) {
	// TODO
}

func (t *TableHeap) UpdateTuple(tup *Tuple, rid common.RID, txn common.Transaction) {
	// TODO
}
