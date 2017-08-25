package user

import (
	"regexp"
)

type RegisterForm struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	Email     string `json:"email"`
}

var emailRegexp = regexp.MustCompile(`^([a-zA-Z0-9][-_.a-zA-Z0-9]*)(@[-_.a-zA-Z0-9]+)?$`)

func (rf *RegisterForm) Validate() bool {
	if rf.FirstName == "" || rf.LastName == "" ||
		len(rf.Password) < 6 || !emailRegexp.MatchString(rf.Email) {
		return false
	}
	return true
}
