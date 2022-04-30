package util

type Payload struct {
	Id   string
	Role string
}

func CreateToken(payload interface{}) string {
	return "not implemented"
}
