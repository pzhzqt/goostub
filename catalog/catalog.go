// Copyright (c) 2021 Qitian Zeng
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package catalog

import (
	"goostub/common"
)

type TableInfo struct {
	//TODO
}

type IndexInfo struct {
	//TODO
}

type Catalog struct {
	//TODO
}

func (c *Catalog) GetTable(tableOid common.TableOID) *TableInfo {
	//TODO
	return nil
}

func (c *Catalog) GetIndex(indexOid common.IndexOID) *IndexInfo {
	//TODO
	return nil
}
