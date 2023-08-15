package model

import (
	"context"
)

type RSSFeed struct {
	Entries []*RSSEntry `xml:"entry"`
}

type RssRequest struct {
	URL string `json:"url" validate:"required,url"`
}

type RssResponse struct {
	InsertCount int    `json:"insertCount,omitempty"`
	Error       string `json:"error,omitempty"`
}

type RSSEntry struct {
	Link struct {
		Href string `xml:"href,attr"`
	} `xml:"link"`
	Thumbnail struct {
		URL string `xml:"url,attr"`
	} `xml:"thumbnail"`
	Title string `xml:"title"`
}

//go:generate mockery --name RssRepository
type RssRepository interface {
	InsertEntries(c context.Context, entries []*RSSEntry) (int, error)
}

//go:generate mockery --name RssUsecase
type RssController interface {
	InsertEntries(c context.Context, entries []*RSSEntry) (int, error)
	GetFeedEntries(c context.Context, url string) ([]*RSSEntry, error)
	ListenForRssRequest(c context.Context, done chan struct{})
}
