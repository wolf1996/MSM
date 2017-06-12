package user

import "regexp"

type LoginForm struct {
	EMail string `json:"eMail"`
	Pass string	 `json:"passWord"`
}

var mailAddressRE = regexp.MustCompile(`^([a-zA-Z0-9][-_.a-zA-Z0-9]*)(@[-_.a-zA-Z0-9]+)?$`)

func (a *LoginForm)Validate() bool {
	var valid bool
	valid = mailAddressRE.MatchString(a.EMail)
	if valid{
		valid = len(a.Pass)>=6
	}
	return valid
}