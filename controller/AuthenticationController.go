package controller

import (
	"fmt"
	"go-gellato-membership/model"
	"go-gellato-membership/service"
	"strings"

	"go-gellato-membership/utility"

	"time"

	"go-gellato-membership/status"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-siris/siris/context"
)

func PostAuth(ctx context.Context) {
	var userInput model.User
	ctx.ReadJSON(&userInput)
	if len(userInput.Password) < 5 || !strings.Contains(userInput.Email, "@") || !strings.Contains(userInput.Email, ".") {
		errValidation := model.BaseResponse{
			Status:  400,
			Message: "Bad Request | Password must be at least 5 characters long",
		}
		ctx.StatusCode(400)
		ctx.JSON(errValidation)
		return
	}

	// Get Password From
	stat, resultUserEmail := service.FindUserByEmail(userInput.Email)
	if stat != 200 {
		ctx.StatusCode(stat)
		ctx.JSON(resultUserEmail)
		return
	}

	user := resultUserEmail.(model.BaseSingleResponse).Value.(model.UserDetailDtoResponse)
	// Do check password
	fmt.Printf("UserInput: %v | hashed: %v\n", userInput.Password, user.Password)
	if !utility.CheckPasswordHash(userInput.Password, user.Password) {
		errWrongPass := model.BaseResponse{
			Status:  401,
			Message: "Unauthorized | Incorrect password",
		}
		ctx.StatusCode(401)
		ctx.JSON(errWrongPass)
		return
	}

	if user.Status != status.HEALTHY {
		errNotConfirmed := model.BaseResponse{
			Status:  403,
			Message: "Forbidden | Account hasn't been confirmed yet",
		}
		ctx.StatusCode(403)
		ctx.JSON(errNotConfirmed)
		return
	}

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.

	// mczal: PleaseRead RFC 7519 for JWT Registered and Private Claim Name
	// link:  https://tools.ietf.org/html/rfc7519#section-4.1
	// yesterday := time.Now().AddDate(0, 0, -1).UnixNano() / 1000000
	nowAddWeek := time.Now().AddDate(0, 0, 7).UnixNano() / 1000000
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.UserId,
		"role": user.Role,
		// "nbf":  yesterday,
		"exp": nowAddWeek,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(utility.Configuration.Secret))
	if err != nil {
		errSignedString := model.BaseResponse{
			Status:  500,
			Message: "Internal Server Error | ErrTokenSignedString: " + err.Error(),
		}
		ctx.StatusCode(500)
		ctx.JSON(errSignedString)
		return
	}

	res := model.BaseSingleResponse{
		Status:  200,
		Message: "Success",
		Value:   tokenString,
	}
	ctx.StatusCode(200)
	ctx.JSON(res)
}
