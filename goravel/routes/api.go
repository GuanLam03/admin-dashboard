package routes

import (
	"github.com/goravel/framework/facades"

	"goravel/app/http/controllers"
	"goravel/app/http/middleware"

	
)

func Api() {
	userController := controllers.NewUserController()
	facades.Route().Get("/users/{id}", userController.Show)

	facades.Route().Middleware(middleware.Auth()).Post("/user/edit", userController.Edit)

	// Auth example
	authController := controllers.NewAuthController()

	facades.Route().Post("/login", authController.Login)
	facades.Route().Post("/register", authController.Register)
	// facades.Route().Get("/profile", authController.Profile)


	
	facades.Route().Middleware(middleware.Auth()).Get("/profile", authController.Profile)
	facades.Route().Middleware(middleware.Auth()).Post("/logout", authController.Logout)




	

}
