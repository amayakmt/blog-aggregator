package main

import (
	"context"
	"fmt"
	"time"

	"github.com/amayakmt/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", feed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Arguments) < 2 {
		return fmt.Errorf("usage: addfeed <name> <url>")
	}

	user, err := s.DB.GetUser(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("could not get current user: %w", err)
	}

	now := time.Now()

	feed, err := s.DB.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      cmd.Arguments[0],
		Url:       cmd.Arguments[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("could not create feed: %w", err)
	}

	// automatically follow the feed just created
	_, err = s.DB.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not follow feed: %w", err)
	}

	fmt.Printf("ID:         %v\n", feed.ID)
	fmt.Printf("CreatedAt:  %v\n", feed.CreatedAt)
	fmt.Printf("UpdatedAt:  %v\n", feed.UpdatedAt)
	fmt.Printf("Name:       %v\n", feed.Name)
	fmt.Printf("URL:        %v\n", feed.Url)
	fmt.Printf("UserID:     %v\n", feed.UserID)
	return nil
}

func handlerGetFeeds(s *state, cmd command) error {
	feeds, err := s.DB.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("Name:  %v\n", feed.Name)
		fmt.Printf("URL:   %v\n", feed.Url)
		fmt.Printf("User:  %v\n", feed.UserName)
		fmt.Println("---")
	}
	return nil
}

func handlerFollow(s *state, cmd command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("usage: follow <url>")
	}

	user, err := s.DB.GetUser(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("could not get current user: %w", err)
	}

	feed, err := s.DB.GetFeedByURL(context.Background(), cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("could not find feed: %w", err)
	}

	now := time.Now()

	follow, err := s.DB.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not follow feed: %w", err)
	}

	fmt.Printf("Feed:  %v\n", follow.FeedName)
	fmt.Printf("User:  %v\n", follow.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command) error {
	user, err := s.DB.GetUser(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("could not get current user: %w", err)
	}

	follows, err := s.DB.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, f := range follows {
		fmt.Printf("* %v\n", f.FeedName)
	}
	return nil
}
