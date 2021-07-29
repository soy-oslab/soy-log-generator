package scheduler

import (
	"log"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

// insertHotString inserts the string in the hot classifier manner
func (s *Scheduler) insertHotString(message Message) error {
	defer func() {
		recover()
	}()
	err := s.hot.Push(message)
	s.hot.Kick <- true
	return err
}

// insertColdString inserts the string in the cold classifier mannner
func (s *Scheduler) insertColdString(message Message) error {
	defer func() {
		recover()
	}()
	start := time.Now()
	timeout := time.Duration(s.config.ColdTimeout) * time.Millisecond
	for {
		ok, _ := s.cold.Offer(message)
		if ok {
			break
		}

		if timeout > 0 && time.Since(start) >= timeout {
			s.cold.Kick <- true
			start = time.Now()
		} else {
			runtime.Gosched()
		}
		if atomic.LoadInt32(&s.IsRun) == 0 {
			break
		}
	}
	return nil
}

// insertString classifies the string state and place to the valid method
func (s *Scheduler) insertString(str string, args interface{}) error {
	filename := args.([]interface{})[0].(string)
	message := Message{}
	message.Info.Timestamp = time.Now().UnixNano()
	message.Info.Filename = filename
	str = strings.Trim(str, "\n")
	message.Info.Length = uint64(len([]byte(str)))
	message.Data = []byte(str)
	if s.isHotString(filename, str) {
		go s.insertHotString(message)
	} else {
		go s.insertColdString(message)
	}
	return nil
}

// registFilesToWatcher regists the files to watcher package in the Scheduler structure
func (s *Scheduler) registFilesToWatcher() error {
	var err error
	for _, file := range s.config.Files {
		err = s.watcher.AddFile(file.Filename, s.insertString)
		if err != nil {
			goto exception
		}
	}
	return nil
exception:
	return err
}

// process processes the message based on the given function
func (s *Scheduler) process(f SubmitFunc, arr []interface{}) error {
	messages := make([]Message, len(arr))
	for i, v := range arr {
		messages[i] = v.(Message)
	}
	return f(messages)
}

// processString prcesses the string based on the scheduling policy
func (s *Scheduler) processString() {
	for {
		if atomic.LoadInt32(&s.IsRun) == 0 {
			break
		}
		select {
		case <-s.hot.Kick:
			s.process(s.submit.Hot, s.hot.Pop(s.config.HotRingThreshold))
			continue
		case <-s.cold.Kick:
			s.process(s.submit.Cold, s.cold.Pop(s.config.ColdRingThreshold))
			continue
		case <-time.After(time.Duration(s.config.PollingInterval) * time.Millisecond):
			s.process(s.submit.Hot, s.hot.Pop(s.config.HotRingThreshold))
			s.process(s.submit.Cold, s.cold.Pop(s.config.ColdRingThreshold))
		}
	}
}

// Run executes the scheduler
func (s *Scheduler) Run() error {
	err := s.registFilesToWatcher()
	if err != nil {
		return err
	}
	go s.processString()
	atomic.StoreInt32(&s.IsRun, 1)
	s.watcher.Wait()
	return nil
}

// isHotString classifies string is hot or not
// Note that if you set the s.customFilter then it will work after keywords check.
func (s *Scheduler) isHotString(filename string, str string) bool {
	str = strings.ToLower(str)
	matcher, ok := s.matcher[filename]
	if !ok {
		log.Panicf("invalid filename detected %v", filename)
	}
	isHot := len(matcher.MatchThreadSafe([]byte(str))) > 0
	if s.customFilter != nil {
		return s.customFilter(str, isHot)
	}
	return isHot
}
