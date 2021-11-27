// Copyright (c) 2021 Qitian Zeng
// 
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package table

import (
	"goostub/common"
	"goostub/metadata"
	"goostub/types"
	"log"
)

type Tuple struct {
	allocated bool       // is allocated?
	rid       common.RID // if pointing to the table heap, the rid is valid
	size      uint32
	data      []byte
}

func NewTuple(args ...interface{}) *Tuple {
	switch len(args) {
	case 0:
		return defaultTuple()
	case 1:
		if rid, ok := args[0].(common.RID); ok {
			return newTupleFromRID(rid)
		}
		if t, ok := args[0].(*Tuple); ok {
			return CopyTuple(t)
		}
		break
	case 2:
		if vals, ok := args[0].([]*types.Value); ok {
			if schema, ok := args[1].(*metadata.Schema); ok {
				return newTupleFromValues(vals, schema)
			}
		}
	}
	log.Fatalln("Invalid argument for tuple constructor")
}

func CopyTuple(other *Tuple) *Tuple {
	// TODO: implement this
	return nil
}

func defaultTuple() *Tuple {
	return &Tuple{
		rid: *common.DefaultRID(),
	}
}

func newTupleFromRID(rid common.RID) *Tuple {
	return &Tuple{
		rid: rid,
	}
}

func newTupleFromValues(vals []*types.Value, schema *metadata.Schema) *Tuple {
	// TODO: implement this
	return nil
}
