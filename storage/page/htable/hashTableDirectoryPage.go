// Copyright (c) 2021 Qitian Zeng
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package htable

import (
	"fmt"
	"github.com/go-kit/kit/log/level"
	"goostub/common"
	"os"
)

const directoryArraySize = 512

/**
 *
 * Directory Page for extendible hash table.
 *
 * Directory format (size in byte):
 * --------------------------------------------------------------------------------------------
 * | PageId (4) | LSN(4) | GlobalDepth(4) | LocalDepths(512) | BucketPageIds(2048) | Free(1524)
 * --------------------------------------------------------------------------------------------
 */
type HashTableDirectoryPage struct {
	pageId        common.PageID
	lsn           common.LSN
	globalDepth   uint32
	localDepth    [directoryArraySize]uint8
	bucketPageIds [directoryArraySize]common.PageID
}

func (p *HashTableDirectoryPage) GetPageId() common.PageID {
	return p.pageId
}

func (p *HashTableDirectoryPage) SetPageId(pid common.PageID) {
	p.pageId = pid
}

func (p *HashTableDirectoryPage) GetLSN() common.LSN {
	return p.lsn
}

func (p *HashTableDirectoryPage) SetLSN(lsn common.LSN) {
	p.lsn = lsn
}

/**
 * Lookup a bucket page using a directory index
 *
 * @param bucketIdx the index in the directory to lookup
 * @return bucket page id corresponding to bucketIdx
 */
func (p *HashTableDirectoryPage) GetBucketPageId(bucketIdx uint32) common.PageID {
	return 0
}

/**
 * Updates the directory index using a bucket index and page_id
 *
 * @param bucket_idx directory index at which to insert page_id
 * @param bucket_page_id page_id to insert
 */
func (p *HashTableDirectoryPage) SetBucketPageId(bucketIdx uint32, bucketPageId common.PageID) {
}

/**
 * Gets the split image of an index
 *
 * @param bucket_idx the directory index for which to find the split image
 * @return the directory index of the split image
 **/
func (p *HashTableDirectoryPage) GetSplitImageIndex(bucketIdx uint32) {
}

/**
 * GetGlobalDepthMask - returns a mask of global_depth 1's and the rest 0's.
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
 * @return mask of global_depth 1's and the rest 0's (with 1's from LSB upwards)
 */
func (p *HashTableDirectoryPage) GetGlobalDepthMask() uint32 {
	return 0
}

/**
 * GetLocalDepthMask - same as global depth mask, except it
 * uses the local depth of the bucket located at bucket_idx
 *
 * @param bucket_idx the index to use for looking up local depth
 * @return mask of local 1's and the rest 0's (with 1's from LSB upwards)
 */
func (p *HashTableDirectoryPage) GetLocalDepthMask(bucketIdx uint32) uint32 {
	return 0
}

/**
 * Get the global depth of the hash table directory
 *
 * @return the global depth of the directory
 */
func (p *HashTableDirectoryPage) GetGlobalDepth() uint32 {
	return 0
}

/**
 * Increase the global depth of the directory
 */
func (p *HashTableDirectoryPage) IncrGlobalDepth() {}

/**
 * Decrease the global depth of the directory
 */
func (p *HashTableDirectoryPage) DecrGlobalDepth() {}

/**
 * @return true if the directory can be shrunk
 */
func (p *HashTableDirectoryPage) CanShrink() bool {
	return false
}

/**
 * @return the current directory size
 */
func (p *HashTableDirectoryPage) Size() uint32 {
	return 0
}

/**
 * Gets the local depth of the bucket at bucket_idx
 *
 * @param bucket_idx the bucket index to lookup
 * @return the local depth of the bucket at bucket_idx
 */
func (p *HashTableDirectoryPage) GetLocalDepth(bucketIdx uint32) uint8 {
	return 0
}

/**
 * Set the local depth of the bucket at bucket_idx to local_depth
 *
 * @param bucket_idx bucket index to update
 * @param local_depth new local depth
 */
func (p *HashTableDirectoryPage) SetLocalDepth(bucketIdx uint32, localDepth uint8) {}

/**
 * Increase the local depth of the bucket at bucket_idx
 * @param bucket_idx bucket index to increment
 */
func (p *HashTableDirectoryPage) IncrLocalDepth(bucketIdx uint32) {}

/**
 * Decrease the local depth of the bucket at bucket_idx
 * @param bucket_idx bucket index to decrement
 */
func (p *HashTableDirectoryPage) DecrLocalDepth(bucketIdx uint32) {}

/**
 * Gets the high bit corresponding to the bucket's local depth.
 * This is not the same as the bucket index itself.  This method
 * is helpful for finding the pair, or "split image", of a bucket.
 *
 * @param bucket_idx bucket index to lookup
 * @return the high bit corresponding to the bucket's local depth
 */
func (p *HashTableDirectoryPage) GetLocalHighBit(bucketIdx uint32) uint32 {
	return 0
}

/**
 * VerifyIntegrity - Use this for debugging but **DO NOT CHANGE**
 *
 * If you want to make changes to this, make a new function and extend it.
 *
 * Verify the following invariants:
 * (1) All LD <= GD.
 * (2) Each bucket has precisely 2^(GD - LD) pointers pointing to it.
 * (3) The LD is the same at each index with the same bucket_page_id
 */
func (p *HashTableDirectoryPage) VerifyIntegrity() {
	//  build maps of {bucket_page_id : pointer_count} and {bucket_page_id : local_depth}
	pageId2Count := make(map[common.PageID]uint32)
	pageId2Ld := make(map[common.PageID]uint8)

	//  verify for each bucket_page_id
	for curIdx := uint32(0); curIdx < directoryArraySize; curIdx++ {
		curPageId := p.bucketPageIds[curIdx]
		curLd := p.localDepth[curIdx]
		common.Assert.LessOrEqual(curLd, p.globalDepth)
		pageId2Count[curPageId]++

		if oldLd, ok := pageId2Ld[curPageId]; ok && oldLd != curLd {
			level.Warn(common.Logger).Log(fmt.Sprintf("Verify Integrity: cur local depth: %d, old local depth %d, for page id: %d", curLd, oldLd, curPageId))
			p.PrintDirectory()
			os.Exit(1)
		} else {
			pageId2Ld[curPageId] = curLd
		}
	}

	for curPageId, curCount := range pageId2Count {
		curLd := uint32(pageId2Ld[curPageId])
		requiredCount := uint32(1 << (p.globalDepth - curLd))
		if curCount != requiredCount {
			level.Warn(common.Logger).Log(fmt.Sprintf("Verify Integrity: cur count: %d, required count: %d, for page id: %d", curCount, requiredCount, curPageId))
			p.PrintDirectory()
			os.Exit(1)
		}
	}
}

/**
 * Prints the current directory
 */
func (p *HashTableDirectoryPage) PrintDirectory() {
	level.Debug(common.Logger).Log(fmt.Sprintf("======== DIRECTORY (global depth: %d) ========", p.globalDepth))
	level.Debug(common.Logger).Log("| bucket idx | page id | local depth |")
	for idx := 0; idx < (1 << p.globalDepth); idx++ {
		level.Debug(common.Logger).Log(fmt.Sprintf("|     %d     |     %d     |     %d     |", idx, p.bucketPageIds[idx], p.localDepth[idx]))
	}
	level.Debug(common.Logger).Log("================ END DIRECTORY ================")
}
