package services

import (
	"context"
	"fmt"
	"log"
	"protobuf/model/user"
	"sync"
	"time"
)

func NewUserService() user.UserServiceServer {
	return &userService{}
}

type userService struct{}

func (u *userService) GreetUser(ctx context.Context, req *user.GreetingRequest) (*user.GreetingResponse, error) {
	message := fmt.Sprintf("Hello %s, %s", req.Name, req.Salutation)

	return &user.GreetingResponse{GreetingMessage: message}, nil
}

func (u *userService) FetchData(req *user.StreamRequest, srv user.UserService_FetchDataServer) error {
	log.Println("fetch response for id:", req.Id)

	var wg sync.WaitGroup
	for i := 0; i <= 5; i++ {
		wg.Add(1)
		go func(count int) {
			defer wg.Done()
			time.Sleep(time.Duration(count) * time.Second)
			resp := user.StreamResponse{Result: fmt.Sprintf("request #%d for id:%d", count, req.Id)}
			if err := srv.Send(&resp); err != nil {
				log.Println(err.Error())
			}
			log.Printf("finishing request number : %d", count)
		}(i)
	}

	wg.Wait()

	return nil
}
