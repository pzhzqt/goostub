// Copyright (c) 2021 Qitian Zeng
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package schema

import (
	"fmt"
	"strings"
)

type Schema struct {
	// size of fixed-length columns, i.e. size of one tuple
	length uint32

	// all columns, inlined and uninlined
	columns []Column

	// true if all columns are inlined
	tupleIsInlined bool

	// indices of all uninlined columns
	uninlinedColumns []int
}

func NewSchema(columns []Column) *Schema {
	schema := &Schema{
		tupleIsInlined: true,
	}

	curOffset := uint32(0)
	for idx := 0; idx < len(columns); idx++ {
		column := columns[idx]
		// handle uninlined column
		if !column.IsInlined() {
			schema.tupleIsInlined = false
			schema.uninlinedColumns = append(schema.uninlinedColumns, idx)
		}

		// set offset
		column.columnOffset = curOffset
		curOffset += column.GetFixedLength()

		// add column
		schema.columns = append(schema.columns, column)
	}

	schema.length = curOffset
	return schema
}

func CopySchema(from *Schema, attrs []uint32) *Schema {
	cols := make([]Column, len(attrs))
	for i := 0; i < len(attrs); i++ {
		cols[i] = from.columns[attrs[i]]
	}
	return NewSchema(cols)
}

func (s *Schema) GetColumns() []Column {
	return s.columns
}

func (s *Schema) GetColumn(colIdx int) *Column {
	return &s.columns[colIdx]
}

// -1 indicates error
func (s *Schema) GetColIdx(colName string) int {
	for i := 0; i < len(s.columns); i++ {
		if s.columns[i].GetName() == colName {
			return i
		}
	}

	return -1
}

func (s *Schema) GetUninlinedColumns() []int {
	return s.uninlinedColumns
}

func (s *Schema) GetColumnCount() int {
	return len(s.columns)
}

func (s *Schema) GetUninlinedColumnCount() int {
	return len(s.uninlinedColumns)
}

func (s *Schema) GetLength() uint32 {
	return s.length
}

func (s *Schema) IsInlined() bool {
	return s.tupleIsInlined
}

func (s Schema) String() string {
	b := &strings.Builder{}
	b.WriteString(
		fmt.Sprint(
			"Schema[", "NumColumns:", s.GetColumnCount(), ", ",
			"IsInlined:", s.tupleIsInlined, ", ",
			"Length:", s.length, "]"),
	)

	b.WriteString(" :: (")

	for i := 0; i < s.GetColumnCount(); i++ {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(s.columns[i].String())
	}
	b.WriteString(")")
	return b.String()
}
