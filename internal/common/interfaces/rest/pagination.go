package commonrest

import "github.com/google/uuid"

type (
	PaginationQuery struct {
		Limit       int       `query:"limit" minimum:"1" maximum:"100" default:"20"`
		AfterCursor uuid.UUID `query:"after"`
	}

	Page[T any] struct {
		Items      []T        `json:"items"`
		NextCursor *uuid.UUID `json:"next,omitempty"`
	}
)
