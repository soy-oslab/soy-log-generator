package server

import (
	"context"
	"fmt"
	"log"

	"github.com/kr/pretty"
	"github.com/smallnest/rpcx/server"
	"github.com/soyoslab/soy_log_collector/pkg/rpc"
	"github.com/soyoslab/soy_log_generator/pkg/compressor"
)

type hotPort int
type coldPort int
type initPort int

var mapTable map[string][]string

// printPacket prints the information in the packet
func printPacket(packet rpc.LogMessage, prefix string, isCompressed bool, compressor compressor.Compressor) {
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

func (p *hotPort) Push(ctx context.Context, args *rpc.LogMessage, reply *rpc.Reply) error {
	(*args).Files.MapTable = mapTable[(*args).Namespace]
	log.Printf("hot: %# v", pretty.Formatter(*args))
	printPacket(*args, "HOT", false, &compressor.GzipComp{})
	return nil
}

func (p *coldPort) Push(ctx context.Context, args *rpc.LogMessage, reply *rpc.Reply) error {
	(*args).Files.MapTable = mapTable[(*args).Namespace]
	log.Printf("cold: %# v", pretty.Formatter(*args))
	printPacket(*args, fmt.Sprintf("COLD(%v)", len((*args).Buffer)), true, &compressor.GzipComp{})
	return nil
}

func (p *initPort) Push(ctx context.Context, args *rpc.LogMessage, reply *rpc.Reply) error {
	mapTable[(*args).Namespace] = (*args).Files.MapTable
	log.Println("mapping table: %# v", pretty.Formatter((*args).Files.MapTable))
	return nil
}

// Run runs the testing server program
func Run() {
	mapTable = make(map[string][]string)
	s := server.NewServer()
	s.RegisterName("HotPort", new(hotPort), "")
	s.RegisterName("ColdPort", new(coldPort), "")
	s.RegisterName("Init", new(initPort), "")
	// Don't change the address of `localhost:8972`
	// Because this program uses in the `pkg/transport/transport_test.go`
	s.Serve("tcp", "localhost:8972")
}
