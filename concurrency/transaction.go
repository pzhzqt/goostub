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

/**
 * Transaction states for 2PL:
 *
 *     _________________________
 *    |                         v
 * GROWING -> SHRINKING -> COMMITTED   ABORTED
 *    |__________|________________________^
 *
 * Transaction states for Non-2PL:
 *     __________
 *    |          v
 * GROWING  -> COMMITTED     ABORTED
 *    |_________________________^
 *
 **/
type TransactionState uint8

const (
	Growing TransactionState = iota
	Shrinking
	Committed
	Aborted
)

type IsolationLevel uint8

const (
	ReadUncommitted IsolationLevel = iota
	RepeatableRead
	ReadCommitted
)

// Type of write operation.
type WType uint8

const (
	Insert WType = iota
	Delete
	Update
)

type TableOID uint32
type IndexOID uint32

type TableWriteRecord struct {
	Rid   common.RID
	Wtype WType
	Tuple table.Tuple      // only used for update operation
	Table *table.TableHeap // specify which table the record is for
}

type IndexWriteRecord struct {
	Rid      common.RID //value stored in the index
	TableOid TableOID
	Wtype    WType
	Tuple    table.Tuple      // used to construct an index key
	OldTuple table.Tuple      // only used for update operation
	IndexOid IndexOID         // the identifier of an index into the index list
	Catalog  *catalog.Catalog // contains metadata required to locate index
}

// Reason to a transaction abortion
type AbortReason uint8

const (
	LockOnShrinking AbortReason = iota
	UnlockOnShrinking
	UpgradeConflict
	Deadlock
	LockSharedOnReadUncommitted
)

type TransactionAbortError struct {
	txnId       common.TxnID
	abortReason AbortReason
}

func (e *TransactionAbortError) GetTransactionId() common.TxnID {
	return e.txnId
}

func (e *TransactionAbortError) GetAbortReason() AbortReason {
	return e.abortReason
}

func (e *TransactionAbortError) Error() string {
	switch e.abortReason {
	case LockOnShrinking:
		return fmt.Sprintln("Transaction ", e.txnId, " aborted because it cannot take locks in the shrinking state")
	case UnlockOnShrinking:
		return fmt.Sprintln("Transaction ", e.txnId, " aborted because it cannot execute unlock in the shrinking state")
	case UpgradeConflict:
		return fmt.Sprintln("Transaction ", e.txnId, " aborted because another transaction is already waiting to upgrade its lock")
	case Deadlock:
		return fmt.Sprintln("Transaction ", e.txnId, " aborted on deadlock")
	case LockSharedOnReadUncommitted:
		return fmt.Sprintln("Transaction ", e.txnId, " aborted on lockshared on ReadUncommitted")
	}

	return "Unknown AbortReason"
}

type Transaction interface {
	// removed unused methods from bustub
	GetTransactionId() common.TxnID
	GetIsolationLevel() IsolationLevel
	GetWriteSet() *deque.Deque                    // deque of TableWriteRecord
	GetIndexWriteSet() *deque.Deque               // deque of IndexWriteRecord
	GetSharedLockSet() map[common.RID]struct{}    // resources under shared lock
	GetExclusiveLockSet() map[common.RID]struct{} // resources under ex lock
	IsSharedLocked(common.RID) bool
	IsExclusiveLocked(common.RID) bool
	GetState() TransactionState
	SetState(TransactionState)
	GetPrevLSN() common.LSN
	SetPrevLSN(common.LSN)
}

func NewTransaction(txnId common.TxnID, isolationLevel ...IsolationLevel) Transaction {
	isoL := RepeatableRead
	if len(isolationLevel) == 1 {
		isoL = isolationLevel[0]
	}
	txn := &txnInstance{
		state:            Growing,
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
	state            TransactionState
	isolationLevel   IsolationLevel
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
func (t *txnInstance) GetIsolationLevel() IsolationLevel {
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
func (t *txnInstance) GetState() TransactionState {
	return t.state
}
func (t *txnInstance) SetState(s TransactionState) {
	t.state = s
}
func (t *txnInstance) GetPrevLSN() common.LSN {
	return t.prevLSN
}
func (t *txnInstance) SetPrevLSN(lsn common.LSN) {
	t.prevLSN = lsn
}
