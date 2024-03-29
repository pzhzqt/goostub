// Copyright (c) 2021 Qitian Zeng
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package common

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
)

type ErrID int

const (
	INVALID ErrID = iota
	// Value out of range
	OUT_OF_RANGE
	// Conversion/casting error
	CONVERSION
	// Unknown type in the type subsystem
	UNKNOWN_TYPE
	// Decimal-related errors
	DECIMAL
	// Type mismatch
	MISMATCH_TYPE
	// Division by 0
	DIVIDE_BY_ZERO
	// Incompatible type
	INCOMPATIBLE_TYPE
	// Method not implemented
	NOT_IMPLEMENTED
)

// GoosTub Error struct
type GTError struct {
	id  ErrID
	msg string
}

func (e *GTError) Error() string {
	return e.msg
}

func NewError(eid ErrID, msgs ...interface{}) *GTError {
	return &GTError{
		id:  eid,
		msg: fmt.Sprint(msgs...),
	}
}

func NewErrorf(eid ErrID, str string, args ...interface{}) *GTError {
	return &GTError{
		id:  eid,
		msg: fmt.Sprintf(str, args...),
	}
}

// Check whether an error is GoosTub Error and the type is matched
func CheckErrorType(err error, id ErrID) bool {
	gterr, ok := err.(*GTError)
	if !ok {
		return false
	}
	if gterr.id != id {
		return false
	}
	return true
}

// helper struct for using assert in non-unittest
type assertHelper struct{}

func (a assertHelper) Errorf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

var Assert = assert.New(assertHelper{})
