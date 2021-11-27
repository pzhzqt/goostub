// Copyright (c) 2021 Qitian Zeng
// 
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package types

import (
	"bytes"
	"encoding/binary"
	"log"
	"strconv"
)

type TinyintType struct {
	IntegerBaseType
}

func newTinyintType() *TinyintType {
	return &TinyintType{
		*newIntegerBaseType(TINYINT),
	}
}

func (t *TinyintType) ToString(v *Value) (string, error) {
	if v.IsNull() {
		return "tinyint_null", nil
	}

	val, ok := v.val.(int8)
	if !ok {
		log.Fatalln("tinyint member function called on non-tinyint value")
	}

	return strconv.FormatInt(int64(val), 10), nil
}

func (t *TinyintType) SerializeTo(v *Value, storage *bytes.Buffer) error {
	val, ok := v.val.(int8)
	if !ok {
		log.Fatalln("bigint member function called on non-bigint value")
	}

	return binary.Write(storage, binary.LittleEndian, val)
}

func (t *TinyintType) DeserializeFrom(storage *bytes.Buffer) (*Value, error) {
	var val int8
	err := binary.Read(storage, binary.LittleEndian, &val)
	if err != nil {
		return nil, err
	}

	return NewValue(t.id, val), nil
}
