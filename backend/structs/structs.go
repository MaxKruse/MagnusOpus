package structs

type Config struct {
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_DB       string
	POSTGRES_URL      string
}

type RippleSelf struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	UserId   int    `json:"id"`
	Username string `json:"username"`
}

type BanchoMe struct {
	Id       int    `json:"id ,omitempty"`
	Username string `json:"username ,omitempty"`
}
