// Copyright (c) 2021 Qitian Zeng
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package catalog

import (
	"goostub/buffer"
	"goostub/common"
	"goostub/concurrency"
	"goostub/hash"
	"goostub/recovery"
	"goostub/schema"
	"goostub/storage/index"
	"goostub/storage/table"
	"sync/atomic"
)

type TableInfo struct {
	Schema schema.Schema
	Name   string
	Table  *table.TableHeap
	Oid    common.TableOID
}

type IndexInfo struct {
	KeySchema schema.Schema // schema for the index key
	Name      string
	Index     *index.Index
	IndexOid  common.IndexOID
	TableName string
	KeySize   uintptr // size of the index key in bytes
}

/**
 * The Catalog is a non-persistent catalog that is designed for
 * use by executors within the DBMS execution engine. It handles
 * table creation, table lookup, index creation, and index lookup.
 */
type Catalog struct {
	bpm          *buffer.BufferPoolManager
	lockManager  *concurrency.LockManager
	logManager   *recovery.LogManager
	tables       map[common.TableOID]*TableInfo
	tableNames   map[string]common.TableOID
	nextTableOid common.TableOID //need atomic operation
	indexes      map[common.IndexOID]*IndexInfo
	//table name->index name->index oid
	indexNames   map[string]map[string]common.IndexOID
	nextIndexOid common.IndexOID //need atomic operation
}

func NewCatalog(bpm *buffer.BufferPoolManager, lockManager *concurrency.LockManager, logManager *recovery.LogManager) *Catalog {
	return &Catalog{
		bpm:          bpm,
		lockManager:  lockManager,
		logManager:   logManager,
		tables:       make(map[common.TableOID]*TableInfo),
		tableNames:   make(map[string]common.TableOID),
		nextTableOid: 0,
		indexes:      make(map[common.IndexOID]*IndexInfo),
		indexNames:   make(map[string]map[string]common.IndexOID),
		nextIndexOid: 0,
	}
}

func (c *Catalog) CreateTable(txn common.Transaction, name string, schema *schema.Schema) *TableInfo {
	if _, ok := c.tableNames[name]; ok {
		// table already exists
		return nil
	}
	// construct the table heap
	table := table.NewTableHeap(c.bpm, c.lockManager, c.logManager, txn)

	// awkward way to do atomic_fetch_add in go
	tableOid := common.TableOID(atomic.AddUint32((*uint32)(&c.nextTableOid), 1) - 1)

	// construct the table info
	meta := &TableInfo{
		Schema: *schema,
		Name:   name,
		Table:  table,
		Oid:    tableOid,
	}

	// update the internal tracking mechanisms
	c.tables[tableOid] = meta
	c.tableNames[name] = tableOid
	c.indexNames[name] = make(map[string]common.IndexOID)

	return meta
}

/**
 * Query table metadata by name.
 * @param table_name The name of the table
 * @return A pointer to the metadata for the table
 */
func (c *Catalog) GetTableByName(name string) *TableInfo {
	if oid, ok := c.tableNames[name]; ok {
		return c.GetTableByOid(oid)
	}
	return nil
}

/**
 * Query table metadata by OID
 * @param table_oid The OID of the table to query
 * @return A pointer to the metadata for the table
 */
func (c *Catalog) GetTableByOid(tableOid common.TableOID) *TableInfo {
	if table, ok := c.tables[tableOid]; ok {
		return table
	}
	return nil
}

/**
 * Create a new index, populate existing data of the table and return its metadata.
 * @param txn The transaction in which the table is being created
 * @param index_name The name of the new index
 * @param table_name The name of the table
 * @param schema The schema of the table
 * @param key_schema The schema of the key
 * @param key_attrs Key attributes
 * @param keysize Size of the key
 * @param (optional)hash_function The hash function for the index, none = default
 * @return A (non-owning) pointer to the metadata of the new table
 */
func (c *Catalog) CreateIndex(txn common.Transaction, indexName string, tableName string, schema *schema.Schema, keySchema *schema.Schema, keyAttrs []string, keysize uintptr, hashFunc ...hash.HashFunc) *IndexInfo {
	// TODO
}

func (c *Catalog) GetIndex(indexOid common.IndexOID) *IndexInfo {
	//TODO
	return nil
}
