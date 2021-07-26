package transport

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	rpcx "github.com/smallnest/rpcx/client"
	"github.com/soyoslab/soy_log_collector/pkg/rpc"
	c "github.com/soyoslab/soy_log_generator/pkg/compressor"
	s "github.com/soyoslab/soy_log_generator/pkg/scheduler"
)

// SubmitFunc is a type for submission the packet to rpcx
type SubmitFunc func(*rpc.LogMessage, rpcx.XClient) error

// BufferingMetadata contains the cold data's buffering information
type BufferingMetadata struct {
	packet    rpc.LogMessage
	start     time.Time
	threshold uint64
	timeout   time.Duration
}

// Port contains the rpcx client and the buffering information
type Port struct {
	xclient rpcx.XClient
	meta    BufferingMetadata
}

// Close closes the rpcx client
func (p *Port) Close() {
	if p.xclient != nil {
		p.xclient.Close()
	}
}

// Transport contains the rpcx and communcation information
type Transport struct {
	scheduler  *s.Scheduler
	hot        Port
	cold       Port
	addr       string
	compressor c.Compressor
	err        error
	submit     SubmitFunc
	fileMap    map[string]uint8
	packetMap  []string
	namespace  string
}

// getAddr returns the address of the rpcx server
func getAddr(ip string, port string) string {
	return fmt.Sprintf("%s:%s", ip, port)
}

// getXClient returns the rpcx client instance
// NewPeer2PeerDiscovery function always returns nil to err
// For this reason, second return parameter doesn't have any meaning
func getXClient(addr string, funcName string) (rpcx.XClient, error) {
	discovery, _ := rpcx.NewPeer2PeerDiscovery("tcp@"+addr, "")
	xclient := rpcx.NewXClient(funcName, rpcx.Failtry, rpcx.RandomSelect, discovery, rpcx.DefaultOption)
	return xclient, nil
}

// exceptionHandler does mapping the error code to the Transport structure's error member
func exceptionHandler(transport *Transport, err error) error {
	if transport != nil && err != nil {
		transport.err = err
		transport.Close()
	}
	return err
}

// InitTransport returns the instance of the Transport structure
func InitTransport(configFileName string, customFilterFunc s.CustomFilterFunc) (*Transport, error) {
	var (
		err       error
		scheduler *s.Scheduler
		config    s.Config
		files     []s.File
	)
	t := new(Transport)
	submitOps := s.SubmitOperations{}
	submitOps.Hot = t.hotSubmitFunc
	submitOps.Cold = t.coldSubmitFunc

	scheduler, err = s.InitScheduler(configFileName, submitOps, customFilterFunc)
	if err != nil {
		goto exception
	}
	t.err = nil
	t.scheduler = scheduler
	t.cold.meta.threshold = scheduler.GetConfig().ColdSendThreshold
	t.cold.meta.timeout = time.Duration(scheduler.GetConfig().ColdTimeout) * time.Millisecond
	t.cold.meta.start = time.Now()
	t.compressor = &c.GzipComp{}
	t.submit = Submit
	t.fileMap = make(map[string]uint8)

	files = t.scheduler.GetConfig().Files
	t.packetMap = make([]string, len(files))

	for idx, file := range files {
		t.fileMap[file.Filename] = uint8(idx)
		t.packetMap[idx] = file.Filename
	}

	config = t.scheduler.GetConfig()
	t.addr = getAddr(config.TargetIP, config.TargetPort)

	t.hot.xclient, _ = getXClient(t.addr, "HotPort")
	t.cold.xclient, _ = getXClient(t.addr, "ColdPort")

	return t, err
exception:
	return nil, exceptionHandler(t, err)
}

// Run executes the scheduler
func (t *Transport) Run() error {
	var err error
	if t.scheduler == nil {
		err = errors.New("scheduler must be allocated")
		goto out
	}
	err = t.scheduler.Run()
out:
	if err != nil {
		t.err = err
	}
	return t.err
}

func getInfo(message s.Message) rpc.LogInfo {
	info := rpc.LogInfo{}
	info.Length = message.Info.Length
	info.Timestamp = message.Info.Timestamp
	return info
}

