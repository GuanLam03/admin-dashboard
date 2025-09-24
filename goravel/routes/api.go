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
	"goravel/app/http/controllers/schedules"
	"goravel/app/http/controllers/googleAuthenticator"
	"goravel/app/http/controllers/gmail"






	
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

	//setting google authenticator
	twofaController := googleAuthenticator.NewTwoFAController()

	//Gmail
	gmailController := gmail.NewGmailController()

	
	facades.Route().Get("/users/{id}", userController.Show)
	facades.Route().Post("/login", authController.Login)
	facades.Route().Post("/login/twofa",authController.VerifyTwoFA)
	facades.Route().Post("/register", authController.Register)


	facades.Route().Middleware(middleware.Auth()).Group(func(router route.Router) {
		router.Get("/profile", authController.Profile)
		router.Post("/logout", authController.Logout)
		router.Get("/users", userController.Index)
	    router.Post("/user/edit", userController.Edit)

		router.Get("/user-management/roles", userManagementController.ShowUserRole)
	    router.Post("/user-management/{id}/assign-role", userManagementController.AssignRole)
		router.Get("/roles/:id", roleController.Show)
		
		// 2 factor authentication
		router.Get("/twofa/qrcode",twofaController.GenerateQRCode)
		router.Post("/twofa/enable",twofaController.ConfirmEnable)
		router.Post("/twofa/disable",twofaController.ConfirmDisable)

		
	})


	//, middleware.CasbinMiddleware()
	facades.Route().Middleware(middleware.Auth(), middleware.CasbinMiddleware()).Group(func(router route.Router) {
		router.Get("/roles", roleController.Index)
		router.Post("/roles", roleController.Store)
		router.Post("/roles/:id", roleController.UpdatePermissions)


		//Gmail routes
		router.Get("/gmail/technical/messages", gmailController.ListMessages)
		router.Get("/gmail/technical/messages/:id", gmailController.ReadMessage)
		router.Post("/gmail/technical/messages/:id/reply", gmailController.ReplyMessage)



		router.Get("/gmail/support/messages", gmailController.ListMessages)
		router.Get("/gmail/support/messages/:id", gmailController.ReadMessage)
		router.Post("/gmail/support/messages/:id/reply", gmailController.ReplyMessage)
		
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




	scheduleController := schedules.NewScheduleController()
	facades.Route().Get("/schedules",scheduleController.ShowSchedule)

	addScheduleController := schedules.NewAddScheduleController()
	facades.Route().Post("/add-schedules",addScheduleController.AddSchedule)

	editScheduleController := schedules.NewEditScheduleController()
	facades.Route().Get("/edit-schedules/:id",editScheduleController.ShowSchedule)
	facades.Route().Post("/edit-schedules/:id",editScheduleController.EditSchedule)



	facades.Route().Get("/gmail/auth", gmailController.RedirectToGoogle)
	facades.Route().Get("/oauth/google/callback", gmailController.HandleCallback)







	





	





	

}
