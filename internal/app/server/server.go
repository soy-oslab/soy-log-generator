package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kr/pretty"
	"github.com/smallnest/rpcx/server"
	"github.com/soyoslab/soy_log_collector/pkg/rpc"
	"github.com/soyoslab/soy_log_generator/pkg/compressor"
	"github.com/soyoslab/soy_log_generator/pkg/transport"
)

type hotPort int
type coldPort int

func (p *hotPort) Push(ctx context.Context, args *rpc.LogMessage, reply *rpc.Reply) error {
	log.Printf("%# v", pretty.Formatter(*args))
	transport.PrintPacket(*args, "HOT", false, &compressor.GzipComp{})
	return nil
}

func (p *coldPort) Push(ctx context.Context, args *rpc.LogMessage, reply *rpc.Reply) error {
	log.Printf("%# v", pretty.Formatter(*args))
	transport.PrintPacket(*args, fmt.Sprintf("COLD(%v)", len((*args).Buffer)), true, &compressor.GzipComp{})
	return nil
}

func main() {
	s := server.NewServer()
	s.RegisterName("HotPort", new(hotPort), "")
	s.RegisterName("ColdPort", new(coldPort), "")
	s.Serve("tcp", "localhost:8972")
}
