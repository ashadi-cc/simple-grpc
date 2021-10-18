package main

import (
	"log"
	"net"
	"protobuf/model/user"
	"protobuf/services"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	user.RegisterUserServiceServer(server, services.NewUserService())

	log.Println("listening on port 5001")
	if err := server.Serve(lis); err != nil {
		log.Fatal(err.Error())
	}
}
