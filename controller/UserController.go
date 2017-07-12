package controller

import (
	"fmt"
	"go-gellato-membership/model"
	"go-gellato-membership/service"
	"go-gellato-membership/utility"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-siris/siris/context"
)

func PostChangePassword(ctx context.Context) {
	var user model.User
	ctx.ReadJSON(&user)
	if len(user.ForgotPasswordToken) == 0 || user.ForgotPasswordToken == "null" || len(user.Password) < 5 {
		badReq := model.BaseResponse{
			Status:  400,
			Message: "Bad Request | ForgotPasswordToken didn't exist | Password didn't exist | Password must be at least 5 characters long",
		}
		ctx.StatusCode(400)
		ctx.JSON(badReq)
		return
	}
	stat, resp := service.FindUserByForgotPasswordToken(user.ForgotPasswordToken)
	if stat != 200 {
		ctx.StatusCode(stat)
		ctx.JSON(resp)
		return
	}
	userRes := resp.(model.BaseSingleResponse).Value.(model.UserDetailDtoResponse)
	statSuc, respSuc := service.UpdateUserPassword(userRes.UserId, user.Password)
	ctx.StatusCode(statSuc)
	ctx.JSON(respSuc)
}

func PostForgotPassword(ctx context.Context) {
	var user model.User
	ctx.ReadJSON(&user)
	if !strings.Contains(user.Email, "@") || !strings.Contains(user.Email, ".") {
		badReq := model.BaseResponse{
			Status:  400,
			Message: "Bad Request | Email or Fullname validation format constraint",
		}
		ctx.StatusCode(400)
		ctx.JSON(badReq)
		return
	}

	randToken, err := utility.GenerateRandomStringURLSafe(10)
	if err != nil {
		errResp := model.BaseResponse{
			Status:  500,
			Message: "Internal Server Error | GenerateRandomStringURLSafe",
		}
		ctx.StatusCode(500)
		ctx.JSON(errResp)
		return
	}

	stat, resp := service.UpdateUserForgotPasswordToken(user.Email, randToken)
	ctx.StatusCode(stat)
	ctx.JSON(resp)
}

func GetUserBydIDWithToken(ctx context.Context) {
	auth := ctx.GetHeader("Authorization") // Or convert directly using: .Values().GetInt/GetInt64 etc...

	token, err := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(utility.Configuration.Secret), nil
	})
	if err != nil {
		errParse := model.BaseResponse{
			Status:  401,
			Message: "Unauthorized | parse: " + err.Error(),
		}
		ctx.StatusCode(401)
		ctx.JSON(errParse)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["sub"].(string)
		statQuery, resQuery := service.FindUserByID(userID)
		ctx.StatusCode(statQuery)
		ctx.JSON(resQuery)
	} else {
		errElseClaim := model.BaseResponse{
			Status:  401,
			Message: "Unauthorized | errElseClaim",
		}
		ctx.StatusCode(401)
		ctx.JSON(errElseClaim)
	}
}
