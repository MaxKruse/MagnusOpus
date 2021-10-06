package structs

type RequestFilter struct {
	Limit  int `json:"Limit"`
	Offset int `json:"offset"`
}

type StaffPost struct {
	Role   string `json:"role"`
	UserId uint   `json:"user_id"`
}
