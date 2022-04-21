// Copyright (c) 2021 Qitian Zeng
// 
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package page

import (
	"goostub/common"
)

const (
	deleteMask = (1 << 31)
)

/**
 * Slotted page format:
 *  ---------------------------------------------------------
 *  | HEADER | ... FREE SPACE ... | ... INSERTED TUPLES ... |
 *  ---------------------------------------------------------
 *                                ^
 *                                free space pointer
 *
 *  Header format (size in bytes):
 *  ----------------------------------------------------------------------------
 *  | PageId (4)| LSN (4)| PrevPageId (4)| NextPageId (4)| FreeSpacePointer(4) |
 *  ----------------------------------------------------------------------------
 *  ----------------------------------------------------------------
 *  | TupleCount (4) | Tuple_1 offset (4) | Tuple_1 size (4) | ... |
 *  ----------------------------------------------------------------
 *
 */

type TablePage interface {
	Page
	Init(common.PageID, uint32, common.PageID)
}

/**
* Initialize the TablePage header.
* @param page_id the page ID of this table page
* @param page_size the size of this table page
* @param prev_page_id the previous table page ID
* @param log_manager the log manager in use
* @param txn the transaction that this page is created in
 */
