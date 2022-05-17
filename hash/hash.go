// Copyright (c) 2021 Qitian Zeng
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package hash

import (
	"github.com/cespare/xxhash/v2"
)

// a hash function takes a byte sequence and returns a uint64
type HashFunc func(key []byte) uint64

// the default hash function
var GoosTubHash HashFunc = xxhash.Sum64
