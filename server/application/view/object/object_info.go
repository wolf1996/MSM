package object

type ObjectInfo struct {
	Id               *int    `json:"id"`
	Name             *string `json:"name"`
	UserId           *int    `json:"user_id"`
	Adres            *string `json:"adres"`
}

type ObjectsInfo []ObjectInfo
