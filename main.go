package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/brotherlogic/kremind/db"
	"github.com/brotherlogic/kremind/server"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	pb "github.com/brotherlogic/kremind/proto"
)

var (
	port        = flag.Int("port", 8080, "Server port for grpc traffic")
	metricsPort = flag.Int("metrics_port", 8081, "Metrics port")
)

func main() {
	flag.Parse()

	ctx := context.Background()

	s := server.NewServer(db.GetDB())

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("kremind is actually unable to listen on the grpc port %v: %v", *port, err)
	}
	gs := grpc.NewServer()
	pb.RegisterKremindServiceServer(gs, s)

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%v", *metricsPort), nil)
		log.Fatalf("kremind is unable to serve metrics: %v", err)
	}()

	err = gs.Serve(lis)
	log.Printf("kfremind is unable to serve http: %v", err)
}
