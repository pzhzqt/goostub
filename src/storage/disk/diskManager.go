package disk

import (
	"common"
    "log"
	"os"
	"strings"
	"github.com/golang-collections/go-datastructures/futures"
    "github.com/go-kit/kit/log/level"
)

/**
 * DiskManager takes care of the allocation and deallocation of pages within a database. It performs the reading and
 * writing of pages to and from disk, providing a logical file layer within the context of a database management system.
 */

var bufferUsed []byte

type DiskManager struct {
    logIO       *os.File
    logName     string
    dbIO        *os.File
    fileName    string
    nextPageID  common.PageID
    numFlushes  int
    numWrites   int
    flushLog    bool
    flushLogF  *futures.Future
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
        level.Debug(common.Logger).Log("wrong file format")
        return nil
    }
    d.logName = dbFile[:n] + ".log"

    var err error
    d.logIO, err = os.OpenFile(d.logName, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        log.Fatalln("can't open dblog file")
    }

    d.dbIO, err = os.OpenFile(dbFile, os.O_RDWR | os.O_CREATE, 0666)
    if err != nil {
        log.Fatalln("can't open db file")
    }

    bufferUsed = nil

    return d
}

/**
* Shut down the disk manager and close all the file resources.
*/
func (d *DiskManager) ShutDown() {
    d.dbIO.Close()
    d.logIO.Close()
}

/**
* Write a page to the database file.
* @param pageID id of the page
* @param pageData raw page data
*/
func (d *DiskManager) WritePage(pageID common.PageID, pageData []byte) {
    offset := int64(pageID * common.PageSize)
    d.numWrites++
    d.dbIO.Seek(offset, os.SEEK_SET)
    bytesWritten, err := d.dbIO.Write(pageData)

    // check for I/O error
    if bytesWritten != common.PageSize || err != nil {
        level.Debug(common.Logger).Log("I/O error while writing")
        return
    }

    // make sure it goes on disk
    d.dbIO.Sync()
}

/**
* Read a page from the database file.
* @param pageID id of the page
* @param[out] pageData output buffer
*/
func (d *DiskManager) ReadPage(pageID common.PageID, pageData []byte) {
    offset := int64(pageID * common.PageSize)
    fInfo, _ := d.dbIO.Stat()
    if offset > fInfo.Size() {
        level.Debug(common.Logger).Log("I/O error past end of file")
        return
    }

    d.dbIO.Seek(offset, os.SEEK_SET)
    bytesRead, err := d.dbIO.Read(pageData)
    if err != nil {
        level.Debug(common.Logger).Log("I/O error while reading")
        return
    }
    // read less than PageSize
    if (bytesRead < common.PageSize) {
        level.Debug(common.Logger).Log("read less than a page")
        // zero out the rest of the page
        for i := 0; i < common.PageSize; i++ {
            pageData[i] = 0
        }
    }
}

/**
* Flush the entire log buffer into disk.
* @param logData raw log data
* @param size size of log entry
*/
func (d *DiskManager) WriteLog(logData []byte, size int) {
    if bufferUsed != nil {
        if &logData[0] == &bufferUsed[0] {
            log.Fatalln("logData == bufferUsed")
        }
    }
    bufferUsed = logData

    if (size == 0) {
        return
    }

    d.flushLog = true

    if d.flushLogF != nil {
        _, err := d.flushLogF.GetResult()
        if err != nil {
            log.Fatalln(err.Error())
        }
    }

    d.numFlushes++

    n, err := d.logIO.Write(logData[:size])

    if n < size || err != nil {
        level.Debug(common.Logger).Log("I/O error while writing log")
        return
    }
    // make sure writing is on disk
    d.logIO.Sync()
    d.flushLog = false;
}

/**
* Read a log entry from the log file.
* @param[out] logData output buffer
* @param size size of the log entry
* @param offset offset of the log entry in the file
* @return true if the read was successful, false otherwise
*/
func (d *DiskManager) ReadLog(logData []byte, size int, offset int64) bool {
    fInfo, _ := d.logIO.Stat()
    if offset > fInfo.Size() {
        level.Debug(common.Logger).Log("I/O error past end of file")
        return false
    }

    d.logIO.Seek(offset, os.SEEK_SET)
    bytesRead, err := d.logIO.Read(logData[:size])
    if err != nil {
        level.Debug(common.Logger).Log("I/O error while reading log")
        return false
    }

    // File ends before size, zero out remaining bytes
    if (bytesRead < size) {
        for i := bytesRead; i < size; i++ {
            logData[i] = 0
        }
    }

    return true
}

/**
* Allocate a page on disk.
* @return the id of the allocated page
*/
func (d *DiskManager) AllocatePage() common.PageID {
    /* stupid Go! I can't just do
     * return d.nextPageID++ */
    ret := d.nextPageID
    d.nextPageID++
    return ret
}

/**
* Deallocate a page on disk.
* @param pageID id of the page to deallocate
* Does nothing for now
*/
func (d *DiskManager) DeallocatePage(pageID common.PageID) {
}

/** @return the number of disk flushes */
func (d *DiskManager) GetNumFlushes() int {
    return d.numFlushes
}

/** @return true iff the in-memory content has not been flushed yet */
func (d *DiskManager) GetFlushState() bool {
    return d.flushLog
}

/** @return the number of disk writes */
func (d *DiskManager) GetNumWrites() int {
    return d.numWrites
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
