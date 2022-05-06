// Copyright (c) 2021 Qitian Zeng
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package concurrency

import (
	"fmt"
	"github.com/gammazero/deque"
	"goostub/catalog"
	"goostub/common"
	"goostub/storage/table"
)

type TableWriteRecord struct {
	Rid   common.RID
	Wtype common.WType
	Tuple table.Tuple      // only used for update operation
	Table *table.TableHeap // specify which table the record is for
}

type IndexWriteRecord struct {
	Rid      common.RID //value stored in the index
	TableOid common.TableOID
	Wtype    common.WType
	Tuple    table.Tuple      // used to construct an index key
	OldTuple table.Tuple      // only used for update operation
	IndexOid common.IndexOID  // the identifier of an index into the index list
	Catalog  *catalog.Catalog // contains metadata required to locate index
}

func NewTransaction(txnId common.TxnID, isolationLevel ...common.IsolationLevel) common.Transaction {
	isoL := common.RepeatableRead
	if len(isolationLevel) == 1 {
		isoL = isolationLevel[0]
	}
	txn := &txnInstance{
		state:            common.Growing,
		isolationLevel:   isoL,
		txnId:            txnId,
		prevLSN:          common.InvalidLSN,
		sharedLockSet:    map[common.RID]struct{}{},
		exclusiveLockSet: map[common.RID]struct{}{},
		tableWriteSet:    &deque.Deque{},
		indexWriteSet:    &deque.Deque{},
	}
	return txn
}

type txnInstance struct {
	state            common.TransactionState
	isolationLevel   common.IsolationLevel
	txnId            common.TxnID
	tableWriteSet    *deque.Deque // deque of TableWriteRecord
	indexWriteSet    *deque.Deque // deque of IndexWriteRecord
	prevLSN          common.LSN
	sharedLockSet    map[common.RID]struct{}
	exclusiveLockSet map[common.RID]struct{}
}

func (t *txnInstance) GetTransactionId() common.TxnID {
	return t.txnId
}
func (t *txnInstance) GetIsolationLevel() common.IsolationLevel {
	return t.isolationLevel
}
func (t *txnInstance) GetWriteSet() *deque.Deque {
	return t.tableWriteSet
}
func (t *txnInstance) GetIndexWriteSet() *deque.Deque {
	return t.indexWriteSet
}
func (t *txnInstance) GetSharedLockSet() map[common.RID]struct{} {
	return t.sharedLockSet
}
func (t *txnInstance) GetExclusiveLockSet() map[common.RID]struct{} {
	return t.exclusiveLockSet
}
func (t *txnInstance) IsSharedLocked(rid common.RID) bool {
	_, ok := t.sharedLockSet[rid]
	return ok
}
func (t *txnInstance) IsExclusiveLocked(rid common.RID) bool {
	_, ok := t.exclusiveLockSet[rid]
	return ok
}
func (t *txnInstance) GetState() common.TransactionState {
	return t.state
}
func (t *txnInstance) SetState(s common.TransactionState) {
	t.state = s
}
func (t *txnInstance) GetPrevLSN() common.LSN {
	return t.prevLSN
}
func (t *txnInstance) SetPrevLSN(lsn common.LSN) {
	t.prevLSN = lsn
}
