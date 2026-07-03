package apimodel

type AssignmentUserResponse struct {
	ID int `json:"ID"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	CompanyName string `json:"companyName"`
	Address string `json:"address"`
	City string `json:"city"`
	County string `json:"county"`
	Postal string `json:"postal"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	Web string `json:"web"`
}

type AssignmentUserListResponse struct {
	TotalRecords      	 int        	            `json:"totalRecords"`
	NoOfRecordsPerPage 	 int               	     `json:"noOfRecordsPerPage"`
	AssignmentUser 	[]AssignmentUserResponse
}

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
