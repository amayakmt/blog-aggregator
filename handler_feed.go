package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/amayakmt/blog-aggregator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("usage: agg <time_between_reqs>")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("invalid duration %q: %w", cmd.Arguments[0], err)
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	feed, err := s.DB.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Printf("error getting next feed: %v\n", err)
		return
	}

	err = s.DB.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		fmt.Printf("error marking feed fetched: %v\n", err)
		return
	}

	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		fmt.Printf("error fetching feed %v: %v\n", feed.Name, err)
		return
	}

	fmt.Printf("Fetching feed: %v\n", feed.Name)
	now := time.Now()

	for _, item := range rssFeed.Channel.Item {
		_, err = s.DB.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  item.Description != "",
			},
			PublishedAt: sql.NullString{
				String: item.PubDate,
				Valid:  item.PubDate != "",
			},
			FeedID: feed.ID,
		})
		if err != nil {
			// unique violation = post already exists, skip silently
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				continue
			}
			fmt.Printf("error saving post %q: %v\n", item.Title, err)
		}
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) < 2 {
		return fmt.Errorf("usage: addfeed <name> <url>")
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

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("usage: follow <url>")
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

func handlerFollowing(s *state, cmd command, user database.User) error {
	follows, err := s.DB.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, f := range follows {
		fmt.Printf("* %v\n", f.FeedName)
	}
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.Arguments) > 0 {
		var err error
		limit, err = strconv.Atoi(cmd.Arguments[0])
		if err != nil {
			return fmt.Errorf("invalid limit %q: %w", cmd.Arguments[0], err)
		}
	}

	posts, err := s.DB.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Printf("Title:  %v\n", post.Title)
		fmt.Printf("URL:    %v\n", post.Url)
		if post.Description.Valid {
			fmt.Printf("Desc:   %v\n", post.Description.String)
		}
		if post.PublishedAt.Valid {
			fmt.Printf("Date:   %v\n", post.PublishedAt.String)
		}
		fmt.Println("---")
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("usage: unfollow <url>")
	}

	feed, err := s.DB.GetFeedByURL(context.Background(), cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("could not find feed: %w", err)
	}

	err = s.DB.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not unfollow feed: %w", err)
	}

	fmt.Printf("Unfollowed %v\n", feed.Name)
	return nil
}
