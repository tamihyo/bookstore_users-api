package users

import (
	"strings"

	"github.com/tamihyo/bookstore_utils-go/rest_errors"
)

const (
	StatusActive = "active"
)

type User struct {
	Id          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
	Password    string `json:"password"`
	/* if we receive password in the field json do not take json
	address to complete password field
	do not place anyfield called password */
}

type Users []User

//function
// func Validate(user *User) rest_errors.RestErr {
// 	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
// 	if user.Email == "" {
// 		return rest_errors.NewBadRequestError("Invalid email address")
// 	}
// 	return nil
// }

//method
/*
assigning this method to the name of the method
the parameter and whateve will return on this method
so in this case, we're assigng validate method to the user struct
so this way, user knows how to validate itself and return whaterver
we have an error or not
*/
func (user *User) Validate() *rest_errors.RestErr {
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if user.Email == "" {
		return rest_errors.NewBadRequestError("Invalid email address")
	}
	user.Password = strings.TrimSpace(user.Password)
	if user.Password == "" {
		return rest_errors.NewBadRequestError("Invalid password")

	}
	return nil
}
