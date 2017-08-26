package user

type UserInfo struct {
	FamilyName          *string `json:"family_name"`
	Name                *string `json:"name"`
	SecondName          *string `json:"second_name"`
	DateReceiving       *string `json:"date_receiving"`
	IssuedBy            *string `json:"issued_by"`
	DivisionNumber      *string `json:"division_number"`
	RegistrationAddress *string `json:"registration_address"`
	MailingAddress      *string `json:"mailing_address"`
	BirthDay            *string `json:"birth_day"`
	Sex                 *bool   `json:"sex"`
	HomePhone           *string `json:"home_phone"`
	MobilePhone         *string `json:"mobile_phone"`
	Citizenship         *string `json:"citizenship"`
	EMail               *string `json:"email"`
}
