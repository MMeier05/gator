package main

import (
	"errors"
	"fmt"
	"github.com/MMeier05/blog_aggregator/internal/database"
	"github.com/google/uuid"
	"time"
	"context"
)

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Print(feed)
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("usage: <name> <url>")
	}

	params := database.CreateFeedParams{
		ID: uuid.New(), 		
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: cmd.args[0],
		Url: cmd.args[1],
		UserID: user.ID,
	}
	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}
	paramsFeedFollow := database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	}
	_, err = s.db.CreateFeedFollow(context.Background(), paramsFeedFollow)
	if err != nil {
		return err
	}

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetUserFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Printf("* %s (%s)\n", feed.FeedName, feed.UserName.String)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return errors.New("usage: <url>")
	}
	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	params := database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	}
	inserted_feed_follow, err:= s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}
	fmt.Printf("Created feed: %s (%s)\n", inserted_feed_follow.FeedName, inserted_feed_follow.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command) error {
	followedFeeds, err := s.db.GetFeedFollowsForUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return err
	}
	fmt.Printf("%s's Feeds:\n", s.config.CurrentUserName)
	for _, followedFeed := range followedFeeds {
		fmt.Printf("* %s\n", followedFeed.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args[0]) < 1 {
		return errors.New("usage: <url>")
	}
	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}
	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
			UserID: user.ID,
			FeedID: feed.ID,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
