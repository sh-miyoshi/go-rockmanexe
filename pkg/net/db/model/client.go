package model

type ClientInfo struct {
	ID  string
	Key string
}

type ClientHandler interface {
	Add(ent ClientInfo) error
	Delete(clientID string) error
	GetAll() ([]ClientInfo, error)
}
