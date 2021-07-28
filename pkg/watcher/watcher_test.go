package watcher

import (
	"errors"
	"io"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
)

func setup() (*Watcher, error) {
	w, err := NewWatcher()
	log.SetFlags(log.Lshortfile)
	return w, err
}

func makeFile(prefix string) *os.File {
	testFile, err := os.CreateTemp("", prefix)
	if err != nil {
		log.Fatalf("test file creation failed: %v", err)
	}
	return testFile
}

func teardown(w *Watcher) {
	for filename := range w.GetFileInfoTable() {
		os.Remove(filename)
	}
	w.Close()
}

func TestNewWatcher(t *testing.T) {
	w, err := setup()
	defer teardown(w)
	if err != nil {
		t.Errorf("watcher generation failed")
	}
}

func TestValidAddFile(t *testing.T) {
	file := makeFile("test-add-file")
	defer file.Close()

	w, _ := setup()
	defer teardown(w)

	err := w.AddFile(file.Name(), func(str string, args interface{}) error {
		return nil
	})
	if err != nil {
		t.Errorf("valid add file test failed (cannot execute the AddFile)")
	}
	info, err := w.GetFileInfo(file.Name())
	if err != nil || info.GetBuffer().GetFile().Name() != file.Name() {
		t.Errorf("valid add file test failed (cannot find added a file)")
	}
}

func TestInvalidFunctionAddFile(t *testing.T) {
	w, _ := setup()
	defer teardown(w)

	file := makeFile("test-add-file-invalid-function")
	defer file.Close()

	err := w.AddFile(file.Name(), nil)
	if err == nil {
		t.Errorf("invalid(nil function) add file test failed (cannot execute")
	}
}

func TestInvalidFileNameAddFile(t *testing.T) {
	w, _ := setup()
	defer teardown(w)

	file := makeFile("test-add-file-invalid-file")
	defer file.Close()
	err := w.AddFile("", func(str string, args interface{}) error {
		return nil
	})
	if err == nil {
		t.Errorf("invalid(empty filename) add file test failed")
	}
}

func TestSpectatorOneFile(t *testing.T) {
	w, _ := setup()
	defer teardown(w)

	file := makeFile("test-spectator")
	defer file.Close()

	sampleString := []string{"test1\n", "test2\n", "test3\n"}

	index := 0
	err := w.AddFile(file.Name(), func(str string, args interface{}) error {
		rv := reflect.ValueOf(args)
		file := rv.Index(1).Interface().(*os.File)
		if cur, _ := file.Seek(0, io.SeekCurrent); cur%6 != 0 {
			t.Errorf("invalid write pointer")
		}
		if sampleString[index] != str {
			t.Errorf("write Detection Failed")
		}
		index++
		return nil
	})
	if err != nil {
		t.Errorf("add file failed")
	}
	for _, str := range sampleString {
		file.WriteString(str)
		file.Sync()
	}
	if index != 3 {
		t.Errorf("event detection failed")
	}
}

func TestSpectatorThreeFile(t *testing.T) {
	w, _ := setup()
	defer teardown(w)

	var files []*os.File = []*os.File{makeFile("test-spectator-1-"), makeFile("test-spectator-2-"), makeFile("test-spectator-3-")}

	sampleString := []string{"test1\n", "test2\n", "test3\n"}
	stringIndex := 0
	for _, file := range files {
		err := w.AddFile(file.Name(), func(str string, args interface{}) error {
			if sampleString[stringIndex] != str {
				t.Errorf("write detection failed")
			}
			stringIndex++
			return nil
		})
		if err != nil {
			t.Errorf("add file failed")
		}
	}

	for idx, str := range sampleString {
		files[idx].WriteString(str)
		files[idx].Close()
	}

	timer := time.AfterFunc(time.Second*60, func() {
		t.Errorf("timeout")
	})

	// infinite loop
	for stringIndex != 3 {
	}

	timer.Stop()

	// finish
}

func TestInvalidGetFileInfo(t *testing.T) {
	w, _ := setup()
	defer teardown(w)
	w.GetFileInfo("")
}

func TestInvalidProcessFile(t *testing.T) {
	w, _ := setup()
	defer teardown(w)
	w.ProcessFile("")
}

func TestValidRemoveTest(t *testing.T) {
	w, _ := setup()
	defer teardown(w)

	var files []*os.File = []*os.File{makeFile("test-spectator-1-"), makeFile("test-spectator-2-"), makeFile("test-spectator-3-")}

	for _, file := range files {
		err := w.AddFile(file.Name(), func(str string, args interface{}) error {
			return nil
		})
		if err != nil {
			t.Errorf("add file failed")
		}
	}
	_, err := w.GetFileInfo(files[0].Name())
	if err != nil {
		t.Errorf("valid get file info test failed")
	}

	w.Remove(files[0].Name())

	_, err = w.GetFileInfo(files[0].Name())
	if err == nil {
		t.Errorf("invalid get file info test failed")
	}

	// finish
}

func TestInvalidRemoveTest(t *testing.T) {
	w, _ := setup()
	defer teardown(w)
	w.Remove("")
}

func TestInvalidEventProcessor(t *testing.T) {
	w, _ := setup()
	defer teardown(w)
	err := errors.New("test")
	w.GetNotifier().Errors <- err
	w.Wait()
	select {
	case _ = <-w.GetErrorChannel():
	default:
		t.Errorf("invalid event processor state test failed")
	}
}

func TestInvalidSpectatorEventName(t *testing.T) {
	w, _ := setup()
	defer teardown(w)
	event := fsnotify.Event{}
	event.Op = fsnotify.Write
	w.GetNotifier().Events <- event
	w.Wait()
	select {
	case _ = <-w.GetErrorChannel():
	default:
		t.Errorf("invalid event processor state test failed")
	}
}
