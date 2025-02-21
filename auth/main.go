package main

import (
	"github.com/lpsaldana/go-appointment-booking-microservices/auth/internal/repositories"
	"github.com/lpsaldana/go-appointment-booking-microservices/auth/internal/services"
)

func main() {
	userRepository := repositories.NewUserRepository()
	userService := services.NewUserService(userRepository)

	userService.ValidateCredentials("test", "1234test")
}