// getPacket converts the custom Message structures to the rpc.LogMessage format
func getPacket(messages []s.Message, fileMap map[string]uint8, packetMap []string) (rpc.LogMessage, error) {
	packet := rpc.LogMessage{}
	packet.Files.MapTable = packetMap

	size := uint64(0)
	for _, message := range messages {
		info := getInfo(message)
		packet.Info = append(packet.Info, info)
		packet.Buffer = append(packet.Buffer, message.Data...)
		if len(message.Info.Filename) == 0 {
			return rpc.LogMessage{}, fmt.Errorf("filename must be specified")
		}
		idx := fileMap[message.Info.Filename]
		packet.Files.Indexes = append(packet.Files.Indexes, idx)
		size += message.Info.Length
	}

	if size != uint64(len(packet.Buffer)) {
		return rpc.LogMessage{}, fmt.Errorf("buffer and info size mismatch (buffer: %d, info: %d)", len(packet.Buffer), size)
	}
	return packet, nil
}

// PrintPacket prints the information in the packet
func PrintPacket(packet rpc.LogMessage, prefix string, isCompressed bool, compressor c.Compressor) {
	if isCompressed {
		buffer, err := compressor.Decompress(packet.Buffer)
		if err != nil {
			log.Panic("Packet decompress failed")
		}
		packet.Buffer = buffer
	}
	wp := uint64(0)
	for i, info := range packet.Info {
		filename := packet.Files.MapTable[packet.Files.Indexes[i]]
		log.Printf("[%s:%s] %s (size: %d)\n", prefix, filename, string(packet.Buffer[wp:wp+info.Length]), len(packet.Buffer))
		wp += info.Length
	}
}

// Submit submits the packet to server by using rpcx
func Submit(packet *rpc.LogMessage, xclient rpcx.XClient) error {
	reply := &rpc.Reply{}
	err := xclient.Call(context.Background(), "Push", packet, reply)
	if err != nil {
		return err
	}
	return nil
}

// hotSubmitFunc submits the hot messages
func (t *Transport) hotSubmitFunc(messages []s.Message) error {
	var (
		packet rpc.LogMessage
		err    error
	)

	packet, err = getPacket(messages, t.fileMap, t.packetMap)
	if err != nil {
		goto exception
	}
	if len(packet.Info) == 0 {
		return nil
	}

	compactPacketMap(&packet)
	packet.Namespace = t.namespace
	err = t.submit(&packet, t.hot.xclient)
	if err != nil {
		goto exception
	}
	return nil
exception:
	return exceptionHandler(t, err)
}

// compactPacketMap compacts the packet mapping table
func compactPacketMap(packet *rpc.LogMessage) {
	seq := uint8(0)
	var indexes []uint8
	mapping := make(map[string]uint8)
	for _, fileIdx := range packet.Files.Indexes {
		if idx, ok := mapping[packet.Files.MapTable[fileIdx]]; !ok {
			mapping[packet.Files.MapTable[fileIdx]] = seq
			indexes = append(indexes, seq)
			seq++
		} else {
			indexes = append(indexes, idx)
		}
	}
	packet.Files.Indexes = indexes
	packet.Files.MapTable = make([]string, len(mapping))
	for k, v := range mapping {
		packet.Files.MapTable[v] = k
	}
}

// coldSubmitFunc submits the cold messages
func (t *Transport) coldSubmitFunc(messages []s.Message) error {
	var (
		meta   *BufferingMetadata
		err    error
		packet rpc.LogMessage
	)

	packet, err = getPacket(messages, t.fileMap, nil)
	if err != nil {
		goto exception
	}
	if len(packet.Info) == 0 {
		return nil
	}
	meta = &t.cold.meta
	meta.packet.Info = append(meta.packet.Info, packet.Info...)
	meta.packet.Buffer = append(meta.packet.Buffer, packet.Buffer...)
	meta.packet.Files.Indexes = append(meta.packet.Files.Indexes, packet.Files.Indexes...)
	if uint64(len(meta.packet.Buffer)) >= meta.threshold || time.Since(meta.start) >= meta.timeout {
		meta.packet.Files.MapTable = t.packetMap
		meta.packet.Buffer, err = t.compressor.Compress(meta.packet.Buffer)
		if err != nil {
			goto exception
		}
		compactPacketMap(&meta.packet)
		meta.packet.Namespace = t.namespace
		err = t.submit(&meta.packet, t.cold.xclient)
		if err != nil {
			goto exception
		}
		meta.packet = rpc.LogMessage{}
	}
	return nil
exception:
	return exceptionHandler(t, err)
}

// Close closes the transport data structure
func (t *Transport) Close() {
	if t.scheduler != nil {
		t.scheduler.Close()
		t.scheduler = nil
	}

	t.cold.Close()
	t.hot.Close()
}
