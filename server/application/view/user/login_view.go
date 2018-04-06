package user

import "regexp"

type LoginForm struct {
	EMail *string `json:"e_mail"`
	Pass  *string `json:"password"`
}

var mailAddressRE = regexp.MustCompile(`^([a-zA-Z0-9][-_.a-zA-Z0-9]*)(@[-_.a-zA-Z0-9]+)?$`)

func (a *LoginForm) Validate() bool {
	var valid bool
	if (a.EMail == nil )|| (a.Pass == nil){
		return false
	}
	valid = mailAddressRE.MatchString(*a.EMail)
	if valid {
		valid = len(*a.Pass) >= 6
	}
	return valid
}
