package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"

	dictionary "github.com/silasbue/implementation/gRPC"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	file, _ := openLogFile("./logs/server.log")

	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)
	log.SetFlags(2 | 3)

	if len(os.Args) != 2 {
		log.Printf("Please input 0 or 1 to run the server on. Fx. inputting 0 will run the primary RM on port 3000 and inputting 1 will run the backup RM on port 3001")
		return
	}

	ownId := os.Args[1]

	listen, _ := net.Listen("tcp", "localhost:300"+ownId)

	convOwnId, _ := strconv.ParseInt(ownId, 10, 32)

	port := int32(3000)

	if ownId == "0" {
		port++
	}

	fmt.Printf("Trying to dial: %v\n", port)
	conn, err := grpc.Dial(fmt.Sprintf(":%v", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("server %v: Could not connect: %s", ownId, err)
	}
	defer conn.Close()

	isLeader := false
	if ownId == "0" {
		isLeader = true
	}
	grpcServer := grpc.NewServer()
	dictionary.RegisterDictionaryServer(grpcServer, &Server{
		id:          int32(convOwnId),
		values:      make(map[string]string),
		otherServer: dictionary.NewDictionaryClient(conn),
		isLeader:    isLeader,
	})

	log.Printf("server listening at %v", listen.Addr())

	grpcServer.Serve(listen)
}

func (s *Server) Add(ctx context.Context, req *dictionary.AddRequest) (*dictionary.AddReply, error) {
	s.values[req.Word] = req.Definition
	log.Printf("RM %v: recieved an Add request for word: %v. definition: %v", s.id, req.GetWord(), req.GetDefinition())

	if s.isLeader {
		log.Printf("RM %v: sending Add request to backup RM", s.id)
		_, err := s.otherServer.Add(ctx, &dictionary.AddRequest{Word: req.GetWord(), Definition: req.GetDefinition()})

		if err != nil {
			log.Printf("Could not copy data to backup RM: ERROR - %v", err)
		} else {
			log.Printf("Backup RM updated!")
		}
	}

	return &dictionary.AddReply{Success: true}, nil
}

func (s *Server) Read(ctx context.Context, req *dictionary.ReadRequest) (*dictionary.ReadReply, error) {
	definition := s.values[req.Word]
	log.Printf("RM %v: value was read. Word recieved: %v, Definition returned: %v", s.id, req.GetWord(), definition)
	return &dictionary.ReadReply{Definition: definition}, nil
}

func openLogFile(path string) (*os.File, error) {
	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return logFile, nil
}

type Server struct {
	dictionary.UnimplementedDictionaryServer
	id          int32
	values      map[string]string
	isLeader    bool
	otherServer dictionary.DictionaryClient
}
