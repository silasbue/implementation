package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	dictionary "github.com/silasbue/implementation/gRPC"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ownPort, _ := strconv.Atoi(os.Args[1])

	file, _ := openLogFile("./logs/client.log")

	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)
	log.SetFlags(2 | 3)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	servers := make([]dictionary.DictionaryClient, 2)

	for i := 0; i < 2; i++ {
		port := int32(3000) + int32(i)

		fmt.Printf("Trying to dial: %v\n", port)
		conn, err := grpc.Dial(fmt.Sprintf(":%v", port), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Front end %v: Could not connect: %s", ownPort, err)
		}
		servers[i] = dictionary.NewDictionaryClient(conn)
		defer conn.Close()
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := strings.Split(scanner.Text(), " ")
		command[0] = strings.ToLower(command[0])
		command[1] = strings.ToLower(command[1])

		if command[0] == "add" {
			command[2] = strings.ToLower(command[2])

			for _, server := range servers {
				res, err := server.Add(ctx, &dictionary.AddRequest{Word: command[1], Definition: command[2]})
				if err != nil {
					continue
				}
				log.Printf("Add succeeded: %v", res.GetSuccess())
				break
			}
		}

		if command[0] == "read" {

			for _, server := range servers {
				res, err := server.Read(ctx, &dictionary.ReadRequest{Word: command[1]})
				if err != nil {
					continue
				}
				log.Printf("Definition: %v", res.GetDefinition())
				break
			}
		}
	}
}

func openLogFile(path string) (*os.File, error) {
	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return logFile, nil
}
