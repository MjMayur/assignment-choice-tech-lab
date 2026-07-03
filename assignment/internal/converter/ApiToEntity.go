package converter

import (
	"project/entity"

	coreAPIModel "project/apimodel/core"
)

func AssignmentUserAPIToAssignmentUserEntity(request *coreAPIModel.AssignmentUser) entity.AssignmentUser {
	e := entity.AssignmentUser{

		FirstName:   request.FirstName,
		LastName:    request.LastName,
		CompanyName: request.CompanyName,
		Address:     request.Address,
		City:        request.City,
		County:      request.County,
		Postal:      request.Postal,
		Phone:       request.Phone,
		Email:       request.Email,
		Web:         request.Web,
	}

	return e
}
func UpdateAssignmentUserAPIRequestToAssignmentUserEntity(request *coreAPIModel.UpdateAssignmentUserRequest) entity.AssignmentUser {
	e := entity.AssignmentUser{

		FirstName:   request.FirstName,
		LastName:    request.LastName,
		CompanyName: request.CompanyName,
		Address:     request.Address,
		City:        request.City,
		County:      request.County,
		Postal:      request.Postal,
		Phone:       request.Phone,
		Email:       request.Email,
		Web:         request.Web,
	}

	return e
}

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
