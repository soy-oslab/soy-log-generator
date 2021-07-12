package scheduler

import (
	"bufio"
	"errors"
	"fmt"
	w "github.com/soyoslab/soy_log_generator/pkg/watcher"
	"log"
	"os"
	"sync/atomic"
	"testing"
	"time"
)

type EvalFunc func(s *Scheduler, err error) bool

const FULL_CONFIG_TEXT = `{
    "targetIp": "localhost",
    "targetPort": "8972",
    "hotRingCapacity": %d,
    "coldRingCapacity": %d,
    "coldTimeoutMilli": 3000,
    "hotRingThreshold": 2,
    "coldRingThreshold": 2,
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

const CONFIG_TEXT = `{
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

func (s *Scheduler) getWatcher() *w.Watcher {
	return s.watcher
}

func setup(prefix string, configText string) string {
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

	return configFile.Name()
}

func teardown(filename string) {
	os.Remove(filename)
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
	filename := setup(testName, text)
	defer teardown(filename)
	s, err := InitScheduler(filename, submit, nil)
	if f(s, err) {
		t.Errorf("%v failed", testName)
	}
	defer s.Close()
}

func TestInitScheduler(t *testing.T) {
	evalFunc := func(s *Scheduler, err error) bool {
		return s == nil || err != nil
	}
	fileContentsTest(t, getSubmit(), "watcher-test-init-scheduler-valid-full", getConfig(FULL_CONFIG_TEXT, 1, 2), evalFunc)
	fileContentsTest(t, getSubmit(), "watcher-test-init-scheduler-valid-partial", getConfig(CONFIG_TEXT, 1, 2), evalFunc)
}

func TestNilClose(t *testing.T) {
	filename := setup("watcher-test-nil-close", getConfig(CONFIG_TEXT, 1, 2))
	defer teardown(filename)
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
	fileContentsTest(t, submit, "watcher-test-init-scheduler-invalid-hot-submit", getConfig(CONFIG_TEXT, 1, 2), evalFunc)
	submit = getSubmit()
	submit.Cold = nil
	fileContentsTest(t, submit, "watcher-test-init-scheduler-invalid-cold-submit", getConfig(CONFIG_TEXT, 1, 2), evalFunc)
	fileContentsTest(t, submit, "watcher-test-init-scheduler-invalid-hot-ring-capacity", getConfig(CONFIG_TEXT, 0, 2), evalFunc)
	fileContentsTest(t, submit, "watcher-test-init-scheduler-invalid-cold-ring-capacity", getConfig(CONFIG_TEXT, 1, 1), evalFunc)
}

func TestIsHotString(t *testing.T) {
	filename := setup("watcher-test-is-hot-string", getConfig(CONFIG_TEXT, 1, 2))
	defer teardown(filename)
	s, _ := InitScheduler(filename, getSubmit(), func(_ string) bool { return false })
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
	log.Fatalf("test-run timeout detected")
}

func TestRun(t *testing.T) {
	var hotCounter int32 = 0
	var coldCounter int32 = 0

	filename := setup("watcher-test-run", getConfig(CONFIG_TEXT, 1, 2))
	defer teardown(filename)
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
