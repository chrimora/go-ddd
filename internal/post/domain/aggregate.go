package domain

import (
	"fmt"
	commondomain "goddd/internal/common/domain"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	commondomain.AggregateRoot
	title       string
	publishDate time.Time
	author      string
}

func NewPost(title, author string) *Post {
	post := &Post{
		AggregateRoot: commondomain.NewAggregateRoot(),
		title:         title,
		publishDate:   time.Now().UTC(),
		author:        author,
	}
	post.AddEvent(NewPostCreatedEvent(post.ID()))
	return post
}

func RehydratePost(
	id uuid.UUID, version int, createdAt, updatedAt time.Time,
	title string, publishDate time.Time, author string,
) *Post {
	return &Post{
		AggregateRoot: commondomain.RehydrateAggregateRoot(id, version, createdAt, updatedAt),
		title:         title,
		publishDate:   publishDate,
		author:        author,
	}
}

func (p *Post) Title() string          { return p.title }
func (p *Post) PublishDate() time.Time { return p.publishDate }
func (p *Post) Author() string         { return p.author }
func (p *Post) String() string         { return fmt.Sprintf("Post[id: %s]", p.ID()) }

func (p *Post) UpdateTitle(title string) {
	p.title = title
	p.AggregateRoot.Update()
	p.AddEvent(NewPostUpdatedEvent(p.ID()))
}

func (p *Post) Clone() *Post {
	return &Post{
		AggregateRoot: p.AggregateRoot.Clone(),
		title:         p.title,
		publishDate:   p.publishDate,
		author:        p.author,
	}
}
