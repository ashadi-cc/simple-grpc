package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"protobuf/model/user"
	"time"

	"google.golang.org/grpc"
)

const userServiceAddress = "localhost:5001"

func main() {
	conn, err := grpc.Dial(userServiceAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	userServiceClient := user.NewUserServiceClient(conn)

	for i := 0; i < 5; i++ {

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		payload := &user.GreetingRequest{
			Name:       fmt.Sprintf("test %d", i),
			Salutation: "do",
		}

		response, err := userServiceClient.GreetUser(ctx, payload)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(response.GreetingMessage)

	}

	in := &user.StreamRequest{Id: 1}
	stream, err := userServiceClient.FetchData(context.Background(), in)
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			if err != nil {
				log.Fatal(err)
			}

			log.Println("Resp stream", resp.Result)
		}
	}()

	<-done
}
