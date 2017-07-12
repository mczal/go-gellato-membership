package routes

import (
	"go-gellato-membership/controller"

	"github.com/go-siris/siris/core/router"
)

func UserParty(userRoutes router.Party) {
	userRoutes.Get("", controller.GetUserBydIDWithToken)
	// userRoutes.Post("/change-password",controller.)
}

func AdminParty(adminRoutes router.Party) {
	adminRoutes.Get("/users/{id:string}", controller.GetUserByID)
	adminRoutes.Get("/users/email/{email:string}", controller.GetUserByEmail)
	adminRoutes.Get("/users", controller.GetScanAllUser)
	adminRoutes.Post("/add-one-point/{id:string}", controller.PostUpdateUserPointAddition)
	// adminRoutes.Post("/users", controller.PostNewUser)
}
