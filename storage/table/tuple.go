// Copyright (c) 2021 Qitian Zeng
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package table

import (
	"bytes"
	"encoding/binary"
	"github.com/go-kit/kit/log/level"
	"goostub/common"
	"goostub/schema"
	"goostub/types"
	"log"
	"reflect"
)

type Tuple struct {
	allocated bool       // is allocated?
	rid       common.RID // if pointing to the table heap, the rid is valid
	data      []byte
}

// valid args: NewTuple(), NewTuple(rid), NewTuple(tuple), NewTuple(values, schema)
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
	case 2:
		if vals, ok := args[0].([]*types.Value); ok {
			if schema, ok := args[1].(*schema.Schema); ok {
				return newTupleFromValues(vals, schema)
			}
		}
	}
	log.Fatalln("Invalid argument for tuple constructor")
	return nil
}

func CopyTuple(other *Tuple) *Tuple {
	newTuple := &Tuple{
		allocated: other.allocated,
		rid:       other.rid,
	}

	if newTuple.allocated {
		newTuple.data = make([]byte, len(other.data))
		if copy(newTuple.data, other.data) != len(other.data) {
			level.Error(common.Logger).Log("Data copy error")
			return nil
		}
	} else {
		newTuple.data = other.data
	}

	return newTuple
}

func defaultTuple() *Tuple {
	return &Tuple{
		rid: common.DefaultRID(),
	}
}

func newTupleFromRID(rid common.RID) *Tuple {
	return &Tuple{
		rid: rid,
	}
}

func newTupleFromValues(vals []*types.Value, schema *schema.Schema) *Tuple {
	if len(vals) != schema.GetColumnCount() {
		log.Fatalln("Value and Schema not matched")
	}

	t := &Tuple{
		allocated: true,
	}

	tupleSize := schema.GetLength()

	for i := range schema.GetUninlinedColumns() {
		// uninlined are varchar columns, need extra bytes to indicate length
		tupleSize += vals[i].GetLength() + uint32(reflect.TypeOf(tupleSize).Size())
	}

	t.data = make([]byte, tupleSize)

	buf := bytes.NewBuffer(t.data[:0])
	// varchar data offset
	offset := schema.GetLength()
	varcharbuf := bytes.NewBuffer(t.data[offset:offset])

	for i := 0; i < schema.GetColumnCount(); i++ {
		col := schema.GetColumn(i)
		if !col.IsInlined() {
			// write data offset to current column position
			binary.Write(buf, binary.LittleEndian, offset)
			// write actual data to offset
			vals[i].SerializeTo(varcharbuf)
			// advance offset for the next varchar column
			offset += vals[i].GetLength() + uint32(reflect.TypeOf(tupleSize).Size())
		} else {
			// just write to current column position
			vals[i].SerializeTo(buf)
		}
	}

	return t
}

func (t *Tuple) SerializeTo(storage *bytes.Buffer) {
	binary.Write(storage, binary.LittleEndian, uint32(len(t.data)))
	binary.Write(storage, binary.LittleEndian, t.data)
}

func (t *Tuple) DeserializeFrom(storage *bytes.Buffer) {
	var size uint32
	binary.Read(storage, binary.LittleEndian, &size)
	t.data = make([]byte, size)
	binary.Read(storage, binary.LittleEndian, t.data)
	t.allocated = true
}

func (t *Tuple) GetData() []byte {
	return t.data
}
