package scheduler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync/atomic"

	"github.com/cloudflare/ahocorasick"
	defaults "github.com/mcuadros/go-defaults"
	w "github.com/soyoslab/soy_log_generator/pkg/watcher"
)

// InitScheduler initializes a Scheduler structure
func InitScheduler(configFilepath string, submitOperations SubmitOperations, customFilter CustomFilterFunc) (*Scheduler, error) {
	var err error
	s := new(Scheduler)
	if err = s.initConfig(configFilepath); err != nil {
		goto exception
	}
	if err = s.initWatcher(); err != nil {
		goto exception
	}
	atomic.StoreInt32(&s.IsRun, 0)
	s.customFilter = customFilter
	s.initHotFilter(s.config.Files)
	if s.config.HotRingCapacity < 1 {
		err = errors.New("hot ring capacity must be over 1")
		goto exception
	}
	s.hot.Init(s.config.HotRingCapacity, "hot")
	if s.config.ColdRingCapacity < 2 {
		err = errors.New("cold ring capacity must be over 2")
		goto exception
	}
	s.cold.Init(s.config.ColdRingCapacity, "cold")
	s.submit = submitOperations
	if s.submit.Hot == nil || s.submit.Cold == nil {
		err = errors.New("invalid submit function")
		goto exception
	}
	return s, err

exception:
	s.Close()
	return nil, err
}

// toLowerStrings makes capital characters to lower
func (s *Scheduler) toLowerStrings(arr []string) []string {
	for i, v := range arr {
		arr[i] = strings.ToLower(v)
	}
	return arr
}

// initHotFilter initializes the hot filtering in a Scheduler structure
// matcher uses the aho-corasick algorithm
func (s *Scheduler) initHotFilter(files []File) {
	s.matcher = make(map[string]*ahocorasick.Matcher)
	for _, file := range files {
		file.HotFilter = s.toLowerStrings(file.HotFilter)
		s.matcher[file.Filename] = ahocorasick.NewStringMatcher(file.HotFilter)
	}
}

// initConfig initializes a Config structure
func (s *Scheduler) initConfig(configFilepath string) error {
	var err error
	var b []byte

	defaults.SetDefaults(&s.config)
	b, err = ioutil.ReadFile(configFilepath)
	if os.IsNotExist(err) {
		goto out
	}
	err = json.Unmarshal(b, &s.config)
	if err != nil {
		goto out
	}

	if len(s.config.Files) == 0 {
		err = errors.New("file must exist at least one file")
		goto out
	}

out:
	return err
}

// initWatcher initializes the watcher package in a Scheduler structure
func (s *Scheduler) initWatcher() error {
	watcher, err := w.NewWatcher()
	s.watcher = watcher
	return err
}

// Close returns the resource related on the scheduling
func (s *Scheduler) Close() {
	if s != nil && s.watcher != nil {
		select {
		case err := <-s.watcher.GetErrorChannel():
			log.Panicln(err)
		default:
			break
		}
		s.watcher.Close()
	}
}
