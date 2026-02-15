
package manualwire

import (
	"github.com/JerryJeager/r3sonance-backend/config"
	"github.com/JerryJeager/r3sonance-backend/internal/http"
	"github.com/JerryJeager/r3sonance-backend/internal/service/users"
)

func GetUserRepository() *users.UserRepo {
	repo := config.GetSession()
	return users.NewUserRepo(repo)
}

func GetUserService(repo users.UserStore) *users.UserServ {
	return users.NewUserService(repo)
}

func GetUserController() *http.UserController {
	repo := GetUserRepository()
	service := GetUserService(repo)
	return http.NewUserController(service)
}

	