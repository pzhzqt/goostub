// Copyright (c) 2021 Qitian Zeng
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package types

import (
	"github.com/go-kit/kit/log/level"
	"github.com/stretchr/testify/assert"
	"goostub/common"
	"testing"
)

var testTypes = [...]TypeID{BOOLEAN, TINYINT, SMALLINT, INTEGER, BIGINT, DECIMAL}

func init() {
	common.InitLogger(level.AllowNone())
}

func TestInvalidType(t *testing.T) {
	typeId := INVALID
	ty := GetInstance(typeId)
	a := assert.New(t)
	a.Equal(typeId, ty.GetTypeID())
	a.False(ty.IsCoercableFrom(typeId))
	_, err := GetTypeSize(typeId)
	a.Error(err)
	a.Nil(GetMinValue(typeId))
	a.Nil(GetMaxValue(typeId))
}

func TestGetInstance(t *testing.T) {
	a := assert.New(t)
	for _, tid := range testTypes {
		ty := GetInstance(tid)
		a.NotNil(ty)
		a.Equal(tid, ty.GetTypeID())
		a.True(ty.IsCoercableFrom(tid))
	}
}

func TestMaxMinValue(t *testing.T) {
	a := assert.New(t)
	for _, tid := range testTypes {
		a.False(GetMaxValue(tid).IsNull())
		a.False(GetMinValue(tid).IsNull())
	}
}

func TestStringIntCmp(t *testing.T) {
	v1 := NewValue(VARCHAR, "32")
	v2 := NewValue(INTEGER, 32)
	res, err := v1.CompareTo(v2)
	a := assert.New(t)
	a.Nil(err)
	a.Equal(CmpEqual, res)
}
