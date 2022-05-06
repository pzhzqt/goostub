// Copyright (c) 2021 Qitian Zeng
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT
package concurrency

import (
	"container/list"
	"goostub/common"
	"sync"
)

type LockMode uint8

const (
	Shared LockMode = iota
	Exclusive
)

type LockRequest struct {
	Mode    LockMode
	TxnId   common.TxnID
	Granted bool
}

type LockRequestQueue struct {
	queue *list.List
	// for notifying blocked transactions on this rid
	cv *sync.Cond
	// txnId of an upgrading transaction (InvalidTxnID = none)
	upgrading common.TxnID
}

type LockManager struct {
	latch     sync.Mutex
	lockTable map[common.RID]LockRequestQueue
}

/*
 * [LOCK_NOTE]: For all locking functions, we:
 * 1. return false if the transaction is aborted; and
 * 2. block on wait, return true when the lock request is granted; and
 * 3. it is undefined behavior to try locking an already locked RID in the
 * same transaction, i.e. the transaction is responsible for keeping track of
 * its current locks.
 */

/**
 * Acquire a lock on RID in shared mode. See [LOCK_NOTE].
 * @param txn the transaction requesting the shared lock
 * @param rid the RID to be locked in shared mode
 * @return true if the lock is granted, false otherwise
 */
func (l *LockManager) LockShared(txn common.Transaction, rid common.RID) bool {
	txn.GetSharedLockSet()[rid] = struct{}{}
	return true
}

/**
 * Acquire a lock on RID in exclusive mode. See [LOCK_NOTE] in header file.
 * @param txn the transaction requesting the exclusive lock
 * @param rid the RID to be locked in exclusive mode
 * @return true if the lock is granted, false otherwise
 */
func (l *LockManager) LockExclusive(txn common.Transaction, rid common.RID) bool {
	txn.GetExclusiveLockSet()[rid] = struct{}{}
	return true
}

/**
 * Upgrade a lock from a shared lock to an exclusive lock.
 * @param txn the transaction requesting the lock upgrade
 * @param rid the RID that should already be locked in shared mode by the
 * requesting transaction
 * @return true if the upgrade is successful, false otherwise
 */
func (l *LockManager) LockUpgrade(txn common.Transaction, rid common.RID) bool {
	delete(txn.GetSharedLockSet(), rid)
	txn.GetExclusiveLockSet()[rid] = struct{}{}
	return true
}

/**
 * Release the lock held by the transaction.
 * @param txn the transaction releasing the lock, it should actually hold the
 * lock
 * @param rid the RID that is locked by the transaction
 * @return true if the unlock is successful, false otherwise
 */
func (l *LockManager) Unlock(txn common.Transaction, rid common.RID) bool {
	delete(txn.GetSharedLockSet(), rid)
	delete(txn.GetExclusiveLockSet(), rid)
	return true
}
