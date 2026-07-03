package converter

import (
	coreAPIModel "project/apimodel/core"
	"project/entity"
)

func AssignmentUserEntityToAssignmentUserAPIModelResponse(e entity.AssignmentUser) coreAPIModel.AssignmentUserResponse {
	list := coreAPIModel.AssignmentUserResponse{

		ID:          e.ID,
		FirstName:   e.FirstName,
		LastName:    e.LastName,
		CompanyName: e.CompanyName,
		Address:     e.Address,
		City:        e.City,
		County:      e.County,
		Postal:      e.Postal,
		Phone:       e.Phone,
		Email:       e.Email,
		Web:         e.Web,
	}
	return list
}

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
