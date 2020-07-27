package table

import (
	"catalog"
	"common"
	"types"
)

type Tuple struct {
	allocated bool       // is allocated?
	rid       common.RID // if pointing to the table heap, the rid is valid
	size      uint32
	data      []byte
}

// TODO: left here
