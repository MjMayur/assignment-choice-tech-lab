package apimodel

type AssignmentUser struct {
	FirstName   string `json:"firstName" binding:"required"`
	LastName    string `json:"lastName" binding:"required"`
	CompanyName string `json:"companyName" binding:"omitempty"`
	Address     string `json:"address" binding:"omitempty"`
	City        string `json:"city" binding:"omitempty"`
	County      string `json:"county" binding:"omitempty"`
	Postal      string `json:"postal" binding:"omitempty"`
	Phone       string `json:"phone" binding:"omitempty"`
	Email       string `json:"email" binding:"omitempty"`
	Web         string `json:"web" binding:"omitempty"`
}

type UpdateAssignmentUserRequest struct {
	FirstName   string `json:"firstName" binding:"omitempty"`
	LastName    string `json:"lastName" binding:"omitempty"`
	CompanyName string `json:"companyName" binding:"omitempty"`
	Address     string `json:"address" binding:"omitempty"`
	City        string `json:"city" binding:"omitempty"`
	County      string `json:"county" binding:"omitempty"`
	Postal      string `json:"postal" binding:"omitempty"`
	Phone       string `json:"phone" binding:"omitempty"`
	Email       string `json:"email" binding:"omitempty"`
	Web         string `json:"web" binding:"omitempty"`
}

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
