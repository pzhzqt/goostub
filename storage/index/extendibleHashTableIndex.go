// Copyright (c) 2021 Qitian Zeng
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package index

import (
	"goostub/buffer"
	"goostub/common"
	"goostub/hash"
	"goostub/storage/page/htable"
	"goostub/storage/table"
)

type ExtendibleHashTableIndex struct {
	baseIndex
	container extendibleHashTable
}

func (i *ExtendibleHashTableIndex) InsertEntry(key *table.Tuple, rid common.RID, transaction common.Transaction) {
	i.container.insert(transaction, key.GetData(), rid)
}
func (i *ExtendibleHashTableIndex) DeleteEntry(key *table.Tuple, rid common.RID, transaction common.Transaction) {
	i.container.remove(transaction, key.GetData(), rid)
}
func (i *ExtendibleHashTableIndex) ScanKey(key *table.Tuple, result *[]common.RID, transaction common.Transaction) {
	i.container.getValue(transaction, key.GetData(), result)
}

func (i *ExtendibleHashTableIndex) createIndex(m *IndexMetadata, bm *buffer.BufferPoolManager, args ...any) *ExtendibleHashTableIndex {
	return &ExtendibleHashTableIndex{
		// implement me
	}
}

/**
Some simplifying changes comparing to bustub

1. 2 fields are removed: comparator, hash func
For comparator, hash index doesn't need an ordered comparator (unless you use an ordered data structure to resolve collision which is not the case here), and therefore we can simply compare whether the byte sequence is the same
Allowing user-defined hash function would be too much trouble for an educational db. we just stick to one specific hash function here.

2. removed template arguments: KeyType, ValueType, ComparatorType
ComparatorType is not needed because of the change 1
ValueType is not needed because we only need RID as value
KeyType is not needed because it doesn't make sense. key type should be decided at runtime by the key attributes instead of at compile time as a template argument. Besides, we can always treat key as a byte sequence then the only variable is the length of the sequence, and however, length is not a type and generics in Go doesn't support non-type argument

Adding another field for convenience: keySize which indicates the size of the key in bytes
*/

/**
 * Implementation of extendible hash table that is backed by a buffer pool
 * manager. Non-unique keys are supported. Supports insert and delete. The
 * table grows/shrinks dynamically as buckets become full/empty.
 */
type extendibleHashTable struct {
	directoryPageId common.PageID
	bufferManager   *buffer.BufferPoolManager
	tableLatch      common.ReaderWriterLatch
	keySize         uint8
}

/**
 * Inserts a key-value pair into the hash table.
 *
 * @param transaction the current transaction
 * @param key the key to create
 * @param value the value to be associated with the key
 * @return true if insert succeeded, false otherwise
 */
func (t *extendibleHashTable) insert(transaction common.Transaction, key []byte, value common.RID) bool {
	return false
}

/**
 * Deletes the associated value for the given key.
 *
 * @param transaction the current transaction
 * @param key the key to delete
 * @param value the value to delete
 * @return true if remove succeeded, false otherwise
 */
func (t *extendibleHashTable) remove(transaction common.Transaction, key []byte, value common.RID) bool {
	return false
}

/**
 * Performs a point query on the hash table.
 *
 * @param transaction the current transaction
 * @param key the key to look up
 * @param[out] result the value(s) associated with a given key
 * @return true if lookup succeeded, false otherwise
 */
func (t *extendibleHashTable) getValue(transaction common.Transaction, key []byte, result *[]common.RID) bool {
	return false
}

/**
 * Returns the global depth.  Do not touch.
 */
func (t *extendibleHashTable) getGlobalDepth() uint32 {
	t.tableLatch.RLock()
	defer t.tableLatch.RUnlock()
	defer common.Assert.True(t.bufferManager.UnpinPage(t.directoryPageId, false, nil))
	return t.fetchDirectoryPage().GetGlobalDepth()
}

/**
 * Helper function to verify the integrity of the extendible hash table's directory.
 * Do not touch.
 */
func (t *extendibleHashTable) verifyIntegrity() {
	t.tableLatch.RLock()
	defer t.tableLatch.RUnlock()
	t.fetchDirectoryPage().VerifyIntegrity()
	common.Assert.True(t.bufferManager.UnpinPage(t.directoryPageId, false, nil))
}

/**
 * Hash - simple helper to downcast 64-bit hash to 32-bit
 * for extendible hashing.
 *
 * @param key the key to hash
 * @return the downcasted 32-bit hash
 */
func (t *extendibleHashTable) hash(key []byte) uint32 {
	return uint32(hash.GoosTubHash(key))
}

/**
 * KeyToDirectoryIndex - maps a key to a directory index
 *
 * In Extendible Hashing we map a key to a directory index
 * using the following hash + mask function.
 *
 * DirectoryIndex = Hash(key) & GLOBAL_DEPTH_MASK
 *
 * where GLOBAL_DEPTH_MASK is a mask with exactly GLOBAL_DEPTH 1's from LSB
 * upwards.  For example, global depth 3 corresponds to 0x00000007 in a 32-bit
 * representation.
 *
 * @param key the key to use for lookup
 * @param dir_page to use for lookup of global depth
 * @return the directory index
 */
func (t *extendibleHashTable) keyToDirectoryIndex(key []byte, dirPage *htable.HashTableDirectoryPage) common.PageID {

}

/**
 * Get the bucket page_id corresponding to a key.
 *
 * @param key the key for lookup
 * @param dir_page a pointer to the hash table's directory page
 * @return the bucket page_id corresponding to the input key
 */
func (t *extendibleHashTable) keyToPageId(key []byte, dirPage *htable.HashTableDirectoryPage) common.PageID {
}

/**
 * Fetches the directory page from the buffer pool manager.
 *
 * @return a pointer to the directory page
 */
func (t *extendibleHashTable) fetchDirectoryPage() *htable.HashTableDirectoryPage {
}

/**
 * Fetches the a bucket page from the buffer pool manager using the bucket's page_id.
 *
 * @param bucket_page_id the page_id to fetch
 * @return a pointer to a bucket page
 */
func (t *extendibleHashTable) fetchBucketPage(bucketPageId common.PageID) *htable.HashTableBucketPage {
}

/**
 * Performs insertion with an optional bucket splitting.
 *
 * @param transaction a pointer to the current transaction
 * @param key the key to insert
 * @param value the value to insert
 * @return whether or not the insertion was successful
 */
func (t *extendibleHashTable) splitInsert(transaction common.Transaction, key []byte, value common.RID) bool {
}

/**
 * Optionally merges an empty bucket into it's pair.  This is called by Remove,
 * if Remove makes a bucket empty.
 *
 * There are three conditions under which we skip the merge:
 * 1. The bucket is no longer empty.
 * 2. The bucket has local depth 0.
 * 3. The bucket's local depth doesn't match its split image's local depth.
 *
 * @param transaction a pointer to the current transaction
 * @param key the key that was removed
 * @param value the value that was removed
 */
func (t *extendibleHashTable) merge(transaction common.Transaction, key []byte, value common.RID) {
}
