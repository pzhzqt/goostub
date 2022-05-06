// Copyright (c) 2021 Qitian Zeng
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package common

import (
	"fmt"
	"github.com/gammazero/deque"
)

// put in common to avoid import cycle

// Struct defined in concurrency/transaction.go
type Transaction interface {
	// removed unused methods from bustub
	GetTransactionId() TxnID
	GetIsolationLevel() IsolationLevel
	GetWriteSet() *deque.Deque             // deque of TableWriteRecord
	GetIndexWriteSet() *deque.Deque        // deque of IndexWriteRecord
	GetSharedLockSet() map[RID]struct{}    // resources under shared lock
	GetExclusiveLockSet() map[RID]struct{} // resources under ex lock
	IsSharedLocked(RID) bool
	IsExclusiveLocked(RID) bool
	GetState() TransactionState
	SetState(TransactionState)
	GetPrevLSN() LSN
	SetPrevLSN(LSN)
}

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
	txnId       TxnID
	abortReason AbortReason
}

func (e *TransactionAbortError) GetTransactionId() TxnID {
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
