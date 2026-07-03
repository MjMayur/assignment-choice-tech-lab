package service

import (
	"context"
	"mime/multipart"

	"project/entity"
	"project/pkg/errors"
)

type AssignmentUserService interface {
	StoreCsvRecords(ctx context.Context, fileHeaders *multipart.FileHeader) (map[string]interface{}, map[string]interface{}, string, errors.Response)
	CreateAssignmentUser(ctx context.Context, assignmentUser entity.AssignmentUser) (int, errors.Response)
	PartialUpdateAssignmentUser(ctx context.Context, assignmentUser entity.AssignmentUser) errors.Response
	DeleteAssignmentUser(ctx context.Context, assignmentUserID int) errors.Response
	GetAssignmentUserbyID(ctx context.Context, assignmentUserID int) (entity.AssignmentUser, errors.Response)
	ListAssignmentUser(ctx context.Context) (int, []entity.AssignmentUser, errors.Response)
}

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
