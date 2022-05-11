// Copyright (c) 2021 Qitian Zeng
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package schema

import (
	"fmt"
	"goostub/execution/expressions"
	"goostub/types"
	"strings"
)

type Column struct {
	columnName     string
	columnType     types.TypeID
	fixedLength    uint32 // for non-inlined column this is size of int
	variableLength uint32 // 0 for inlined column
	columnOffset   uint32 // column offset in the tuple

	expr *expressions.AbstractExpression // expression used to create column
}

func NewColumn(colName string, id types.TypeID, expr *expressions.AbstractExpression) *Column {
	return &Column{
		columnName:  colName,
		columnType:  id,
		fixedLength: typeSize(id),
		expr:        expr,
	}
}

// Different from BusTub: there's not a separate constructor for varchar type where variablelength is specified. The column struct is a part of schema which is supposed to be metadata and can exist without data. Therefore, variableLength should always start from 0 and increase as inserting tuples to the table and thus a separate constructor is redundant

func (c *Column) GetName() string {
	return c.columnName
}

func (c *Column) GetLength() uint32 {
	if c.IsInlined() {
		return c.fixedLength
	}
	return c.variableLength
}

func (c *Column) GetFixedLength() uint32 {
	return c.fixedLength
}

func (c *Column) GetVariableLength() uint32 {
	return c.variableLength
}

func (c *Column) GetOffset() uint32 {
	return c.columnOffset
}

func (c *Column) IsInlined() bool {
	return c.columnType != types.VARCHAR
}

func (c Column) String() string {
	b := &strings.Builder{}
	b.WriteString("Column[" + c.columnName + ", " + types.TypeIDToString(c.columnType) + ", Offset:" + fmt.Sprint(c.columnOffset) + ", ")
	if c.IsInlined() {
		b.WriteString("FixedLength:" + fmt.Sprint(c.fixedLength))
	} else {
		b.WriteString("VarLength:" + fmt.Sprint(c.variableLength))
	}
	b.WriteString("]")
	return b.String()
}

func (c *Column) GetExpr() *expressions.AbstractExpression {
	return c.expr
}

func typeSize(id types.TypeID) uint32 {
	switch id {
	case types.BOOLEAN:
		return 1
	case types.TINYINT:
		return 1
	case types.SMALLINT:
		return 2
	case types.INTEGER:
		return 4
	case types.BIGINT:
	case types.DECIMAL:
	case types.TIMESTAMP:
		return 8
	case types.VARCHAR:
		return 4
	}
	return 0
}
