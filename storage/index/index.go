// Copyright (c) 2021 Qitian Zeng
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package index

import (
	"goostub/buffer"
	"goostub/common"
	"goostub/schema"
	"goostub/storage/table"
)

type IndexMetadata struct {
	// strings in Go are immutable
	// so there's no need for making them private and use getter
	Name      string         // index name
	TableName string         // table name
	keyAttrs  []uint32       // mapping relation between index and table schema
	keySchema *schema.Schema // schema of index key
}

func NewIndexMetadata(name string, tableName string, keyAttrs []uint32, s *schema.Schema) *IndexMetadata {
	return &IndexMetadata{
		Name:      name,
		TableName: tableName,
		keyAttrs:  keyAttrs,
		keySchema: schema.CopySchema(s, keyAttrs),
	}
}

// These 2 getters are pointless right now, just here for future flexibility
func (im *IndexMetadata) GetKeyAttrs() []uint32 {
	return im.keyAttrs
}

func (im *IndexMetadata) GetKeySchema() *schema.Schema {
	return im.keySchema
}

func (im *IndexMetadata) GetIndexColumnCount() uint32 {
	return uint32(len(im.keyAttrs))
}

func (im IndexMetadata) String() string {
	return "IndexMetadata[Name = " + im.Name + ", Type = B+Tree, Table name = " + im.TableName + "] :: " + im.keySchema.String()
}

type Index interface {
	GetMetadata() *IndexMetadata
	GetIndexColumnCount() uint32
	GetName() string
	GetKeySchema() *schema.Schema
	GetKeyAttrs() []uint32
	String() string

	/**
	 * Insert an entry into the index.
	 * @param key The index key
	 * @param rid The RID associated with the key (unused)
	 * @param transaction The transaction context
	 */
	InsertEntry(key *table.Tuple, rid common.RID, transaction common.Transaction)

	/**
	 * Delete an index entry by key.
	 * @param key The index key
	 * @param rid The RID associated with the key (unused)
	 * @param transaction The transaction context
	 */
	DeleteEntry(key *table.Tuple, rid common.RID, transaction common.Transaction)

	/**
	 * Search the index for the provided key.
	 * @param key The index key
	 * @param result The collection of RIDs that is populated with results of the search
	 * @param transaction The transaction context
	 */
	ScanKey(key *table.Tuple, result *[]common.RID, transaction common.Transaction)
}

// base struct for all indexes
type baseIndex struct {
	metadata *IndexMetadata
}

func (bi *baseIndex) GetMetadata() *IndexMetadata {
	return bi.metadata
}

func (bi *baseIndex) GetIndexColumnCount() uint32 {
	return bi.metadata.GetIndexColumnCount()
}

func (bi *baseIndex) GetName() string {
	return bi.metadata.Name
}

func (bi *baseIndex) GetKeySchema() *schema.Schema {
	return bi.metadata.GetKeySchema()
}

func (bi *baseIndex) GetKeyAttrs() []uint32 {
	return bi.metadata.GetKeyAttrs()
}

func (bi *baseIndex) String() string {
	return "INDEX: (" + bi.metadata.String() + ")"
}

// clumsy way to implement factory function that returns different type of index
// if you define a new index type, then the index must implement the indexFactory
type indexFactory interface {
	createIndex(*IndexMetadata, *buffer.BufferPoolManager, ...any) Index
}

// IndexTypePtr: *IndexType
func NewIndex[IndexTypePtr indexFactory](m *IndexMetadata, bm *buffer.BufferPoolManager, args ...any) Index {
	var tmp IndexTypePtr
	return tmp.createIndex(m, bm, args...)
}
