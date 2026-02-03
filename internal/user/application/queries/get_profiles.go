package queries

import (
	"context"
	commonapplication "goddd/internal/common/application"
	"goddd/internal/user/domain"
	usersql "goddd/internal/user/infrastructure/sql"
	"log/slog"

	"github.com/google/uuid"
)

type GetProfilesInput struct {
	Limit int
	After *uuid.UUID
	Name  *string
}

type GetProfilesOutput struct {
	Profiles []domain.User
	Next     *uuid.UUID
}

type GetProfilesQuery commonapplication.QueryI[GetProfilesInput, GetProfilesOutput]

func NewGetProfilesQuery(
	log *slog.Logger,
	userSql *usersql.Queries,
) GetProfilesQuery {
	return commonapplication.NewQuery(log, &getProfiles{
		userSql: userSql,
	})
}

type getProfiles struct {
	userSql *usersql.Queries
}

func (u *getProfiles) Handle(
	ctx context.Context, log *slog.Logger, input GetProfilesInput,
) (GetProfilesOutput, error) {
	rows, err := u.userSql.ListUsers(
		ctx,
		usersql.ListUsersParams{
			LimitPlusOne: int32(input.Limit + 1),
			After:        input.After,
			Name:         input.Name,
		},
	)
	if err != nil {
		return GetProfilesOutput{}, err
	}

	var next *uuid.UUID
	if len(rows) > input.Limit {
		rows = rows[:input.Limit]
		next = &rows[input.Limit-1].ID
	}

	users := make([]domain.User, len(rows))
	for i, row := range rows {
		users[i] = *domain.RehydrateUser(
			row.ID, int(row.Version), row.CreatedAt, row.UpdatedAt, row.Name,
		)
	}
	return GetProfilesOutput{Profiles: users, Next: next}, err
}
