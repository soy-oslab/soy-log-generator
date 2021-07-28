package scheduler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
	if s.config.PollingInterval > 1000 {
		log.Fatalf("Please set the polling interval to below than 1000ms (current: %dms)\n", s.config.PollingInterval)
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

func getConfigFiles(filenames []string, meta File) []File {
	files := []File{}
	for _, filename := range filenames {
		file := File{filename, meta.HotFilter}
		files = append(files, file)
	}
	return files
}

func configPatternTranslation(metaList []File) ([]File, error) {
	files := []File{}
	for _, meta := range metaList {
		matches, err := filepath.Glob(meta.Filename)
		if len(matches) == 0 || err != nil {
			return nil, fmt.Errorf("matches error detected (str:%s;err:%v;matches:%v)", meta.Filename, err, matches)
		}
		files = append(files, getConfigFiles(matches, meta)...)
	}
	return files, nil
}

// initConfig initializes a Config structure
func (s *Scheduler) initConfig(configFilepath string) error {
	var (
		err    error
		b      []byte
		dupMap map[string]bool
	)

	defaults.SetDefaults(&s.config)
	b, err = ioutil.ReadFile(configFilepath)
	if os.IsNotExist(err) {
		goto out
	}
	err = json.Unmarshal(b, &s.config)
	if err != nil {
		goto out
	}

	s.config.Files, err = configPatternTranslation(s.config.Files)
	if err != nil {
		goto out
	}

	if len(s.config.Files) == 0 || len(s.config.Files) >= 256 {
		err = errors.New("number of files is over 0 and under 256")
		goto out
	}

	dupMap = make(map[string]bool)
	for _, file := range s.config.Files {
		if _, ok := dupMap[file.Filename]; ok {
			err = errors.New("duplicated file doesn't allow")
			goto out
		}
		dupMap[file.Filename] = true
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
	if s == nil {
		return
	}
	atomic.StoreInt32(&s.IsRun, 0)
	if s.watcher == nil {
		return
	}
	select {
	case err := <-s.watcher.GetErrorChannel():
		if !strings.Contains(err.Error(), os.ErrClosed.Error()) && err != nil {
			log.Panicln(err, os.ErrClosed)
		}
	default:
		break
	}
	s.hot.Close()
	s.cold.Close()
	s.watcher.Close()
	defer func() {
		recover()
	}()
	select {
	case _, ok := <-s.hot.Kick:
		if ok {
			close(s.hot.Kick)
		}
	case _, ok := <-s.cold.Kick:
		if ok {
			close(s.cold.Kick)
		}
	default:
	}
}
