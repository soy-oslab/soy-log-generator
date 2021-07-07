package buffering

import (
	"bufio"
	"errors"
	"io"
	"os"
)

// Buffering structure contains the file information and line processing function
type Buffering struct {
	file                   *os.File
	name                   string
	reader                 *bufio.Reader
	lineProcessingFunction func(string, interface{}) error
}

// NewBuffering makes a new structure based on Buffering type
func NewBuffering(filename string, processFunction func(string, interface{}) error) (*Buffering, error) {
	var err error = nil

	buffering := new(Buffering)

	buffering.name = filename
	if processFunction == nil {
		err = errors.New("buffering's process function must be specified")
		goto functionNilException
	}
	buffering.lineProcessingFunction = processFunction
	buffering.file, err = os.Open(filename)
	if err != nil {
		goto openException
	}
	_, err = buffering.file.Seek(0, io.SeekEnd)
	if err != nil {
		goto seekException
	}
	buffering.reader = bufio.NewReader(buffering.file)

	return buffering, err

seekException:
	buffering.Close()
openException:
functionNilException:
	return nil, err
}

// Close collects the resources in the Buffering structure
func (b *Buffering) Close() error {
	err := b.file.Close()
	return err
}

// GetFile returns *os.File
func (b *Buffering) GetFile() *os.File {
	return b.file
}

// GetFileSize returns the size of a file in buffering structure
func (b *Buffering) GetFileSize() (int64, error) {
	file := b.file
	stat, err := file.Stat()
	return stat.Size(), err
}

// IsValidFileSize checks the file size is changed
func (b *Buffering) IsValidFileSize() (bool, error) {
	var err error
	var fileSize, currentOffset int64 = 0, 0

	file := b.file
	fileSize, err = b.GetFileSize()
	if err != nil {
		goto exception
	}
	currentOffset, err = file.Seek(0, io.SeekCurrent)
exception:
	return currentOffset <= fileSize, err
}

// UpdateToValidOffset changes current offset to valid offset
func (b *Buffering) UpdateToValidOffset() {
	if isValid, _ := b.IsValidFileSize(); !isValid {
		b.file.Seek(0, io.SeekEnd)
	}
}

// DoReadLines does the read "lines" until encountering the EOF
// Note that DoReadLines()'s args directly pass to buffering's line processing functions.
// In other words, line processing function can hold the `[]interface{}` not `interface{}`.
// Therfore, you must think this to when you create the line processing function.
func (b *Buffering) DoReadLines(args ...interface{}) (int64, error) {
	var str string
	var err error = nil
	file := b.file
	offset, _ := file.Seek(0, io.SeekCurrent)

	if b.lineProcessingFunction == nil {
		err = errors.New("line processing function is not specified")
		goto exception
	}
	for {
		str, err = b.reader.ReadString('\n')
		if err == io.EOF {
			offset, err = file.Seek(offset, io.SeekStart)
			break
		} else if err != nil {
			goto exception
		}
		err = b.lineProcessingFunction(str, args)
		if err != nil {
			goto exception
		}
		offset += (int64(len(str)))
	}
exception:
	return offset, err
}

// SetProcessingFunction sets the processing line function
func (b *Buffering) SetProcessingFunction(processFunction func(string, interface{}) error) {
	b.lineProcessingFunction = processFunction
}
