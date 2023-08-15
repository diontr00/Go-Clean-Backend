package controller

import "fmt"

type RssRequestError struct {
	err error
}

func (r RssRequestError) Error() string {
	return fmt.Sprintf("Rss Request Error : %v", r.err)
}

type RssParsingError struct {
	err error
}

func (r RssParsingError) Error() string {
	return fmt.Sprintf("Rss Parsing Error : %v", r.err)
}
