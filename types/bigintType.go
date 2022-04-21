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

type BigintType struct {
	IntegerBaseType
}

func newBigintType() *BigintType {
	return &BigintType{
		*newIntegerBaseType(BIGINT),
	}
}

func (t *BigintType) ToString(v *Value) (string, error) {
	if v.IsNull() {
		return "bigint_null", nil
	}

	val, ok := v.val.(int64)
	if !ok {
		log.Fatalln("bigint member function called on non-bigint value")
	}

	return strconv.FormatInt(int64(val), 10), nil
}

func (t *BigintType) SerializeTo(v *Value, storage *bytes.Buffer) error {
	val, ok := v.val.(int64)
	if !ok {
		log.Fatalln("bigint member function called on non-bigint value")
	}

	return binary.Write(storage, binary.LittleEndian, val)
}

func (t *BigintType) DeserializeFrom(storage *bytes.Buffer) (*Value, error) {
	var val int64
	err := binary.Read(storage, binary.LittleEndian, &val)
	if err != nil {
		return nil, err
	}

	return NewValue(t.id, val), nil
}
