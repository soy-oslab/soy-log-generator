package watcher

import (
	"errors"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/soyoslab/soy_log_generator/pkg/buffering"
)

// FileInfo structure contains the Buffering structure's pointer
type FileInfo struct {
	buffer *buffering.Buffering
}

// GetBuffer returns the pointer of the buffering structure
func (f *FileInfo) GetBuffer() *buffering.Buffering {
	return f.buffer
}

// Watcher structure contains the fsnotify structures and being watched files information
type Watcher struct {
	infoTable    map[string]FileInfo
	notifier     *fsnotify.Watcher
	workingGroup *sync.WaitGroup
	stop         chan bool
	errors       chan error
}

// NewWatcher creates Watcher structure
// Note that this function start to run the Spectator goroutine
func NewWatcher() (*Watcher, error) {
	var err error
	watcher := new(Watcher)
	watcher.workingGroup = new(sync.WaitGroup)
	watcher.infoTable = make(map[string]FileInfo)
	watcher.notifier, err = fsnotify.NewWatcher()

	watcher.workingGroup.Add(1)
	watcher.stop = make(chan bool)
	watcher.errors = make(chan error)
	go watcher.Spectator()

	return watcher, err
}

// GetErrorChannel returns the errors channel
func (w *Watcher) GetErrorChannel() <-chan error {
	return w.errors
}

// AddFile adds a file to the Watcher.
// During the adding file, this also creates the buffering structure
func (w *Watcher) AddFile(filename string, lineProcessingFunction func(string, interface{}) error) error {
	buffer, err := buffering.NewBuffering(filename, lineProcessingFunction)
	if err != nil {
		goto exception
	}
	w.infoTable[filename] = FileInfo{buffer}
	err = w.notifier.Add(filename)

exception:
	return err
}

// GetFileInfoTable return FileInfoTable
func (w *Watcher) GetFileInfoTable() map[string]FileInfo {
	return w.infoTable
}

// GetFileInfo returns the FileInfo structure pointer in the Watcher structure
func (w *Watcher) GetFileInfo(filename string) (FileInfo, error) {
	if info, ok := w.infoTable[filename]; ok {
		return info, nil
	}
	return FileInfo{}, errors.New("cannot find the buffering structure")
}

// ProcessFile prcesses each lines in a file
func (w *Watcher) ProcessFile(filename string) error {
	info, err := w.GetFileInfo(filename)
	if err != nil {
		goto exception
	}
	info.buffer.UpdateToValidOffset()
	_, err = info.buffer.DoReadLines(filename, info.buffer.GetFile())
exception:
	return err
}

// EventProcessor is generic event processor.
// Note that you must take care of other case always returns it is valid state
func (w *Watcher) EventProcessor(event fsnotify.Event) error {
	var err error = nil
	if event.Op&fsnotify.Write == fsnotify.Write {
		err = w.ProcessFile(event.Name)
	}
	return err
}

// GetNotifier returns the fsnotify.Watcher's address
func (w *Watcher) GetNotifier() *fsnotify.Watcher {
	return w.notifier
}

// Spectator runs the infinitely before it receives the stop signal
// Also it passes an event to the EventProcessor
func (w *Watcher) Spectator() {
	var err error = nil
	for {
		select {
		case _ = <-w.stop:
			err = nil
			goto exception
		case event, _ := <-w.notifier.Events:
			err = w.EventProcessor(event)
			if err != nil {
				goto exception
			}
			// something to do
		case err, _ = <-w.notifier.Errors:
			goto exception
		}
	}
exception:
	w.workingGroup.Done()
	w.errors <- err
}

// Remove removes the specific file hook
func (w *Watcher) Remove(filename string) error {
	info, err := w.GetFileInfo(filename)
	if err != nil {
		return err
	}
	info.GetBuffer().Close()
	delete(w.infoTable, filename)
	return w.notifier.Remove(filename)
}

// Wait waits the working group in the watcher
func (w *Watcher) Wait() {
	w.workingGroup.Wait()
}

// Stop signals the stop signal to the EventProcessor
func (w *Watcher) Stop() {
	w.stop <- true
}

// Close frees resources
func (w *Watcher) Close() error {
	for _, info := range w.infoTable {
		info.buffer.Close()
	}
	close(w.stop)
	return w.notifier.Close()
}
