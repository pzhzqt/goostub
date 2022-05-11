// Copyright (c) 2021 Qitian Zeng
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package transaction

import (
	"goostub/common"
	"goostub/concurrency"
	"goostub/recovery"
	"sync"
	"sync/atomic"
)

var TxnMap = make(map[common.TxnID]common.Transaction)
var TxnMapMutex = sync.RWMutex{}

type TransactionManager struct {
	nextTxnId      common.TxnID //operation needs to be atomic on this
	lockManager    *concurrency.LockManager
	logManager     *recovery.LogManager
	globalTxnLatch sync.RWMutex
}

func (tm *TransactionManager) Begin(txn common.Transaction, isoLevel ...common.IsolationLevel) common.Transaction {
	tm.globalTxnLatch.RLock()
	if txn == nil {
		txn = NewTransaction(tm.nextTxnId, isoLevel...)
		atomic.AddInt32((*int32)(&tm.nextTxnId), 1)
	}
	TxnMapMutex.Lock()
	TxnMap[txn.GetTransactionId()] = txn
	TxnMapMutex.Unlock()
	return txn
}

func (tm *TransactionManager) Commit(txn common.Transaction) {
	txn.SetState(common.Committed)
	writeSet := txn.GetWriteSet()

	// Perform all deletes before we commit.
	for writeSet.Len() > 0 {
		item := writeSet.Back().(TableWriteRecord)
		table := item.Table
		if item.Wtype == common.Delete {
			// Note that this also releases the lock when holding the page latch.
			table.ApplyDelete(item.Rid, txn)
		}
		writeSet.PopBack()
	}
	writeSet.Clear()

	// Release all the locks
	tm.releaseLocks(txn)
	tm.globalTxnLatch.RUnlock()
}

func (tm *TransactionManager) Abort(txn common.Transaction) {
	txn.SetState(common.Aborted)

	// Rollback before releasing the lock
	tableWriteSet := txn.GetWriteSet()
	for tableWriteSet.Len() > 0 {
		item := tableWriteSet.Back().(TableWriteRecord)
		table := item.Table
		if item.Wtype == common.Delete {
			table.RollbackDelete(item.Rid, txn)
		} else if item.Wtype == common.Insert {
			// Note that this also releases the lock when holding the page latch.
			table.ApplyDelete(item.Rid, txn)
		} else if item.Wtype == common.Update {
			table.UpdateTuple(&item.Tuple, item.Rid, txn)
		}
		tableWriteSet.PopBack()
	}
	tableWriteSet.Clear()

	// Rollback index updates
	indexWriteSet := txn.GetIndexWriteSet()
	for indexWriteSet.Len() > 0 {
		item := indexWriteSet.Back().(IndexWriteRecord)
		catalog := item.Catalog
		// Metadata identifying the table that should be deleted from.
		tableInfo := catalog.GetTable(item.TableOid)
		indexInfo := catalog.GetIndex(item.IndexOid)
		// TODO: left here
	}

}

func (tm *TransactionManager) releaseLocks(txn common.Transaction) {
	lockSet := make(map[common.RID]struct{})
	for item := range txn.GetExclusiveLockSet() {
		lockSet[item] = struct{}{}
	}
	for item := range txn.GetSharedLockSet() {
		lockSet[item] = struct{}{}
	}
	for item := range lockSet {
		tm.lockManager.Unlock(txn, item)
	}
}
