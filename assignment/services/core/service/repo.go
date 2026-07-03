package service

import (
	"context"

	"project/entity"
	"project/pkg/errors"
)

type LoginRepo interface {
	GetUserDetailsForLogin(ctx context.Context, email string) (string, errors.Response)
}

type HomeRepo interface {
	Home(ctx context.Context) errors.Response
}

type AssignmentUserRepo interface {
	BulkInsertAssignmentUsers(ctx context.Context, users []entity.AssignmentUser) (int, error)
	CreateAssignmentUser(ctx context.Context, assignmentUser entity.AssignmentUser) (int, errors.Response)
	PartialUpdateAssignmentUser(ctx context.Context, assignmentUser entity.AssignmentUser) errors.Response
	DeleteAssignmentUser(ctx context.Context, assignmentUserID int) errors.Response
	GetAssignmentUserbyID(ctx context.Context, assignmentUserID int) (entity.AssignmentUser, errors.Response)
	ListAssignmentUser(ctx context.Context) (int, []entity.AssignmentUser, errors.Response)
}

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
