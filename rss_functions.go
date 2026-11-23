package main

import (
	"context"
	"net/http"
	"io"
	"encoding/xml"
	"time"
	"html"
	"database/sql"
	"strings"
	"github.com/MMeier05/blog_aggregator/internal/database"
	"github.com/google/uuid"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}


func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var feed RSSFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, err
	}
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for _, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}
	return &feed, nil
	
}

func scrapeFeeds(s *state, user database.User) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background(), user.ID)
	if err != nil {
	return err
	}
	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{ 
				ID: nextFeed.ID,
				LastFetchedAt: sql.NullTime{
				Time: time.Now(),
				Valid: true,
			},
		},
	)
	if err != nil {
	return err
	}
	feed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return err
	}
	for _, feedItem := range feed.Channel.Item {
		parsedPubDate, err := time.Parse(time.RFC1123Z, feedItem.PubDate)
		pubDateStruct := sql.NullTime{
			Time: parsedPubDate,
			Valid: true,
		}
		if err != nil {
			pubDateStruct.Valid = false
		}
		
		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title: feedItem.Title,
			Url: feedItem.Link,
			Description: sql.NullString{
				String: feedItem.Description,
				Valid: true,
			},
			PublishedAt: pubDateStruct,
			FeedID: nextFeed.ID, 
		},
			)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value") {
			continue
			}
			return err
		}
	}
	return nil
}
