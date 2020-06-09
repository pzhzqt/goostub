package types

import (
    "log"
)



type Value struct {
    // use interface{} for union type
    val         interface{}
    size        interface{}
    manage_data bool
    typeID      TypeID
}

func NewValue(id TypeID, data ...interface{}) *Value {
    if len(data) == 0 {
        return newNullValue(id)
    } else if len(data) == 1 {
        switch data[0].(type) {
        case int8:
            return newValueFromInt8(id, data[0].(int8))
        }
    }

    log.Fatalln("Wrong input format")
    return nil
}

// specific value constructors

func newNullValue(id TypeID) *Value {
    return &Value {
        typeID :      id,
        size :        GOOSTUB_VALUE_NULL,
        manage_data : false,
    }
}

func newValueFromInt8(id TypeID, i int8) *Value {
    value := newNullValue(id)
    switch (id) {
    case BOOLEAN:
        value.val = i
        if value.val == GOOSTUB_BOOLEAN_NULL {
            value.size = GOOSTUB_VALUE_NULL
        } else {
            value.size = 0
        }
        break
    }
    return value
}
