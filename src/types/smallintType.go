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

type SmallintType struct {
	IntegerBaseType
}

func newSmallintType() *SmallintType {
	return &SmallintType{
		*newIntegerBaseType(TINYINT),
	}
}

func (t *SmallintType) ToString(v *Value) (string, error) {
	if v.IsNull() {
		return "smallint_null", nil
	}

	val, ok := v.val.(int16)
	if !ok {
		log.Fatalln("smallint member function called on non-smallint value")
	}

	return strconv.FormatInt(int64(val), 10), nil
}

func (t *SmallintType) SerializeTo(v *Value, storage *bytes.Buffer) error {
	val, ok := v.val.(int16)
	if !ok {
		log.Fatalln("bigint member function called on non-bigint value")
	}

	return binary.Write(storage, binary.LittleEndian, val)
}

func (t *SmallintType) DeserializeFrom(storage *bytes.Buffer) (*Value, error) {
	var val int16
	err := binary.Read(storage, binary.LittleEndian, &val)
	if err != nil {
		return nil, err
	}

	return NewValue(t.id, val), nil
}
