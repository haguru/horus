package interfaces

import "context"

type DbClient interface {
	CreateUser(username string, email string, password string) error
	GetUser(email string) error
	UpdatePassword(email string) error
	Disconnect(context.Context) error
	DeleteUser(email string) error
}
