package storage

import "github.com/dfryer1193/basic-web-authentication/models"

type UserStore interface {
	Set(username string, user models.User)
	Get(username string) (models.User, bool)
}
