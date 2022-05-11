// Copyright (c) 2021 Qitian Zeng
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package table

import (
	"goostub/buffer"
	"goostub/common"
	"goostub/concurrency"
	"goostub/recovery"
)

type TableHeap struct {
	// TODO
}

/**
* Create a table heap with a transaction. (create table)
* @param buffer_pool_manager the buffer pool manager
* @param lock_manager the lock manager
* @param log_manager the log manager
* @param txn the creating transaction
 */
func NewTableHeap(bpm *buffer.BufferPoolManager, lockM *concurrency.LockManager, logM *recovery.LogManager, txn common.Transaction) *TableHeap {
	return nil
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
