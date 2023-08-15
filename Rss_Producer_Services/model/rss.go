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
	Count int    `json:"insertCount,omitempty"`
	Error string `json:"error,omitempty"`
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

//go:generate mockery --name RssUsecase
type RssUsecase interface {
	FetchAndInsert(c context.Context, request *RssRequest) error
}
