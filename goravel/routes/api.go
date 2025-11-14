package routes

import (
	"goravel/app/http/controllers"
	"goravel/app/http/controllers/adsCampaign" 
	"goravel/app/http/controllers/adsTracking"
	"goravel/app/http/controllers/documents"
	"goravel/app/http/controllers/gmail"
	"goravel/app/http/controllers/googleAuthenticator"
	"goravel/app/http/controllers/googleCalendar"
	"goravel/app/http/controllers/googleDocument"
	"goravel/app/http/controllers/role"
	"goravel/app/http/controllers/schedules"
	"goravel/app/http/controllers/userManagement"
	"goravel/app/http/middleware"

	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"
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
	//document
	documentController := docuements.NewDocumentController()
	//google document
	googleDocumentController := googleDocument.NewGoogleDocumentController()
	addGoogleDocumentController := googleDocument.NewAddGoogleDocumentController()
	editGoogleDocumentController := googleDocument.NewEditGoogleDocumentController()
	//schedule
	scheduleController := schedules.NewScheduleController()
	addScheduleController := schedules.NewAddScheduleController()
	editScheduleController := schedules.NewEditScheduleController()
	googleCalendarAccountController := googleCalendar.NewGoogleCalendarAccountController()
	// gmail
	gmailController := gmail.NewGmailController()
	replyGmailController := gmail.NewReplyGmailController()
	forwardGmailController := gmail.NewForwardGmailController()
	templateController := gmail.NewTemplateController()
	// ads tracking
	addAdsCampaignController := adsCampaign.NewAddAdsCampaignController()
	editAdsCampaignController := adsCampaign.NewEditAdsCampaignController()
	adsCampaignController := adsCampaign.NewAdsCampaignController()
	reportAdsCampaignController := adsCampaign.NewReportAdsCampaignController()
	adsTrackingController := adsTracking.NewAdsTrackingController()
	//setting google authenticator
	twofaController := googleAuthenticator.NewTwoFAController()

	facades.Route().Post("/login", authController.Login)
	facades.Route().Post("/login/twofa",authController.VerifyTwoFA)
	facades.Route().Post("/register", authController.Register)

	facades.Route().Middleware(middleware.ActivityLogger(),middleware.Auth()).Group(func(router route.Router) {
		//profile
		router.Get("/profile", authController.Profile)
		router.Post("/logout", authController.Logout)
		router.Get("/users", userController.Index)
	    router.Post("/user/edit", userController.Edit)
		//user management
		router.Get("/user-management/roles", userManagementController.ShowUserRole)
	    router.Post("/user-management/{id}/assign-role", userManagementController.AssignRole)
		router.Get("/roles/:id", roleController.Show)
		
		// 2 factor authentication
		router.Get("/twofa/qrcode",twofaController.GenerateQRCode)
		router.Post("/twofa/enable",twofaController.ConfirmEnable)
		router.Post("/twofa/disable",twofaController.ConfirmDisable)


		// permission
		router.Get("/permissions", permissionController.Index)

		//documents
		router.Get("/documents", documentController.Index)
		router.Post("/documents/upload", documentController.Store)
		router.Get("/documents/download/:id", documentController.Download)

		//google documents
		router.Post("/add-google-documents", addGoogleDocumentController.AddGoogleDocument)
		router.Get("/edit-google-documents/:id", editGoogleDocumentController.ShowGoogleDocument)
		router.Post("/edit-google-documents/:id", editGoogleDocumentController.EditGoogleDocument)
		router.Post("/remove-google-documents/:id", editGoogleDocumentController.RemoveGoogleDocument)
		router.Get("/google-documents", googleDocumentController.ListGoogleDocuments)
		router.Get("/google-documents/:id", googleDocumentController.ShowGoogleDocument)


		//schedules
		router.Get("/schedules",scheduleController.ShowSchedule)
		router.Post("/add-schedules",addScheduleController.AddSchedule)
		router.Get("/edit-schedules/:id",editScheduleController.ShowSchedule)
		router.Post("/edit-schedules/:id",editScheduleController.EditSchedule)

		//email
		router.Get("/gmail/auth", gmailController.RedirectToGoogle)
		router.Get("/gmail/accounts", gmailController.ListAccounts)
		router.Post("/gmail/remove-accounts/:id", gmailController.DeleteAccount)
		router.Get("/gmail/accounts/teams", gmailController.GetGmailAccountTeams)
		router.Post("/gmail/technical/messages/:id/star", gmailController.ToggleStar)
		router.Post("/gmail/support/messages/:id/star", gmailController.ToggleStar)
		router.Get("/gmail/templates", templateController.ShowTemplates)
		router.Post("/gmail/templates", templateController.AddTemplate)
		router.Post("/gmail/templates/edit/:id", templateController.EditTemplate)
		router.Post("/gmail/templates/remove/:id", templateController.RemoveTemplate)

		//Ads Tracking
		// Add ads campaign
		router.Post("/add-ads-campaign",addAdsCampaignController.AddAdsCampaign)
		//List ads campaign
		router.Get("/ads-campaign",adsCampaignController.ListAdsCampaigns)
		// Show edit ads campaign
		router.Get("/edit-ads-campaign/:id",editAdsCampaignController.ShowAdsCampaign)
		router.Post("/edit-ads-campaign/:id",editAdsCampaignController.EditAdsCampaign)
		//report ads campaign
		router.Get("/ads-campaign/report/:campaign_id",reportAdsCampaignController.ShowReportAdsCampaign)
		router.Get("/ads-campaign/report/:campaign_id/filter",reportAdsCampaignController.ShowReportAdsFilterCampaign)
		// show support parameter
		router.Get("/add-ads-campaign/support-parameters",addAdsCampaignController.ShowSupportParameter)





		
	})


	//rbac casbin middleware.CasbinMiddleware()
	facades.Route().Middleware(middleware.Auth(), middleware.CasbinMiddleware()).Group(func(router route.Router) {
		router.Get("/roles", roleController.Index)
		router.Post("/roles", roleController.Store)
		router.Post("/roles/:id", roleController.UpdatePermissions)

		//Gmail routes
		router.Get("/gmail/technical/messages", gmailController.ListMessages)
		router.Get("/gmail/technical/messages/:id", gmailController.ReadMessage)
		router.Post("/gmail/technical/messages/:id/reply", replyGmailController.ReplyMessage)
		router.Post("/gmail/technical/messages/forward", forwardGmailController.ForwardMessage)

		router.Get("/gmail/support/messages", gmailController.ListMessages)
		router.Get("/gmail/support/messages/:id", gmailController.ReadMessage)
		router.Post("/gmail/support/messages/:id/reply", replyGmailController.ReplyMessage)
		router.Post("/gmail/support/messages/forward", forwardGmailController.ForwardMessage)

	})
	
	//gmail login callback
	facades.Route().Get("/oauth/google/callback", gmailController.HandleCallback)

	facades.Route().Get("/:code",adsTrackingController.Track)
	facades.Route().Post("/postback",adsTrackingController.PostBackAdsTracking)



	
	facades.Route().Get("/ads-campaign/report/ads-log-details/:campaign_id",reportAdsCampaignController.ShowReportAdsLogDetailsCampaign)
	
	facades.Route().Get("/google/schedule/account", googleCalendarAccountController.ShowGoogleAccount)
	facades.Route().Get("/google/auth/schedule/url", googleCalendarAccountController.AuthURL)
	facades.Route().Get("/oauth/google/schedule/callback", googleCalendarAccountController.Callback)



}
