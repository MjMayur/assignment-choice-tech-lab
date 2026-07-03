package converter

import (
	"project/entity"
	"project/model"
)

func AssignmentUserModelToAssignmentUserEntity(m model.AssignmentUser) entity.AssignmentUser {
	e := entity.AssignmentUser{

		ID:          m.ID,
		FirstName:   m.FirstName,
		LastName:    m.LastName,
		CompanyName: m.CompanyName.String,
		Address:     m.Address.String,
		City:        m.City.String,
		County:      m.County.String,
		Postal:      m.Postal.String,
		Phone:       m.Phone.String,
		Email:       m.Email.String,
		Web:         m.Web.String,
	}
	return e
}

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
