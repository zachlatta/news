package hn

import (
	"net/url"
	"time"
)

type PostType int

const (
	Standard PostType = iota
	Question
)

// Post represents a post from HN
type Post struct {
	ID       int64
	Title    string
	Author   User
	URL      url.URL
	Body     string // Body and URL are mutually exclusive
	Points   int
	Comments []Comment
}

// Type returns the type of the post. The post is a Question if it has a body,
// else it's a standard post.
func (p *Post) Type() PostType {
	if p.Body == "" {
		return Question
	}
	return Standard
}

// Comment represents a comment from HN
type Comment struct {
	ID     int64
	Author User
	Body   string
}

// User represents a user on HN
type User struct {
	Username     string
	Created      time.Time
	Karma        int
	AverageKarma float32
	About        string

	Submissions []Post
	Comments    []Comment
}
