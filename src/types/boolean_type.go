package types

import (
    "log"
    "common"
    "github.com/go-kit/kit/log/level"
    "unsafe"
)

type BooleanType struct {
    BaseType
}

func newBooleanType() Type {
    return &BooleanType {
        BaseType {
            id : BOOLEAN,
        },
    }
}

func (t *BooleanType) Compare(l *Value, r *Value) CmpResult {
    if t.GetTypeID() != BOOLEAN {
        log.Fatalln("BooleanType member function called from non-boolean type")
    }

    if !l.CheckComparable(r) {
        level.Error(common.Logger).Log("Not comparable")
        return nil
    }

    if l.IsNull() || r.IsNull() {
        return nil
    }

    lval := l.CastAs(BOOLEAN).val.(int8)
    rval := r.CastAs(BOOLEAN).val.(int8)

    var ret int

    if lval < rval {
        ret = -1
    } else if lval > rval {
        ret = 1
    } else {
        ret = 0
    }

    return &ret
}

func (t *BooleanType) IsInlined(v *Value) int8 {
    return 1
}

func (t *BooleanType) ToString(v *Value) string {
    if t.GetTypeID() != BOOLEAN {
        log.Fatalln("BooleanType member function called from non-boolean type")
    }

    if v.val.(int8) == 1 {
        return "true"
    } else if v.val.(int8) == 0 {
        return "false"
    }

    return "boolean_null"
}

func (t *BooleanType) SerializeTo(v *Value, storage *byte) {
    *(*int8)(unsafe.Pointer(storage)) = v.val.(int8)
}

func (t *BooleanType) DeserializeFrom(storage *byte) *Value {
    val := *(*int8)(unsafe.Pointer(storage))
    return NewValue(BOOLEAN, val)
}

func (t *BooleanType) Copy(v *Value) *Value {
    return NewValue(BOOLEAN, v.val.(int8))
}

func (t *BooleanType) CastAs(v *Value, id TypeID) *Value {
    switch (id) {
    case BOOLEAN:
        return t.Copy(v)
    case VARCHAR:
        if (v.IsNull()) {
            return NewValue(VARCHAR, nil, false)
        }
        return NewValue(VARCHAR, v.ToString())
    default:
        break
    }
    level.Error(common.Logger).Log("Boolean is not coercable to ",TypeIDToString(id))
    return nil
}
