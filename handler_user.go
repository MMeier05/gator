package main

import (
	"errors"
	"fmt"
	"github.com/MMeier05/gator/internal/database"
	"github.com/google/uuid"
	"time"
	"context"
)


func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Required 'username'\n")  
	}
	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}
	err = s.config.SetUser(user.Name)
	if err != nil {
		return err
	}
	fmt.Printf("User set to '%s' successfully\n", cmd.args[0])
	return nil
}


func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Not enough arguments. Required: 'username'\n")
	}
	if _, err := s.db.GetUser(context.Background(), cmd.args[0]); err == nil {
		return errors.New("User already exists\n")
	}
	parameters := database.CreateUserParams{
		ID : uuid.New(),
		CreatedAt : time.Now(),
		UpdatedAt : time.Now(),
		Name : cmd.args[0],
	}
	user, err := s.db.CreateUser(context.Background(), parameters)
	if err != nil {
		return err
	}
	err = s.config.SetUser(user.Name)
	if err != nil {
		return errors.New("couldn't set current user")
	}
	fmt.Printf("%s\n", user.Name)
	return nil
}


func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		fmt.Print("1\n")
		return err
	}
	fmt.Print("0\n")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, u := range users {
		if u.Name == s.config.CurrentUserName {
			fmt.Printf("* %s (current)\n", u.Name)
		} else {
			fmt.Printf("* %s\n", u.Name)
		}
	}
	return nil
}

