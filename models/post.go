package models

import (
	"html/template"
	"time"
)

type Post struct {
	Id          int64         `json:"id"`
	Title       string        `json:"title"`
	Slug        string        `json:"slug"`
	ParentId    int64         `json:"parent_id"`
	Content     template.HTML `json:"content"`
	Raw         string        `json:"raw"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	PublishedAt time.Time     `json:"published_at"`
	HeadMatter  HeadMatter    `json:"head_matter"`
	Filename    string        `json:"filename"`
	Directory   string        `json:"directory"`
	Type        string        `json:"type"`
}
