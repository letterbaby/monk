package common

const (
	FUNC_PING = 1
	FUNC_PONG = 2
)

type ServerInfo struct {
	Id  string `json:"id"` //ID
	IId int64
}
