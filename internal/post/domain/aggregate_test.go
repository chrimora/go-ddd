//go:build unit

package domain_test

import (
	"goddd/internal/post/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewPost(t *testing.T) {
	post := domain.NewPost("Hello World", "Alice")

	assert.NotEmpty(t, post.ID())
	assert.Equal(t, "Hello World", post.Title())
	assert.Equal(t, "Alice", post.Author())
	assert.WithinDuration(t, time.Now(), post.PublishDate(), time.Second)

	events := post.PullEvents()
	assert.Len(t, events, 1)

	event, ok := events[0].(domain.PostCreatedEvent)
	assert.True(t, ok)
	assert.Equal(t, post.ID(), event.GetAggregateId())
}

func TestUpdateTitle(t *testing.T) {
	post := domain.NewPost("Original Title", "Alice")
	post.PullEvents() // clear creation event
	oldUpdatedAt := post.UpdatedAt()

	time.Sleep(1 * time.Millisecond)
	post.UpdateTitle("Updated Title")

	assert.Equal(t, "Updated Title", post.Title())
	assert.Greater(t, post.UpdatedAt(), oldUpdatedAt)

	events := post.PullEvents()
	assert.Len(t, events, 1)

	event, ok := events[0].(domain.PostUpdatedEvent)
	assert.True(t, ok)
	assert.Equal(t, post.ID(), event.GetAggregateId())
}
