package main

import (
	"context"
	"fmt"
	"log"
	"logger/data"
	"logger/proto/pb"
	"net"

	"google.golang.org/grpc"
) 

type LogServer struct {
	pb.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *pb.LogRequest) (*pb.LogResponse, error) {
	input := req.GetLogEntry()

	// write the log
	logEntry := data.LogEntry {
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &pb.LogResponse{Result: "failed"}
		return res, err
	}

	// return resp
	res := &pb.LogResponse{Result: "logged"}
	return res, nil

}

func (app *Config) gRPCListen() {
	log.Println("Co listen grpc")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
	server := grpc.NewServer()

	pb.RegisterLogServiceServer(server, &LogServer{Models: app.Models})

	log.Printf("gRPC Server starting on port %s", gRpcPort)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for gRPC: %v",err )
	}
}