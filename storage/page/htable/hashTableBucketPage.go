// Copyright (c) 2021-2022 Qitian Zeng
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package htable

import (
	"goostub/common"
	"goostub/hash"
	"goostub/storage/page"
	"unsafe"
)

/**
* Store indexed key and and value together within bucket page. Supports
* non-unique keys.
*
* page format (keys are stored in order):
*  ----------------------------------------------------------------
*  occupied(n bits) | readable(n bits) | KEY(1) + VALUE(1) | ... | KEY(n) + VALUE(n)
*  ----------------------------------------------------------------
*
*  Here '+' means concatenation.
 */

// apart from keySize, fields are references to the actual page in buffer pool
type HashTableBucketPage struct {
	keySize  uint32
	occupied []byte
	readable []byte
	kvArray  []byte
}

// get a bucket page pointer to existing page
func PageAsBucketPage(page page.Page, keySize uint32) *HashTableBucketPage {
	d := page.GetData()
	p := &HashTableBucketPage{
		keySize: keySize,
	}
	p.occupied = d[:p.bitArraySize()]
	p.readable = d[p.bitArraySize() : 2*p.bitArraySize()]
	p.kvArray = d[2*p.bitArraySize():]
	return p
}

/**
 * Scan the bucket and collect values that have the matching key
 *
 * @return true if at least one key matched
 */
func (h *HashTableBucketPage) GetValue(key []byte, result *[]common.RID) bool {
	return false
}

/**
 * Attempts to insert a key and value in the bucket.  Uses the occupied_
 * and readable_ arrays to keep track of each slot's availability.
 *
 * @param key key to insert
 * @param value value to insert
 * @return true if inserted, false if duplicate KV pair or bucket is full
 */
func (p *HashTableBucketPage) Insert(key []byte, value common.RID) bool {
	return false
}

/**
 * Removes a key and value.
 *
 * @return true if removed, false if not found
 */
func (p *HashTableBucketPage) Remove(key []byte, value common.RID) bool {
	return false
}

/**
 * Remove the KV pair at bucket_idx
 */
func (p *HashTableBucketPage) RemoveAt(bucketIdx uint32) {
}

// helper functions

func (p *HashTableBucketPage) kvSize() uint32 {
	return p.keySize + uint32(unsafe.Sizeof(common.RID{}))
}

func (p *HashTableBucketPage) bucketArraySize() uint32 {
	//     2*((x-1)/8+1)      +   (kv size)*x   <=  pageSize
	// occupied + readable         kvArray
	return (4*common.PageSize - 7) / (4*p.kvSize() + 1)
}

func (p *HashTableBucketPage) bitArraySize() uint32 {
	return (p.bucketArraySize()-1)/8 + 1
}

// get a reference to the key at bucketIdx in this page
func (p *HashTableBucketPage) keyAt(bucketIdx uint32) []byte {
	return p.kvArray[bucketIdx*p.kvSize() : bucketIdx*p.kvSize()+p.keySize]
}

// get the value (RID) at bucketIdx in this page, use a pointer so it's mutable
func (p *HashTableBucketPage) valueAt(bucketIdx uint32) *common.RID {
	return (*common.RID)(unsafe.Pointer(&p.kvArray[bucketIdx*p.kvSize()+p.keySize]))
}
