package controller

import (
	"go-gellato-membership/model"
	"go-gellato-membership/service"
	"go-gellato-membership/utility"
	"strconv"
	"strings"

	"github.com/go-siris/siris/context"
	qrcode "github.com/skip2/go-qrcode"
)

func PostConfirmAccount(ctx context.Context) {
	var user model.User
	ctx.ReadJSON(&user)
	if len(user.ConfirmationToken) == 0 || user.ConfirmationToken == "null" {
		errValidation := model.BaseResponse{
			Status:  400,
			Message: "Bad Request | ConfirmationToken attribute didn't exist",
		}
		ctx.StatusCode(400)
		ctx.JSON(errValidation)
		return
	}

	stat, resp := service.FindUserByConfirmationToken(user.ConfirmationToken)
	if stat != 200 {
		ctx.StatusCode(stat)
		ctx.JSON(resp)
		return
	}

	userResp := resp.(model.BaseSingleResponse).Value.(model.UserDetailDtoResponse)

	err := qrcode.WriteFile(userResp.UserId, qrcode.Medium, 256, "./public/"+userResp.UserId+".png")
	if err != nil {
		errQrCode := model.BaseResponse{
			Status:  500,
			Message: "Internal Server Error | QrcodeWriteFile: " + err.Error(),
		}
		ctx.StatusCode(500)
		ctx.JSON(errQrCode)
		return
	}
	stat, resp = service.UpdateUserStatusHealthy(userResp.UserId)
	ctx.StatusCode(stat)
	ctx.JSON(resp)
}

func PostResendConfirmationEmail(ctx context.Context) {
	var user model.User
	ctx.ReadJSON(&user)
	if !strings.Contains(user.Email, "@") || !strings.Contains(user.Email, ".") {
		ctx.StatusCode(400)
		errValidation := model.BaseResponse{
			Status:  400,
			Message: "Bad Request | EmailValidationConstraint",
		}
		ctx.JSON(errValidation)
		return
	}
	stat, resultUser := service.FindUserByEmail(user.Email)
	if stat != 200 {
		errValidation := model.BaseResponse{
			Status:  404,
			Message: "Not Found | User with email: " + user.Email + " didn't exist.",
		}
		ctx.StatusCode(404)
		ctx.JSON(errValidation)
		return
	}

	userEnt := resultUser.(model.BaseSingleResponse).Value.(model.UserDetailDtoResponse)
	service.SendConfirmationMail(userEnt.FullName, userEnt.Email, userEnt.ConfirmationToken)
	successResp := model.BaseResponse{
		Status:  200,
		Message: "Success",
	}
	ctx.StatusCode(200)
	ctx.JSON(successResp)
}

func validateNewUserRequest(user *model.User) (bool, string) {
	// "email": "mczal@nawar.in",
	// "fullName": "Fahrizal Septrianto",
	// "password": "123",
	// "address": "Jl Jawa 66",
	// "birthDate": 812246400000, 12 length
	// "phoneNumber": "0856272"
	if !strings.Contains(user.Email, "@") || !strings.Contains(user.Email, ".") {
		return false, "EmailValidationConstraint: "
	}
	if len(user.FullName) <= 5 {
		return false, "FullnameValidationConstraint: Fullname length must be more than 5 characters long"
	}
	if len(user.Password) < 5 {
		return false, "PasswordValidationConstraint: Password length must be at least 5 characters long"
	}
	if len(user.Address) < 5 {
		return false, "AddressValidationConstraint: Address length must be at least 5 characters long"
	}
	if len(user.PhoneNumber) < 5 {
		return false, "PhoneNumberValidationConstraint: PhoneNumber length must be at least 10 characters long"
	}
	if len(strconv.Itoa(int(user.BirthDate))) != 12 {
		return false, "BirthDateValidationConstraint: BirthDate digit length must be 12 characters long and form time in a millisecond epoch format"
	}

	return true, ""
}

func PostNewUser(ctx context.Context) {
	var user model.User
	ctx.ReadJSON(&user)
	ok, strerr := validateNewUserRequest(&user)
	if !ok {
		errValid := model.BaseResponse{
			Status:  400,
			Message: "BadRequest | " + strerr,
		}
		ctx.StatusCode(400)
		ctx.JSON(errValid)
		return
	}
	confirmationToken, err := utility.GenerateRandomStringURLSafe(10)
	if err != nil {
		ctx.StatusCode(500)
		errGenRandom := model.BaseResponse{
			Status:  500,
			Message: "Internal Server Error | " + err.Error(),
		}
		ctx.JSON(errGenRandom)
		return
	}
	user.ConfirmationToken = confirmationToken
	stat, result := service.NewUser(&user)
	if stat == 200 {
		service.SendConfirmationMail(user.FullName, user.Email, confirmationToken)
	}
	ctx.StatusCode(stat)
	ctx.JSON(result)
}
