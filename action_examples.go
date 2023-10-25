package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

// SimpleAction is an action that runs a function
type SimpleAction struct {
	Func func(ctx context.Context) error
}

// Run runs the function
func (sa *SimpleAction) Run(ctx context.Context) error {
	return sa.Func(ctx)
}

type User struct {
	ID    int
	Name  string
	Email string
	Error error
}

type (
	Users         []User
	UsersParallel []User
)

func getUsersParalel() UsersParallel {
	return []User(getUsers())
}

func getUsers() Users {
	return []User{
		{
			ID:    1,
			Name:  "John Doe",
			Email: "john.doe@gmail.com",
		},
		{
			ID:    2,
			Name:  "Jane Smith",
			Email: "jane.smith@yahoo.com",
		},
		{
			ID:    3,
			Name:  "Alex Brown",
			Email: "alex.brown@hotmail.com",
			// Error: errors.New("error user 3"),
		},
		{
			ID:    4,
			Name:  "Emily Johnson",
			Email: "emily.johnson@outlook.com",
		},
		{
			ID:    5,
			Name:  "Michael Davis",
			Email: "michael.davis@gmail.com",
		},
	}
}

func (u *User) Run(ctx context.Context) error {
	fmt.Println(u)
	return nil
}

func (u Users) Run(ctx context.Context) error {
	for _, user := range u {
		fmt.Printf("Sequential user: %v\n", user)
	}
	return nil
}

func (u UsersParallel) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, user := range u {
		currUser := user
		g.Go(func() error {
			// Check if context is done
			select {
			case <-ctx.Done():
				return ctx.Err() // Return the error from the context, e.g., context.Canceled
			case <-time.After(time.Second * 1):
				fmt.Printf("Parallel user: %v\n", currUser)
				return currUser.Error
			}
		})
	}
	return g.Wait()
}

type Loop struct {
	Count int
}

func NewLoop(count int) *Loop {
	return &Loop{Count: count}
}

func (l *Loop) Run(ctx context.Context) error {
	for i := 0; i < l.Count; i++ {
		fmt.Printf("Loop: %d\n", i)
	}
	return nil
}
