package transport

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	rpcx "github.com/smallnest/rpcx/client"
	"github.com/soyoslab/soy_log_collector/pkg/rpc"
	c "github.com/soyoslab/soy_log_generator/pkg/compressor"
	s "github.com/soyoslab/soy_log_generator/pkg/scheduler"
)

const ConfigTest = `{
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

func TestPortClose(t *testing.T) {
	addr := getAddr("localhost", "8972")
	xclient, err := getXClient(addr, "HotPort")
	if err != nil {
		t.Errorf("port close failed: %v", err)
	}
	port := Port{xclient, BufferingMetadata{}}
	port.Close()
}

func (t *Transport) getError() error {
	return t.err
}

func TestExceptionHandler(t *testing.T) {
	sample := &Transport{}
	err := errors.New("Test")
	_ = exceptionHandler(sample, err)
	if sample.getError() != err {
		t.Errorf("exception handler doesn't work correctly (%s<->%s)", sample.getError(), err)
	}
}

func setup() (string, string) {
	testFile, err := os.CreateTemp("", "transport-init-test")
	if err != nil {
		log.Fatalf("config file open failed: %v", err)
	}
	configFile, err := os.OpenFile(testFile.Name(), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.FileMode(0644))
	if err != nil {
		log.Fatalf("config file open failed: %v", err)
	}
	defer configFile.Close()
	configText := fmt.Sprintf(ConfigTest, testFile.Name())
	writer := bufio.NewWriter(configFile)
	writer.WriteString(configText)
	writer.Flush()
	return testFile.Name(), configFile.Name()
}

func TestInit(t *testing.T) {
	testFileName, configFileName := setup()
	defer func() { os.Remove(testFileName) }()
	trans, err := InitTransport(configFileName, nil)
	if err != nil {
		t.Errorf("initialize the transport failed %v", err)
	}
	trans.Close()
}

func TestInitFailed(t *testing.T) {
	_, err := InitTransport("", nil)
	if err == nil {
		t.Errorf("invalid initialize parameter but it works")
	}
}

func TestNilRun(t *testing.T) {
	sample := &Transport{}
	err := sample.Run()
	if err == nil {
		t.Errorf("invalid run state but it runs")
	}
}

func messageGeneration(data string, length uint64, isCompressed bool) []s.Message {
	var messages []s.Message
	val := s.Message{}
	val.Data = []byte(data)
	val.Info.Length = length
	val.Info.Filename = "test1.txt"
	messages = append(messages, val)
	if isCompressed {
		compressor := &c.GzipComp{}
		messages[0].Data, _ = compressor.Compress(messages[0].Data)
	}
	return messages
}

func TestGetPacket(t *testing.T) {
	messages := messageGeneration("test", 4, false)
	fileMap := make(map[string]uint8)
	fileMap["test"] = 0
	packetMap := []string{"test"}
	_, err := getPacket(messages, fileMap, packetMap)
	if err != nil {
		t.Errorf("get packet method doesn't work correctly: %v", err)
	}
}

func TestGetPacketInvalid(t *testing.T) {
	messages := messageGeneration("test", 0, false)
	fileMap := make(map[string]uint8)
	fileMap["test"] = 0
	packetMap := []string{"test"}
	_, err := getPacket(messages, fileMap, packetMap)
	if err == nil {
		t.Errorf("get packet receives invalid data but it works")
	}
}

func TestPrintPacket(t *testing.T) {
	messages := messageGeneration("test", 4, false)
	fileMap := make(map[string]uint8)
	fileMap["test"] = 0
	packetMap := []string{"test"}
	packet, err := getPacket(messages, fileMap, packetMap)
	if err != nil {
		t.Errorf("get packet failed %v", err)
	}
	PrintPacket(packet, "test", false, nil)
}

func TestPrintCompressedPacket(t *testing.T) {
	messages := messageGeneration("test", 4, false)
	fileMap := make(map[string]uint8)
	fileMap["test"] = 0
	packetMap := []string{"test"}
	packet, err := getPacket(messages, fileMap, packetMap)
	if err != nil {
		t.Errorf("get packet failed %v", err)
	}
	compressor := &c.GzipComp{}
	packet.Buffer, _ = compressor.Compress(packet.Buffer)
	PrintPacket(packet, "test", true, compressor)
}

func TestInvalidSubmit(t *testing.T) {
	messages := messageGeneration("test", 4, false)
	fileMap := make(map[string]uint8)
	fileMap["test"] = 0
	packetMap := []string{"test"}
	packet, _ := getPacket(messages, fileMap, packetMap)
	xclient, _ := getXClient("localhost:8972", "HotPort")
	err := Submit(&packet, xclient)
	if err == nil {
		t.Errorf("invalid submit requested but it works")
	}
}

func validSubmitFunc(t *testing.T, trans *Transport, f func([]s.Message) error) {
	messages := messageGeneration("test", 4, false)
	trans.submit = func(msg *rpc.LogMessage, _ rpcx.XClient) error {
		if msg == nil {
			t.Errorf("invalid hot submit detected")
		}
		return nil
	}
	err := f(messages)
	if err != nil {
		t.Errorf("hot data submission failed: %v", err)
	}
}

func TestHotSubmitFunc(t *testing.T) {
	trans := Transport{}
	trans.fileMap = make(map[string]uint8)
	trans.fileMap["test"] = 0
	trans.namespace = "test"
	trans.packetMap = []string{"test"}
	validSubmitFunc(t, &trans, trans.hotSubmitFunc)
}

func TestColdSubmitFunc(t *testing.T) {
	trans := Transport{}
	trans.fileMap = make(map[string]uint8)
	trans.fileMap["test"] = 0
	trans.packetMap = []string{"test"}
	trans.namespace = "test"
	trans.compressor = &c.GzipComp{}
	validSubmitFunc(t, &trans, trans.coldSubmitFunc)
}

func invalidSubmitFunc(t *testing.T, target string, isCompressed bool, trans *Transport, f func([]s.Message) error) {
	trans.submit = func(msg *rpc.LogMessage, _ rpcx.XClient) error {
		if msg == nil {
			t.Errorf("invalid hot submit detected")
		}
		return nil
	}
	messages := messageGeneration("", 0, isCompressed)
	err := f(messages)
	if err != nil {
		t.Errorf("%s info length is zero valid returns but it evaluates invalid: %v", target, err)
	}

	messages = messageGeneration("test", 3, isCompressed)
	err = f(messages)
	if err == nil {
		t.Errorf("%s invalid message info but it evaluates valid", target)
	}

	trans.submit = func(msg *rpc.LogMessage, _ rpcx.XClient) error {
		return errors.New("test")
	}
	messages = messageGeneration("test", 4, isCompressed)
	err = f(messages)
	if err == nil && target != "cold" {
		t.Errorf("%s invalid submit returns but it evaluates valid", target)
	}
}

func TestHotInvalidSubmit(t *testing.T) {
	trans := Transport{}
	trans.fileMap = make(map[string]uint8)
	trans.fileMap["test"] = 0
	trans.packetMap = []string{"test"}
	trans.namespace = "test"
	invalidSubmitFunc(t, "hot", false, &trans, trans.hotSubmitFunc)
}

func TestColdInvalidSubmit(t *testing.T) {
	testFileName, configFileName := setup()
	defer func() { os.Remove(testFileName) }()
	trans, _ := InitTransport(configFileName, nil)
	invalidSubmitFunc(t, "cold", false, trans, trans.coldSubmitFunc)
	trans.Close()
}
