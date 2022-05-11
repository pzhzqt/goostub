// Copyright (c) 2021 Qitian Zeng
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package hash

import (
	"github.com/cespare/xxhash/v2"
	"goostub/common"
	"unsafe"
)

// a hash function takes a byte sequence and returns a uint64
type HashFunc func(key []byte) uint64

// the default hash function
var GoosTubHash HashFunc = xxhash.Sum64

// A reader for a concatenated key value pair
type KeyValReader struct {
	keySize uint8  // the size of the key in bytes
	kvRef   []byte // reference to the key value pair
}

func NewKeyValReader(keySize uint8, kvRef []byte) *KeyValReader {
	return &KeyValReader{
		keySize: keySize,
		kvRef:   kvRef,
	}
}

func (kvr *KeyValReader) GetKey() []byte {
	return kvr.kvRef[:kvr.keySize]
}

func (kvr *KeyValReader) GetVal() common.RID {
	return *(*common.RID)(unsafe.Pointer(&kvr.kvRef[kvr.keySize]))
}
