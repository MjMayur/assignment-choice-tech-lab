
package entity

import "time"

type AssignmentUser struct {
	ID int 
	FirstName string 
	LastName string 
	CompanyName string 
	Address string 
	City string 
	County string 
	Postal string 
	Phone string 
	Email string 
	Web string 
	CreatedAt time.Time 
	UpdatedAt time.Time 
	DeletedAt time.Time 
}
