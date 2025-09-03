package routes

import (
	"github.com/goravel/framework/facades"

	"goravel/app/http/controllers"
	"goravel/app/http/middleware"
	// "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"goravel/app/http/controllers/documents"
	"goravel/app/http/controllers/googleDocument"
	"goravel/app/http/controllers/userManagement"
	"goravel/app/http/controllers/role"


	
)

func Api() {
	// facades.Route().GlobalMiddleware(
	// 	middleware.RecoverNotify(), 
	// )


	// users
	userController := controllers.NewUserController()
	// Auth
	authController := controllers.NewAuthController()
	userManagementController := userManagement.NewUserManagementController()
	// roles
	roleController := role.NewRoleController()
	// permissions
	permissionController := controllers.NewPermissionController()

	
	facades.Route().Get("/users/{id}", userController.Show)
	facades.Route().Post("/login", authController.Login)
	facades.Route().Post("/register", authController.Register)

	

	facades.Route().Middleware(middleware.Auth()).Group(func(router route.Router) {
		router.Get("/profile", authController.Profile)
		router.Post("/logout", authController.Logout)
		router.Get("/users", userController.Index)
	    router.Post("/user/edit", userController.Edit)

		router.Get("/user-management/roles", userManagementController.ShowUserRole)
	    router.Post("/user-management/{id}/assign-role", userManagementController.AssignRole)
		router.Get("/roles/:id", roleController.Show)
	})


	facades.Route().Middleware(middleware.Auth(), middleware.CasbinMiddleware()).Group(func(router route.Router) {
		router.Get("/roles", roleController.Index)
		router.Post("/roles", roleController.Store)
		router.Post("/roles/:id", roleController.UpdatePermissions)
	})
	
	facades.Route().Get("/permissions", permissionController.Index)


	documentController := docuements.NewDocumentController()
	facades.Route().Get("/documents", documentController.Index)
	facades.Route().Post("/documents/upload", documentController.Store)
	facades.Route().Get("/documents/download/:id", documentController.Download)

	
	addGoogleDocumentController := googleDocument.NewAddGoogleDocumentController()
	facades.Route().Post("/add-google-documents", addGoogleDocumentController.AddGoogleDocument)

	editGoogleDocumentController := googleDocument.NewEditGoogleDocumentController()
	facades.Route().Get("/edit-google-documents/:id", editGoogleDocumentController.ShowGoogleDocument)
	facades.Route().Post("/edit-google-documents/:id", editGoogleDocumentController.EditGoogleDocument)
	facades.Route().Post("/remove-google-documents/:id", editGoogleDocumentController.RemoveGoogleDocument)

	googleDocumentController := googleDocument.NewGoogleDocumentController()
	facades.Route().Get("/google-documents", googleDocumentController.ListGoogleDocuments)
	facades.Route().Get("/google-documents/:id", googleDocumentController.ShowGoogleDocument)


	






	





	

}
