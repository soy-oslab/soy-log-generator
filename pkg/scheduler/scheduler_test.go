package scheduler

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"
	"sync/atomic"
	"testing"
	"time"

	w "github.com/soyoslab/soy_log_generator/pkg/watcher"
)

type EvalFunc func(s *Scheduler, err error) bool

const FullConfigText = `{
    "targetIp": "localhost",
    "targetPort": "8972",
    "hotRingCapacity": %d,
    "coldRingCapacity": %d,
    "coldTimeoutMilli": 1000,
    "hotRingThreshold": 0,
    "coldRingThreshold": 0,
    "pollingIntervalMilli": 1000,
    "files": [
        {
            "filename": "%v",
            "hotFilter": [
                "error",
                "critical"
            ]
        },
        {
            "filename": "%v",
            "hotFilter": [
                "critical",
                "warn"
            ]
        }
    ]
  }`

const ConfigText = `{
    "hotRingThreshold": 0,
    "coldRingThreshold": 0,
    "hotRingCapacity": %d,
    "coldRingCapacity": %d,
    "coldTimeoutMilli": 30,
    "pollingIntervalMilli": 100,
    "files": [
        {
            "filename": "%v",
            "hotFilter": [
                "error",
                "critical"
            ]
        }
    ]
  }`

func (s *Scheduler) getWatcher() *w.Watcher {
	return s.watcher
}

func setup(prefix string, configText string) (string, string) {
	log.SetFlags(log.Lshortfile)
	testFile, err := os.CreateTemp("", prefix)
	if err != nil {
		log.Fatalf("config file creation failed: %v", err)
	}
	defer testFile.Close()
	configFile, err := os.OpenFile(testFile.Name(), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.FileMode(0644))
	if err != nil {
		log.Fatalf("config file open failed: %v", err)
	}
	defer configFile.Close()
	writer := bufio.NewWriter(configFile)
	writer.WriteString(configText)
	writer.Flush()

	return testFile.Name(), configFile.Name()
}

func teardown(filenames []string) {
	for _, filename := range filenames {
		os.Remove(filename)
	}
}

func getConfig(configText string, hotRingCapacity uint64, coldRingCapacity uint64) string {
	testFile1, err := os.CreateTemp("", "test1.txt")
	if err != nil {
		log.Fatalf("temproary file-1 creation failed")
	}
	testFile2, err := os.CreateTemp("", "test2.txt")
	if err != nil {
		log.Fatalf("temproary file-2 creation failed")
	}
	return fmt.Sprintf(configText, hotRingCapacity, coldRingCapacity, testFile1.Name(), testFile2.Name())
}

func getSubmit() SubmitOperations {
	hot := func(message []Message) error {
		return nil
	}
	cold := func(message []Message) error {
		return nil
	}
	submit := SubmitOperations{hot, cold}
	return submit
}

func fileContentsTest(t *testing.T, submit SubmitOperations, testName string, text string, f EvalFunc) {
	testFilename, filename := setup(testName, text)
	defer teardown([]string{filename, testFilename})
	s, err := InitScheduler(filename, submit, nil)
	if f(s, err) {
		t.Errorf("%v failed", testName)
	}
	defer s.Close()
}

func TestDuplication(t *testing.T) {
	evalFunc := func(_ *Scheduler, err error) bool { return err == nil }
	testFile, err := os.CreateTemp("", "test1.txt")
	if err != nil {
		log.Fatalf("temproary file-1 creation failed")
	}
	defer teardown([]string{testFile.Name()})
	config := fmt.Sprintf(FullConfigText, 1, 2, testFile.Name(), testFile.Name())
	fileContentsTest(t, getSubmit(), "watcher-test-duplication-test", config, evalFunc)
}

func getPatternConfig(pattern string) ([]string, string) {
	file1, err := os.CreateTemp("", "pattern-test.txt")
	if err != nil {
		log.Fatalf("temproary file-1 creation failed")
	}
	file2, err := os.CreateTemp("", "pattern-test.txt")
	if err != nil {
		log.Fatalf("temproary file-2 creation failed")
	}
	config := fmt.Sprintf(ConfigText, 1, 2, pattern)
	return []string{file1.Name(), file2.Name()}, config
}

func TestPatternFiles(t *testing.T) {
	evalFunc := func(s *Scheduler, err error) bool {
		if err != nil {
			return true
		}
		if len(s.GetConfig().Files) != 2 {
			return true
		}
		return false
	}
	filenames, config := getPatternConfig("/tmp/pattern-test.txt*")
	defer teardown(filenames)
	fileContentsTest(t, getSubmit(), "pattern-valid-test", config, evalFunc)
}

func TestPatternFilesInvalid(t *testing.T) {
	evalFunc := func(s *Scheduler, err error) bool {
		if err == nil {
			return true
		}
		return false
	}
	filenames, config := getPatternConfig("/tmp/*pattern-test.txt")
	defer teardown(filenames)
	fileContentsTest(t, getSubmit(), "pattern-invalid-test", config, evalFunc)
}

func TestInitScheduler(t *testing.T) {
	evalFunc := func(s *Scheduler, err error) bool {
		return s == nil || err != nil
	}
	fileContentsTest(t, getSubmit(), "watcher-test-init-scheduler-valid-full", getConfig(FullConfigText, 1, 2), evalFunc)
	fileContentsTest(t, getSubmit(), "watcher-test-init-scheduler-valid-partial", getConfig(FullConfigText, 1, 2), evalFunc)
}

