package buffering_test

import (
	"errors"
	"io"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/soyoslab/soy_log_generator/pkg/buffering"
)

func setup(prefix string) (*buffering.Buffering, error) {
	testFile, err := os.CreateTemp("", prefix)
	if err != nil {
		log.Fatalf("test file creation failed: %v", err)
	}
	buffering, err := buffering.NewBuffering(testFile.Name(), func(str string, args interface{}) error {
		s, ok := args.([]interface{})
		if !ok {
			log.Println(strings.Trim(str, "\n"), s)
		} else {
			log.Println(strings.Trim(str, "\n"), args)
		}
		return nil
	})
	return buffering, err
}

func teardown(buffering *buffering.Buffering) {
	buffering.Close()
	os.Remove(buffering.GetFile().Name())
}

func TestNewBufferingSuccess(t *testing.T) {
	buffering, err := setup("test-new-buffering-success")
	defer teardown(buffering)
	if err != nil || buffering == nil {
		t.Errorf("buffering generation failed %v", err)
	}
}

func TestNewBufferingFailed_InvalidFileName(t *testing.T) {
	buffering, err := buffering.NewBuffering("", func(str string, args interface{}) error { return nil })
	if buffering != nil && err == nil {
		t.Errorf("invalid file name test failed %v", err)
	}
}

func TestNewBufferingFailed_NilFunction(t *testing.T) {
	testFile, err := os.CreateTemp("", "test-new-buffering-success")
	if err != nil {
		log.Fatalf("test file creation failed: %v", err)
	}
	defer testFile.Close()
	defer os.Remove(testFile.Name())
	buffering, err := buffering.NewBuffering(testFile.Name(), nil)
	if buffering != nil && err == nil {
		t.Errorf("nil function test failed %v", err)
	}
}

func TestGetFileSize(t *testing.T) {
	buffering, _ := setup("test-get-file-size")
	defer teardown(buffering)

	_, err := buffering.GetFileSize()
	if err != nil {
		t.Errorf("get file size failed %v", err)
	}
}

func TestFileSizeChanged(t *testing.T) {
	buffering, _ := setup("test-file-size-changed")
	defer teardown(buffering)

	targetFile, _ := os.OpenFile(buffering.GetFile().Name(), os.O_WRONLY|os.O_CREATE, 0755)
	defer targetFile.Close()
	n, err := targetFile.WriteString("01234567890123456789\n")
	if err != nil || n == 0 {
		log.Fatalf("testFile write failed %v\n", err)
	}

	buffering.GetFile().Seek(0, io.SeekEnd)
	if v, _ := buffering.IsValidFileSize(); !v {
		t.Errorf("is valid but return is invalid")
	}
	// underflow location move to end
	buffering.GetFile().Seek(0, io.SeekStart)
	if v, _ := buffering.IsValidFileSize(); !v {
		t.Errorf("is valid but return is invalid")
	}
	buffering.GetFile().Seek(1, io.SeekEnd)
	if v, _ := buffering.IsValidFileSize(); v {
		t.Errorf("is invalid but return is valid")
	}
}

func TestUpdateToValidOffset(t *testing.T) {
	buffering, _ := setup("test-update-to-valid-offset")
	defer teardown(buffering)
	targetFile, _ := os.OpenFile(buffering.GetFile().Name(), os.O_WRONLY|os.O_CREATE, 0755)
	defer targetFile.Close()
	n, err := targetFile.WriteString("01234567890123456789\n")
	if err != nil || n == 0 {
		log.Fatalf("testFile write failed %v\n", err)
	}

	target, _ := buffering.GetFile().Seek(0, io.SeekEnd)
	buffering.GetFile().Seek(1, io.SeekEnd)
	buffering.UpdateToValidOffset()
	current, _ := buffering.GetFile().Seek(0, io.SeekCurrent)
	if current != target {
		t.Errorf("update to valid offset is failed %v->%v", current, target)
	}
}

func writeFiles(filename string) []string {
	targetFile, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0755)
	stringList := []string{"0\n", "11\n", "222\n", "3333\n", "44444\n"}
	for _, str := range stringList {
		n, err := targetFile.WriteString(str)
		if err != nil || n == 0 {
			log.Fatalf("testFile write failed %v\n", err)
		}
	}
	targetFile.Close()
	return stringList
}

func TestDoReadLines(t *testing.T) {
	buffering, _ := setup("test-do-read-lines")
	defer teardown(buffering)
	stringList := writeFiles(buffering.GetFile().Name())

	i := 0
	buffering.SetProcessingFunction(func(str string, _ interface{}) error {
		if strings.Compare(str, stringList[i]) != 0 {
			t.Errorf("string is not equal %s <> %s", str, stringList[i])
		}
		i++
		return nil
	})
	_, err := buffering.DoReadLines()
	if err != nil {
		t.Errorf("do readlines failed %v", err)
	}
}

func TestDoReadLinesNil(t *testing.T) {
	buffering, _ := setup("test-do-read-lines-nil")
	defer teardown(buffering)
	buffering.SetProcessingFunction(nil)
	_, err := buffering.DoReadLines()
	if err == nil {
		t.Errorf("do readlines nil function failed %v", err)
	}
}

func TestDoReadLinesInvalidFunction(t *testing.T) {
	buffering, _ := setup("test-do-read-lines-invalid-function")
	defer teardown(buffering)
	_ = writeFiles(buffering.GetFile().Name())
	buffering.SetProcessingFunction(func(_ string, _ interface{}) error {
		return errors.New("sample error")
	})
	_, err := buffering.DoReadLines()
	if err == nil {
		t.Errorf("do readlines invalid function %v", err)
	}
}
