package disk

import (
	"common"
	"os"
	"strings"
	"github.com/golang-collections/go-datastructures/futures"
)

/**
 * DiskManager takes care of the allocation and deallocation of pages within a database. It performs the reading and
 * writing of pages to and from disk, providing a logical file layer within the context of a database management system.
 */

type DiskManager struct {
    logIO       os.File
    logName     string
    dbIO        os.File
    fileName    string
    nextPageID  common.PageID
    numFlushes  int
    numWrites   int
    flushLog    bool
    flushLogF  *futures.Future
}

func (d *DiskManager) getFileSize(fileName *string) {
}

/**
* Creates a new disk manager that writes to the specified database file.
* @param dbFile the file name of the database file to write to
*/
func NewDiskManager(dbFile string) *DiskManager {
    d := &DiskManager{
        fileName : dbFile,
        nextPageID : 0,
        numFlushes : 0,
        numWrites : 0,
        flushLog : false,
        flushLogF : nil,
    }
    n := strings.LastIndex(dbFile, ".")
    if n == -1 {
        LOG_DEBUG
}

/**
* Shut down the disk manager and close all the file resources.
*/
func (d *DiskManager) ShutDown() {
}

/**
* Write a page to the database file.
* @param pageID id of the page
* @param pageData raw page data
*/
func (d *DiskManager) WritePage(pageID common.PageID, pageData []byte) {
}

/**
* Read a page from the database file.
* @param pageID id of the page
* @param[out] pageData output buffer
*/
func (d *DiskManager) ReadPage(pageID common.PageID, pageData []byte) {
}

/**
* Flush the entire log buffer into disk.
* @param logData raw log data
* @param size size of log entry
*/
func (d *DiskManager) WriteLog(logData []byte, size int) {
}

/**
* Read a log entry from the log file.
* @param[out] logData output buffer
* @param size size of the log entry
* @param offset offset of the log entry in the file
* @return true if the read was successful, false otherwise
*/
func (d *DiskManager) ReadLog(logData []byte, size int, offset int) bool {
}

/**
* Allocate a page on disk.
* @return the id of the allocated page
*/
func (d *DiskManager) AllocatePage() common.PageID {
}

/**
* Deallocate a page on disk.
* @param pageID id of the page to deallocate
*/
func (d *DiskManager) DeallocatePage(pageID common.PageID) {
}

/** @return the number of disk flushes */
func (d *DiskManager) GetNumFlushes() int {
}

/** @return true iff the in-memory content has not been flushed yet */
func (d *DiskManager) GetFlushState() bool {
}

/** @return the number of disk writes */
func (d *DiskManager) GetNumWrites() int {
}

/**
* Sets the future which is used to check for non-blocking flushes.
* @param f the non-blocking flush check
*/
func (d *DiskManager) SetFlushLogFuture(f *futures.Future) {
    d.flushLogF = f
}

/** Checks if the non-blocking flush future was set. */
func (d *DiskManager) HasFlushLogFuture() bool {
    return d.flushLogF != nil
}
