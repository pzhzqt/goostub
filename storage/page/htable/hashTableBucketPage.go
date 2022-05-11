// Copyright (c) 2021-2022 Qitian Zeng
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package htable

import (
	"goostub/common"
	"goostub/hash"
	"reflect"
	"unsafe"
)

/**
 * Store indexed key and and value together within bucket page. Supports
 * non-unique keys.
 *
 * Bucket page starts with one byte indicating the key size in byte (255 is big enough for our need)
 * page format (keys are stored in order):
 *  ----------------------------------------------------------------
 * keySize | occupied(n bits) | readable(n bits) | KEY(1) + VALUE(1) | ... | KEY(n) + VALUE(n)
 *  ----------------------------------------------------------------
 *
 *  Here '+' means concatenation.
 */
type HashTableBucketPage struct {
	keySize uint8
	/*
	   logical layout:
	      occupied [(bucketArraySize-1)/8+1]byte
	      readable [(bucketArraySize-1)/8+1]byte
	      kvArray [bucketArraySize * size of keyVal]byte
	*/
}

/**
 * Scan the bucket and collect values that have the matching key
 *
 * @return true if at least one key matched
 */
func (h *HashTableBucketPage) GetValue(key []byte, result *[]common.RID) bool {
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
}

/**
 * Removes a key and value.
 *
 * @return true if removed, false if not found
 */
func (p *HashTableBucketPage) Remove(key []byte, value common.RID) bool {
}

/**
 * Gets the key at an index in the bucket.
 *
 * @param bucket_idx the index in the bucket to get the key at
 * @return key at index bucket_idx of the bucket
 */
func (p *HashTableBucketPage) KeyAt(bucketIdx int) []byte {
}

/**
 * Gets the value at an index in the bucket.
 *
 * @param bucket_idx the index in the bucket to get the value at
 * @return value at index bucket_idx of the bucket
 */
func (h *HashTableBucketPage) ValueAt(bucketIdx int) common.RID {
}

/**
 * Remove the KV pair at bucket_idx
 */
func (p *HashTableBucketPage) RemoveAt(bucketIdx int) {
}

// helper functions

func (p *HashTableBucketPage) kvSize() int {
	return int(p.keySize) + int(unsafe.Sizeof(common.RID{}))
}

func (p *HashTableBucketPage) bucketArraySize() int {
	//    1    +        2*((x-1)/8+1)      +   (key size + value size)*x   <=  pageSize
	// keySize      occupied + readable               kvArray
	return (4*common.PageSize - 11) / (4*p.kvSize() + 1)
}

func (p *HashTableBucketPage) bitArraySize() int {
	return int((p.bucketArraySize()-1)/8 + 1)
}

// return reference to the occupied bit array as a byte slice
// you can think of it as p.occupied
func (p *HashTableBucketPage) getOccupied(bucketIdx int) []byte {
	sh := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(p)) + 1,
		Len:  p.bitArraySize(),
		Cap:  p.bitArraySize(),
	}
	return *(*[]byte)(unsafe.Pointer(&sh))
}

// return reference to the readable bit array as a byte slice
// you can think of it as p.readable
func (p *HashTableBucketPage) getReadable() []byte {
	sh := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(p)) + 1 + uintptr(p.bitArraySize()),
		Len:  p.bitArraySize(),
		Cap:  p.bitArraySize(),
	}
	return *(*[]byte)(unsafe.Pointer(&sh))
}
