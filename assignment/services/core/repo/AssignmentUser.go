package repo

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"project/entity"
	"project/internal/converter"
	"project/model"
	falconContext "project/pkg/context"
	"project/pkg/errors"
	"project/pkg/log"

	"github.com/jmoiron/sqlx"
)

type AssignmentUserRepoImpl struct {
	db *sqlx.DB
}

func NewAssignmentUserRepoImpl(db *sqlx.DB) (*AssignmentUserRepoImpl, error) {
	repo := &AssignmentUserRepoImpl{db: db}
	if err := repo.ensureTableExists(context.Background()); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *AssignmentUserRepoImpl) ensureTableExists(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS assignment_users (
		id INT NOT NULL AUTO_INCREMENT,
		first_name VARCHAR(100) DEFAULT '',
		last_name VARCHAR(100) DEFAULT '',
		company_name VARCHAR(150) DEFAULT '',
		address VARCHAR(250) DEFAULT '',
		city VARCHAR(100) DEFAULT '',
		county VARCHAR(100) DEFAULT '',
		postal VARCHAR(50) DEFAULT '',
		phone VARCHAR(50) DEFAULT '',
		email VARCHAR(150) DEFAULT '',
		web VARCHAR(255) DEFAULT '',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP NULL DEFAULT NULL,
		PRIMARY KEY (id)
	)`

	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *AssignmentUserRepoImpl) BulkInsertAssignmentUsers(ctx context.Context, users []entity.AssignmentUser) (int, error) {
	reqID, _ := falconContext.GetRequestIDFromContext(ctx)
	log.Info(fmt.Sprintf("repo: BulkInsertAssignmentUsers started, count: %d", len(users)), reqID)

	if len(users) == 0 {
		return 0, nil
	}

	query := `INSERT INTO assignment_users 
		(first_name, last_name, company_name, address, city, county, postal, phone, email, web) 
		VALUES `

	valueStrings := []string{}
	valueArgs := []interface{}{}

	for _, user := range users {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs,
			user.FirstName,
			user.LastName,
			user.CompanyName,
			user.Address,
			user.City,
			user.County,
			user.Postal,
			user.Phone,
			user.Email,
			user.Web,
		)
	}

	query += strings.Join(valueStrings, ",")

	result, err := r.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		log.Info(fmt.Sprintf("Failed to bulk insert assignment users: %v", err), reqID)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Info(fmt.Sprintf("Failed to get rows affected: %v", err), reqID)
		return 0, err
	}

	log.Info(fmt.Sprintf("repo: BulkInsertAssignmentUsers completed, inserted: %d", rowsAffected), reqID)
	return int(rowsAffected), nil
}

func (r *AssignmentUserRepoImpl) CreateAssignmentUser(ctx context.Context, assignmentUser entity.AssignmentUser) (int, errors.Response) {
	reqID, _ := falconContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>assignmentUser: CreateAssignmentUser started", reqID)

	query := "INSERT INTO assignment_users (first_name, last_name, company_name, address, city, county, postal, phone, email, web) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	result, err := r.db.ExecContext(ctx, query, assignmentUser.FirstName, assignmentUser.LastName, assignmentUser.CompanyName, assignmentUser.Address, assignmentUser.City, assignmentUser.County, assignmentUser.Postal, assignmentUser.Phone, assignmentUser.Email, assignmentUser.Web)
	if err != nil {
		log.Error(err.Error(), reqID)
		return 0, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	assignmentUserID, err := result.LastInsertId()
	if err != nil {
		log.Error(err.Error(), reqID)
		return 0, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	log.Info("core>repo>assignmentUser: CreateAssignmentUser completed & assignment user id is "+strconv.Itoa(int(assignmentUserID)), reqID)
	return int(assignmentUserID), nil
}

func (r *AssignmentUserRepoImpl) PartialUpdateAssignmentUser(ctx context.Context, assignmentUser entity.AssignmentUser) errors.Response {
	reqID, _ := falconContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>assignmentUser: PartialUpdateAssignmentUser started for assignment user id "+strconv.Itoa(assignmentUser.ID), reqID)

	columns := []string{}
	args := []interface{}{}

	if assignmentUser.Address != "" {
		columns = append(columns, "address=?")
		args = append(args, assignmentUser.Address)
	}
	if assignmentUser.City != "" {
		columns = append(columns, "city=?")
		args = append(args, assignmentUser.City)
	}
	if assignmentUser.CompanyName != "" {
		columns = append(columns, "company_name=?")
		args = append(args, assignmentUser.CompanyName)
	}
	if assignmentUser.County != "" {
		columns = append(columns, "county=?")
		args = append(args, assignmentUser.County)
	}
	if assignmentUser.Email != "" {
		columns = append(columns, "email=?")
		args = append(args, assignmentUser.Email)
	}
	if assignmentUser.FirstName != "" {
		columns = append(columns, "first_name=?")
		args = append(args, assignmentUser.FirstName)
	}
	if assignmentUser.LastName != "" {
		columns = append(columns, "last_name=?")
		args = append(args, assignmentUser.LastName)
	}
	if assignmentUser.Phone != "" {
		columns = append(columns, "phone=?")
		args = append(args, assignmentUser.Phone)
	}
	if assignmentUser.Postal != "" {
		columns = append(columns, "postal=?")
		args = append(args, assignmentUser.Postal)
	}
	if assignmentUser.Web != "" {
		columns = append(columns, "web=?")
		args = append(args, assignmentUser.Web)
	}

	if len(columns) == 0 {
		return nil
	}

	args = append(args, assignmentUser.ID)
	query := "UPDATE assignment_users SET " + strings.Join(columns, ", ") + " WHERE id=? AND assignment_users.deleted_at IS NULL"

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		log.Error(err.Error(), reqID)
		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	log.Info("core>repo>assignmentUser: PartialUpdateAssignmentUser completed for assignment user id "+strconv.Itoa(assignmentUser.ID), reqID)
	return nil
}

func (r *AssignmentUserRepoImpl) DeleteAssignmentUser(ctx context.Context, assignmentUserID int) errors.Response {
	reqID, _ := falconContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>assignmentUser: DeleteAssignmentUser started for assignment user id "+strconv.Itoa(assignmentUserID), reqID)

	query := "UPDATE assignment_users SET deleted_at = ? WHERE id = ?"

	_, err := r.db.ExecContext(ctx, query, time.Now(), assignmentUserID)
	if err != nil {
		log.Error(err.Error(), reqID)
		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	log.Info("core>repo>assignmentUser: DeleteAssignmentUser completed for assignment user id "+strconv.Itoa(assignmentUserID), reqID)
	return nil
}

func (r *AssignmentUserRepoImpl) GetAssignmentUserbyID(ctx context.Context, assignmentUserID int) (entity.AssignmentUser, errors.Response) {
	reqID, _ := falconContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>assignmentUser: GetAssignmentUserbyID started for assignment user id "+strconv.Itoa(assignmentUserID), reqID)

	query := "SELECT * FROM assignment_users WHERE assignment_users.id=? AND assignment_users.deleted_at IS NULL"

	assignmentUserModel := model.AssignmentUser{}
	assignmentUserEntity := entity.AssignmentUser{}

	err := r.db.Get(&assignmentUserModel, query, assignmentUserID)
	if err != nil {
		log.Error(err.Error(), reqID)
		return assignmentUserEntity, errors.ResponseNotFoundError(errors.NOT_FOUND_ERROR)
	}

	assignmentUserEntity = converter.AssignmentUserModelToAssignmentUserEntity(assignmentUserModel)

	log.Info("core>repo>assignmentUser: GetAssignmentUserbyID completed for assignment user id "+strconv.Itoa(assignmentUserID), reqID)
	return assignmentUserEntity, nil
}

func (r *AssignmentUserRepoImpl) ListAssignmentUser(ctx context.Context) (int, []entity.AssignmentUser, errors.Response) {
	reqID, _ := falconContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>assignmentUser: ListAssignmentUser started", reqID)

	queryStatement := "SELECT * FROM assignment_users WHERE assignment_users.deleted_at IS NULL ORDER BY id DESC"
	countQuery := "SELECT COUNT(id) as totalRecords FROM assignment_users WHERE assignment_users.deleted_at IS NULL"

	var count int
	err := r.db.Get(&count, countQuery)
	if err != nil {
		log.Error(err.Error(), reqID)
		return 0, nil, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	assignmentUsersModel := []model.AssignmentUser{}
	err = r.db.Select(&assignmentUsersModel, queryStatement)
	if err != nil {
		log.Error(err.Error(), reqID)
		return 0, nil, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	assignmentUserEntities := []entity.AssignmentUser{}
	for _, assignmentUserModel := range assignmentUsersModel {
		assignmentUserEntities = append(assignmentUserEntities, converter.AssignmentUserModelToAssignmentUserEntity(assignmentUserModel))
	}

	log.Info("core>repo>assignmentUser: ListAssignmentUser completed", reqID)
	return count, assignmentUserEntities, nil
}