func TestNilClose(t *testing.T) {
	testFilename, filename := setup("watcher-test-nil-close", getConfig(FullConfigText, 1, 2))
	defer teardown([]string{testFilename, filename})
	s, _ := InitScheduler(filename, getSubmit(), nil)
	watcher := s.getWatcher()
	defer func() {
		_ = recover()
	}()
	watcher.GetNotifier().Errors <- errors.New("TestNilClose")
	watcher.Wait()
	s.Close()
	t.Errorf("error detected in close sequence but it is ignored")
}

func TestInitSchedulerInvalid(t *testing.T) {
	evalFunc := func(_ *Scheduler, err error) bool { return err == nil }
	fileContentsTest(t, getSubmit(), "watcher-test-init-scheduler-invalid-contents", "", evalFunc)
	fileContentsTest(t, getSubmit(), "watcher-test-init-scheduler-invalid-files-count", `{"files":[]}`, evalFunc)
	submit := getSubmit()
	submit.Hot = nil
	fileContentsTest(t, submit, "watcher-test-init-scheduler-invalid-hot-submit", getConfig(FullConfigText, 1, 2), evalFunc)
	submit = getSubmit()
	submit.Cold = nil
	fileContentsTest(t, submit, "watcher-test-init-scheduler-invalid-cold-submit", getConfig(FullConfigText, 1, 2), evalFunc)
	fileContentsTest(t, submit, "watcher-test-init-scheduler-invalid-hot-ring-capacity", getConfig(FullConfigText, 0, 2), evalFunc)
	fileContentsTest(t, submit, "watcher-test-init-scheduler-invalid-cold-ring-capacity", getConfig(FullConfigText, 1, 1), evalFunc)
}

func TestIsHotString(t *testing.T) {
	testFilename, filename := setup("watcher-test-is-hot-string", getConfig(FullConfigText, 1, 2))
	defer teardown([]string{testFilename, filename})
	s, _ := InitScheduler(filename, getSubmit(), func(_ string, isHot bool) bool { return isHot })
	if s.isHotString(s.GetConfig().Files[0].Filename, "warn e r r o r") {
		t.Errorf("cold sentence is evaluated to hot")
	}
	if !s.isHotString(s.GetConfig().Files[0].Filename, "error w a r n") {
		t.Errorf("hot sentence is evaluated to cold")
	}
	s.Close()

	s, _ = InitScheduler(filename, getSubmit(), nil)
	if s.isHotString(s.GetConfig().Files[1].Filename, "error w a r n") {
		t.Errorf("cold sentence is evaluated to hot")
	}
	if s.isHotString(s.GetConfig().Files[1].Filename, "") {
		t.Errorf("empty sentence is evaluated to hot")
	}
	if !s.isHotString(s.GetConfig().Files[1].Filename, "warn e r r o r") {
		t.Errorf("hot sentence is evaluated to cold")
	}

	defer func() {
		_ = recover()
		s.Close()
	}()
	s.isHotString("", "")
}

func background(s *Scheduler) {
	filename := s.GetConfig().Files[0].Filename
	file, err := os.OpenFile(filename, os.O_WRONLY, os.FileMode(0644))
	if err != nil {
		log.Fatalf("hooked file open failed")
	}
	for {
		if atomic.LoadInt32(&s.IsRun) == 1 {
			break
		}
	}
	log.Println(filename)
	file.WriteString("critical1\n")
	file.Sync()
	file.WriteString("cold1\n")
	file.Sync()
	file.WriteString("cold2\n")
	file.Sync()
	file.WriteString("cold3\n")
	file.Sync()
	file.WriteString("cold4\n")
	file.Sync()
	file.WriteString("critical2\n")
	file.Sync()
	file.Close()
}

func counter(t *testing.T, hotCounter *int32, coldCounter *int32, s *Scheduler) {
	start := time.Now()
	timeout := time.Duration(5) * time.Second
	for {
		if time.Since(start) > timeout {
			break
		}
		currentHotCount := atomic.LoadInt32(hotCounter)
		currentColdCount := atomic.LoadInt32(coldCounter)
		if currentHotCount == 2 && currentColdCount == 4 {
			s.Close()
		}
	}
	pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
	log.Fatalf("test-run timeout detected")
}

func TestRun(t *testing.T) {
	var hotCounter int32 = 0
	var coldCounter int32 = 0

	testFilename, filename := setup("watcher-test-run", getConfig(FullConfigText, 1, 2))
	defer teardown([]string{testFilename, filename})
	hot := func(messages []Message) error {
		for _, message := range messages {
			log.Println(string(message.Data))
			atomic.AddInt32(&hotCounter, 1)
		}
		return nil
	}
	cold := func(messages []Message) error {
		for _, message := range messages {
			log.Println(string(message.Data))
			atomic.AddInt32(&coldCounter, 1)
		}
		return nil
	}
	submit := SubmitOperations{hot, cold}
	s, _ := InitScheduler(filename, submit, nil)
	go background(s)
	go counter(t, &hotCounter, &coldCounter, s)
	s.Run()
}

func TestMain(m *testing.M) {
	files, _ := filepath.Glob("/tmp/pattern-test.txt*")
	for _, f := range files {
		os.Remove(f)
	}
	os.Exit(m.Run())
}
