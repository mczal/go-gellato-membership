package controller

import (
	"go-gellato-membership/model"
	"go-gellato-membership/service"
	"strings"

	"github.com/go-siris/siris/context"
)

func PostUpdateUserPointAddition(ctx context.Context) {
	userID := ctx.Params().Get("id")
	stat, resp := service.UpdateUserPointAddition(userID)
	ctx.StatusCode(stat)
	ctx.JSON(resp)
}

func GetScanAllUser(ctx context.Context) {
	statScan, resScan := service.ScanUser()
	ctx.StatusCode(statScan)
	ctx.JSON(resScan)
}

func GetUserByID(ctx context.Context) {
	userID := ctx.Params().Get("id") // Or convert directly using: .Values().GetInt/GetInt64 etc...
	statQuery, resQuery := service.FindUserByID(userID)
	ctx.StatusCode(statQuery)
	ctx.JSON(resQuery)
}

func GetUserByEmail(ctx context.Context) {
	userEmail := ctx.Params().Get("email")
	if !strings.Contains(userEmail, "@") || !strings.Contains(userEmail, ".") {
		badReq := model.BaseResponse{
			Status:  400,
			Message: "Bad Request | Email validation format constraint",
		}
		ctx.StatusCode(400)
		ctx.JSON(badReq)
		return
	}
	statQuery, resQuery := service.FindUserByEmail(userEmail)
	ctx.StatusCode(statQuery)
	ctx.JSON(resQuery)
}
