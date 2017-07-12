package main

import (
	"go-gellato-membership/controller"
	"go-gellato-membership/middleware"
	"go-gellato-membership/routes"
	"go-gellato-membership/service"

	"go-gellato-membership/utility"

	"github.com/go-siris/siris"
	"github.com/go-siris/siris/context"
)

func main() {
	utility.InitializeConfig()
	// fmt.Printf("Config: %v", utility.Configuration)
	service.GenerateDynamoDBSvc()

	app := siris.New()

	// Regster custom handler for specific http errors.
	app.OnErrorCode(siris.StatusInternalServerError, func(ctx context.Context) {
		// .Values are used to communicate between handlers, middleware.
		errMessage := ctx.Values().GetString("error")
		if errMessage != "" {
			ctx.Writef("Internal server error: %s", errMessage)
			return
		}

		ctx.Writef("(Unexpected) internal server error")
	})

	// mczal: Applying Log Middleware To All Request
	app.Use(middleware.LogMiddleware)

	// app.Done(func(ctx context.Context) {})

	app.Post("/register", controller.PostNewUser)
	app.Post("/confirm", controller.PostConfirmAccount)
	app.Post("/resend-confirmation-email", controller.PostResendConfirmationEmail)

	app.Post("/forgot-password", controller.PostForgotPassword)
	app.Post("/change-password", controller.PostChangePassword)

	app.Post("/auth", controller.PostAuth)

	routes.UserParty(app.Party("/user", middleware.ValidateTokenUser))
	routes.AdminParty(app.Party("/admin", middleware.ValidateTokenAdmin))

	app.StaticWeb("/static", "./public")

	// Listen for incoming HTTP/1.x & HTTP/2 clients on localhost port 8080.
	app.Run(siris.Addr(":"+utility.Configuration.Port), siris.WithCharset("UTF-8"))
}
