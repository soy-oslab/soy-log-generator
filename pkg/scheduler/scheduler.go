package scheduler

import (
	"github.com/cloudflare/ahocorasick"
	"github.com/soyoslab/soy_log_generator/pkg/ring"
	w "github.com/soyoslab/soy_log_generator/pkg/watcher"
)

// SubmitFunc is the function pointer of the submit message
type SubmitFunc func(messages []Message) error

// CustomFilterFunc is the function pointer of the custom hot/cold filtering
type CustomFilterFunc func(str string, isHot bool) bool

// File contains the each file's information in json manner
type File struct {
	Filename  string   `json:"filename"`
	HotFilter []string `json:"hotFilter"`
}

// Config contains the application running configurations in json manner
type Config struct {
	Namespace         string `json:"namespace" default:"anonymous"`
	TargetIP          string `json:"targetIp" default:"localhost"`
	TargetPort        string `json:"targetPort" default:"8972"`
	HotRingCapacity   uint64 `json:"hotRingCapacity" default:"32"`
	ColdRingCapacity  uint64 `json:"coldRingCapacity" default:"32"`
	ColdTimeout       uint64 `json:"coldTimeoutMilli" default:"5000"`
	PollingInterval   uint64 `json:"pollingIntervalMilli" default:"1000"`
	Files             []File `json:"files"`
	HotRingThreshold  uint64 `json:"hotRingThreshold" default:"0"`
	ColdRingThreshold uint64 `json:"coldRingThreshold" default:"0"`
	ColdSendThreshold uint64 `json:"coldSendThresholdBytes" default:"4096"`
}

// FileInfo contains the file data block metadata
type FileInfo struct {
	Timestamp int64
	Filename  string
	Length    uint64
}

// Message structure is used to transport with log-collector
type Message struct {
	Info FileInfo
	Data []byte
}

// SubmitOperations contains functions which contain the transport logic
type SubmitOperations struct {
	Hot  SubmitFunc
	Cold SubmitFunc
}

// Scheduler contains the scheduling information
type Scheduler struct {
	watcher      *w.Watcher
	config       Config
	hot          ring.Ring
	cold         ring.Ring
	matcher      map[string]*ahocorasick.Matcher
	submit       SubmitOperations
	customFilter CustomFilterFunc
	IsRun        int32
}

// GetConfig returns Config structure in Scheduler
func (s *Scheduler) GetConfig() Config {
	return s.config
}
