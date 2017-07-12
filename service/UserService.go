package service

import (
	"fmt"
	"go-gellato-membership/model"
	"go-gellato-membership/status"
	"go-gellato-membership/utility"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	uuid "github.com/satori/go.uuid"
)

// func UpdateUserPassword(newPassword, token string) (int, interface{}) {
// ForgotPasswordToken-index
// stat, resp :=
// }

func UpdateUserPassword(userID, newPass string) (int, interface{}) {
	hashed, errH := utility.HashPassword(newPass)
	if errH != nil {
		errRes := model.BaseResponse{
			Status:  500,
			Message: "Internal Server Errror | Hashing",
		}
		return 500, errRes
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#P": aws.String("Password"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":p": {
				S: aws.String(hashed),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"UserId": {
				S: aws.String(userID),
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		TableName:        aws.String(utility.Configuration.Dynamodb_dbname),
		UpdateExpression: aws.String("SET #P = :p"),
	}

	_, err := svc.UpdateItem(input)
	if err != nil {
		errUpdateItem := model.BaseResponse{
			Status:  500,
			Message: "Internal Server Error: " + err.Error() + " | ",
		}
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				fmt.Println(dynamodb.ErrCodeConditionalCheckFailedException, aerr.Error())
				errUpdateItem.Message += "ErrCodeConditionalCheckFailedException: " + aerr.Error()
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
				errUpdateItem.Message += "ErrCodeProvisionedThroughputExceededException: " + aerr.Error()
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
				errUpdateItem.Message += "ErrCodeResourceNotFoundException: " + aerr.Error()
			case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
				fmt.Println(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
				errUpdateItem.Message += "ErrCodeItemCollectionSizeLimitExceededException: " + aerr.Error()
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
				errUpdateItem.Message += "ErrCodeInternalServerError: " + aerr.Error()
			default:
				errUpdateItem.Message += aerr.Error()
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return 500, errUpdateItem
	}

	succResp := model.BaseResponse{
		Status:  200,
		Message: "Success change password",
	}
	return 200, succResp
}

func FindUserByForgotPasswordToken(token string) (int, interface{}) {
	inputCheck := &dynamodb.QueryInput{
		TableName: aws.String(utility.Configuration.Dynamodb_dbname),
		IndexName: aws.String("ForgotPasswordToken-index"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":t": {
				S: aws.String(token),
			},
		},
		KeyConditionExpression: aws.String("ForgotPasswordToken = :t"),
	}

	result, err := svc.Query(inputCheck)
	if err != nil {
		errQuerying := model.BaseResponse{
			Status:  500,
			Message: "Internal Server Error: Err Querying | ",
		}
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
				errQuerying.Message += "ErrCodeProvisionedThroughputExceededException: " + aerr.Error()
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
				errQuerying.Message += "ErrCodeResourceNotFoundException: " + aerr.Error()
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
				errQuerying.Message += "ErrCodeInternalServerError: " + aerr.Error()
			default:
				fmt.Println(aerr.Error())
				errQuerying.Message += aerr.Error()
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return 500, errQuerying
	}

	users := []model.UserDetailDtoResponse{}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		errUnmarshal := model.BaseResponse{
			Status:  500,
			Message: "Internal Server Error: Err UnmarshalListOfMaps | failed to unmarshal Dynamodb Query Items: " + err.Error(),
		}
		return 500, errUnmarshal
	}

	if len(users) == 0 {
		errNotFound := model.BaseResponse{
			Message: "Error ForgotPasswordToken: " + token + " didn't exist!",
			Status:  404,
		}
		return 404, errNotFound
	}
	succResp := model.BaseSingleResponse{
		Message: "Success: " + token,
		Status:  200,
		Value:   users[0],
	}
	return 200, succResp

}

func UpdateUserForgotPasswordToken(email, token string) (int, interface{}) {
	stat, resp := FindUserByEmail(email)
	if stat != 200 {
		return stat, resp
	}
	user := resp.(model.BaseSingleResponse).Value.(model.UserDetailDtoResponse)

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#FPT": aws.String("ForgotPasswordToken"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":fpt": {
				S: aws.String(token),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"UserId": {
				S: aws.String(user.UserId),
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		TableName:        aws.String(utility.Configuration.Dynamodb_dbname),
		UpdateExpression: aws.String("SET #FPT = :fpt"),
	}
	_, err := svc.UpdateItem(input)
	if err != nil {
		errUpdateItem := model.BaseResponse{
			Status:  500,
			Message: "Internal Server Error: " + err.Error() + " | ",
		}
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				fmt.Println(dynamodb.ErrCodeConditionalCheckFailedException, aerr.Error())
				errUpdateItem.Message += "ErrCodeConditionalCheckFailedException: " + aerr.Error()
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
				errUpdateItem.Message += "ErrCodeProvisionedThroughputExceededException: " + aerr.Error()
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
				errUpdateItem.Message += "ErrCodeResourceNotFoundException: " + aerr.Error()
			case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
				fmt.Println(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
				errUpdateItem.Message += "ErrCodeItemCollectionSizeLimitExceededException: " + aerr.Error()
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
				errUpdateItem.Message += "ErrCodeInternalServerError: " + aerr.Error()
			default:
				errUpdateItem.Message += aerr.Error()
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return 500, errUpdateItem
	}

	SendPasswordConfirmation(user.FullName, email, token)

	succResp := model.BaseResponse{
		Status:  200,
		Message: "Success send forgot password token to " + email,
	}
	return 200, succResp
}

func UpdateUserPointAddition(userID string) (int, interface{}) {
	nowString := strconv.Itoa(int(time.Now().UnixNano() / 1000000))

	stat, resp := FindUserByID(userID)
	if stat != 200 {
		return stat, resp
	}
	user := resp.(model.BaseSingleResponse).Value.(model.UserDetailDtoResponse)

	var (
		input        *dynamodb.UpdateItemInput
		currentPoint int
	)

	if user.Point+1 == utility.Configuration.BonusPoint {
		currentPoint = 0
		input = &dynamodb.UpdateItemInput{
			ExpressionAttributeNames: map[string]*string{
				"#PP": aws.String("PointStatusLogs"),
				"#P":  aws.String("Point"),
				"#PR": aws.String("PrestigeRank"),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":pp": {
					L: []*dynamodb.AttributeValue{
						{
							M: map[string]*dynamodb.AttributeValue{
								"CreatedDate": {
									N: aws.String(nowString),
								},
								"Point": {
									S: aws.String("-" + strconv.Itoa(utility.Configuration.BonusPoint)),
								},
							},
						},
						{
							M: map[string]*dynamodb.AttributeValue{
								"CreatedDate": {
									N: aws.String(nowString),
								},
								"Point": {
									S: aws.String("+1"),
								},
							},
						},
					},
				},
				":p": {
					N: aws.String("0"),
				},
				":pr": {
					N: aws.String("1"),
				},
			},
			Key: map[string]*dynamodb.AttributeValue{
				"UserId": {
					S: aws.String(userID),
				},
			},
			ReturnValues:     aws.String("ALL_NEW"),
			TableName:        aws.String(utility.Configuration.Dynamodb_dbname),
			UpdateExpression: aws.String("SET #PP = list_append(:pp,#PP), #P = :p ADD #PR :pr"),
		}
	} else {
		currentPoint = user.Point + 1
		input = &dynamodb.UpdateItemInput{
			ExpressionAttributeNames: map[string]*string{
				"#PP": aws.String("PointStatusLogs"),
				"#P":  aws.String("Point"),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":pp": {
					L: []*dynamodb.AttributeValue{
						{
							M: map[string]*dynamodb.AttributeValue{
								"CreatedDate": {
									N: aws.String(nowString),
								},
								"PointAddition": {
									S: aws.String("+1"),
								},
								"PointSubtraction": {
									S: aws.String("0"),
								},
							},
						},
					},
				},
				":p": {
					N: aws.String("1"),
				},
			},
			Key: map[string]*dynamodb.AttributeValue{
				"UserId": {
					S: aws.String(userID),
				},
			},
			ReturnValues:     aws.String("ALL_NEW"),
			TableName:        aws.String(utility.Configuration.Dynamodb_dbname),
			UpdateExpression: aws.String("SET #PP = list_append(:pp,#PP) ADD #P :p"),
		}
	}

	_, err := svc.UpdateItem(input)
	if err != nil {
		errUpdateItem := model.BaseResponse{
			Status:  500,
			Message: "Internal Server Error: " + err.Error() + " | ",
		}
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				fmt.Println(dynamodb.ErrCodeConditionalCheckFailedException, aerr.Error())
				errUpdateItem.Message += "ErrCodeConditionalCheckFailedException: " + aerr.Error()
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
				errUpdateItem.Message += "ErrCodeProvisionedThroughputExceededException: " + aerr.Error()
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
				errUpdateItem.Message += "ErrCodeResourceNotFoundException: " + aerr.Error()
			case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
				fmt.Println(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
				errUpdateItem.Message += "ErrCodeItemCollectionSizeLimitExceededException: " + aerr.Error()
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
				errUpdateItem.Message += "ErrCodeInternalServerError: " + aerr.Error()
			default:
				errUpdateItem.Message += aerr.Error()
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return 500, errUpdateItem
	}

	successRess := model.BaseResponse{
		Status:  200,
		Message: "Success update user point by +1 | current point: " + strconv.Itoa(currentPoint),
	}
	if utility.ContainsInt(utility.Configuration.NotifWhenPointAchieve, user.Point+1) {
		SendNotifPoint(user.FullName, user.Email, user.Point+1)
		successRess.Message += " | EmailNotifSent: for point = " + strconv.Itoa(user.Point+1)
	}
	if currentPoint == 0 {
		successRess.Status = 202
		successRess.Message += " | Congratulations you've used your bonus | PrestigeRank: " + strconv.Itoa(user.PrestigeRank+1)
	}
	return successRess.Status, successRess
}

func UpdateUserStatusHealthy(userID string) (int, interface{}) {
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#S": aws.String("Status"),
			"#T": aws.String("ConfirmationToken"),
			"#Q": aws.String("Qrcode"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":s": {
				S: aws.String(status.HEALTHY),
			},
			":t": {
				S: aws.String("null"),
			},
			":q": {
				S: aws.String(utility.Configuration.StaticBasePathLocation + "/static/" + userID + ".png"),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"UserId": {
				S: aws.String(userID),
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		TableName:        aws.String(utility.Configuration.Dynamodb_dbname),
		UpdateExpression: aws.String("SET #S = :s, #T = :t, #Q = :q"),
	}
	_, err := svc.UpdateItem(input)
	if err != nil {
		errUpdateItem := model.BaseResponse{
			Status:  500,
			Message: "Internal Server Error: " + err.Error() + " | ",
		}
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				fmt.Println(dynamodb.ErrCodeConditionalCheckFailedException, aerr.Error())
				errUpdateItem.Message += "ErrCodeConditionalCheckFailedException: " + aerr.Error()
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
				errUpdateItem.Message += "ErrCodeProvisionedThroughputExceededException: " + aerr.Error()
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
				errUpdateItem.Message += "ErrCodeResourceNotFoundException: " + aerr.Error()
			case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
				fmt.Println(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
				errUpdateItem.Message += "ErrCodeItemCollectionSizeLimitExceededException: " + aerr.Error()
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
				errUpdateItem.Message += "ErrCodeInternalServerError: " + aerr.Error()
			default:
				errUpdateItem.Message += aerr.Error()
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return 500, errUpdateItem
	}
	successRess := model.BaseResponse{
		Status:  200,
		Message: "Success update status user",
	}
	return 200, successRess
}

func FindUserByConfirmationToken(confirmationToken string) (int, interface{}) {
	input := &dynamodb.QueryInput{
		TableName: aws.String(utility.Configuration.Dynamodb_dbname),
		IndexName: aws.String("ConfirmationToken-index"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":token": {
				S: aws.String(confirmationToken),
			},
		},
		KeyConditionExpression: aws.String("ConfirmationToken = :token"),
	}

	resultChecker, err := svc.Query(input)
	if err != nil {
		errQuerying := model.BaseResponse{
			Status:  500,
			Message: "Internal Server Error: Err Querying | ",
		}
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
				errQuerying.Message += "ErrCodeProvisionedThroughputExceededException: " + aerr.Error()
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
				errQuerying.Message += "ErrCodeResourceNotFoundException: " + aerr.Error()
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
				errQuerying.Message += "ErrCodeInternalServerError: " + aerr.Error()
			default:
				fmt.Println(aerr.Error())
				errQuerying.Message += aerr.Error()
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return 500, errQuerying
	}
	users := []model.UserDetailDtoResponse{}
	err = dynamodbattribute.UnmarshalListOfMaps(resultChecker.Items, &users)
	if err != nil {
		errUnmarshal := model.BaseResponse{
			Status:  500,
			Message: "Internal Server Error: Err UnmarshalListOfMaps | " + err.Error(),
		}
		return 500, errUnmarshal
	}
	if len(users) == 0 {
		errNotFound := model.BaseResponse{
			Status:  404,
			Message: "Not Found | Token \"" + confirmationToken + "\" doesn't exist!",
		}
		return 404, errNotFound
	}

	result := model.BaseSingleResponse{
		Status:  200,
		Message: "Success find user by confirmationToken: " + confirmationToken,
		Value:   users[0],
	}
	return 200, result

}

func NewUser(user *model.User) (int, interface{}) {

	statGetByEmail, _ := FindUserByEmail(user.Email)
	if statGetByEmail == 200 {
		errUserExist := model.BaseResponse{
			Status:  400,
			Message: "Bad Request | User with email \"" + user.Email + "\" already exist!",
		}
		return 400, errUserExist
	}

	userID := uuid.NewV1().String()

	hashedPassword, err := utility.HashPassword(user.Password)
	if err != nil {
		errHashing := model.BaseResponse{
			Status:  500,
			Message: "Internal Server Error | ErrorHashing: " + err.Error(),
		}
		return 500, errHashing
	}

	nowString := strconv.Itoa(int(time.Now().UnixNano() / 1000000)) // millis
	birthDateString := strconv.Itoa(int(user.BirthDate))
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"UserId": {
				S: aws.String(userID),
			},
			"Email": {
				S: aws.String(user.Email),
			},
			"Password": {
				S: aws.String(hashedPassword),
			},
			"CreatedDate": {
				N: aws.String(nowString),
			},
			"UpdatedDate": {
				N: aws.String(nowString),
			},
			"BirthDate": {
				N: aws.String(birthDateString),
			},
			"FullName": {
				S: aws.String(user.FullName),
			},
			"Address": {
				S: aws.String(user.Address),
			},
			"PhoneNumber": {
				S: aws.String(user.PhoneNumber),
			},
			"Point": {
				N: aws.String("0"),
			},
			"PrestigeRank": {
				N: aws.String("0"),
			},
			"Status": {
				S: aws.String(status.UNCONFIRMED),
			},
			"Role": {
				S: aws.String(status.USER),
			},
			"ConfirmationToken": {
				S: aws.String(user.ConfirmationToken),
			},
			"ForgotPasswordToken": {
				S: aws.String("null"),
			},
			// "ForgetPasswordToken": {
			// 	S: aws.String(""),
			// },
			"PointStatusLogs": {
				L: []*dynamodb.AttributeValue{
				// {
				// 	M: map[string]*dynamodb.AttributeValue{},
				// },
				},
			},
		},
		ReturnConsumedCapacity: aws.String("TOTAL"),
		TableName:              aws.String(utility.Configuration.Dynamodb_dbname),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		errResponse := model.BaseResponse{}
		errResponse.Status = 500
		errResponse.Message = "Internal Server Error: "
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				fmt.Println(dynamodb.ErrCodeConditionalCheckFailedException, aerr.Error())
				errResponse.Message += dynamodb.ErrCodeConditionalCheckFailedException
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
				errResponse.Message += dynamodb.ErrCodeProvisionedThroughputExceededException
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
				errResponse.Message += dynamodb.ErrCodeResourceNotFoundException
			case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
				fmt.Println(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
				errResponse.Message += dynamodb.ErrCodeItemCollectionSizeLimitExceededException
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
				errResponse.Message += dynamodb.ErrCodeInternalServerError
			default:
				fmt.Println(aerr.Error())
			}
			errResponse.Message += " | " + aerr.Error()
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
			errResponse.Message = err.Error()
		}
		return 500, errResponse
	}
	resSuccess := model.BaseResponse{
		Status:  200,
		Message: "Success create new user, email: " + user.Email,
	}
	return 200, resSuccess
}

func FindUserByEmail(email string) (int, interface{}) {
	inputCheck := &dynamodb.QueryInput{
		TableName: aws.String(utility.Configuration.Dynamodb_dbname),
		IndexName: aws.String("Email-index"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":email": {
				S: aws.String(email),
			},
		},
		KeyConditionExpression: aws.String("Email = :email"),
	}

	resultChecker, err := svc.Query(inputCheck)
	if err != nil {
		errQuerying := model.BaseResponse{
			Status:  500,
			Message: "Internal Server Error: Err Querying | ",
		}
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
				errQuerying.Message += "ErrCodeProvisionedThroughputExceededException: " + aerr.Error()
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
				errQuerying.Message += "ErrCodeResourceNotFoundException: " + aerr.Error()
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
				errQuerying.Message += "ErrCodeInternalServerError: " + aerr.Error()
			default:
				fmt.Println(aerr.Error())
				errQuerying.Message += aerr.Error()
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return 500, errQuerying
	}
	usersChecker := []model.UserDetailDtoResponse{}
	err = dynamodbattribute.UnmarshalListOfMaps(resultChecker.Items, &usersChecker)
	if err != nil {
		errUnmarshal := model.BaseResponse{
			Status:  500,
			Message: "Internal Server Error: Err UnmarshalListOfMaps | " + err.Error(),
		}
		return 500, errUnmarshal
	}
	if len(usersChecker) == 0 {
		errNotFound := model.BaseResponse{
			Status:  404,
			Message: "Not Found | User with email \"" + email + "\" doesn't exist!",
		}
		return 404, errNotFound
	}

	result := model.BaseSingleResponse{
		Status:  200,
		Message: "Success get user by email: " + email,
		Value:   usersChecker[0],
	}
	return 200, result
}

func ScanUser() (int, interface{}) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(utility.Configuration.Dynamodb_dbname),
	}
	var users []model.UserSimpleDtoResponse
	err := svc.ScanPages(input, func(page *dynamodb.ScanOutput, last bool) bool {
		usrs := []model.UserSimpleDtoResponse{}
		err := dynamodbattribute.UnmarshalListOfMaps(page.Items, &usrs)
		if err != nil {
			panic(fmt.Sprintf("failed to unmarshal Dynamodb Scan Items, %v", err))
		}
		users = append(users, usrs...)
		return true
	})
	if err != nil {
		errScanning := model.BaseResponse{
			Status:  500,
			Message: "Internal Server Error: ErrScanning | " + err.Error() + " ",
		}
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
				errScanning.Message += "ErrCodeProvisionedThroughputExceededException: " + aerr.Error()
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
				errScanning.Message += "ErrCodeResourceNotFoundException: " + aerr.Error()
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
				errScanning.Message += "ErrCodeInternalServerError: " + aerr.Error()
			default:
				fmt.Println(aerr.Error())
				errScanning.Message += aerr.Error()
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return 500, errScanning
	}
	usTmps := make([]interface{}, len(users))
	for i, v := range users {
		usTmps[i] = v
	}
	resSuccess := model.BaseListResponse{
		Status:  200,
		Message: "Success",
		Content: usTmps,
	}
	return 200, resSuccess
}

func FindUserByID(userID string) (int, interface{}) {
	// your own db fetch here instead of user :=...
	input := &dynamodb.QueryInput{
		TableName: aws.String(utility.Configuration.Dynamodb_dbname),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":id": {
				S: aws.String(userID),
			},
		},
		KeyConditionExpression: aws.String("UserId = :id"),
	}

	result, err := svc.Query(input)
	if err != nil {
		errQuerying := model.BaseResponse{
			Status:  500,
			Message: "Internal Server Error: Err Querying | ",
		}
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
				errQuerying.Message += "ErrCodeProvisionedThroughputExceededException: " + aerr.Error()
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
				errQuerying.Message += "ErrCodeResourceNotFoundException: " + aerr.Error()
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
				errQuerying.Message += "ErrCodeInternalServerError: " + aerr.Error()
			default:
				fmt.Println(aerr.Error())
				errQuerying.Message += "" + aerr.Error()
			}
		} else {
			fmt.Println(err.Error())
		}
		return 500, errQuerying
	}

	users := []model.UserDetailDtoResponse{}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		errUnmarshal := model.BaseResponse{
			Status:  500,
			Message: "Internal Server Error: Err UnmarshalListOfMaps | failed to unmarshal Dynamodb Query Items: " + err.Error(),
		}
		return 500, errUnmarshal
	}

	if len(users) != 0 {
		succRes := model.BaseSingleResponse{
			Status:  200,
			Message: "Success get user by id: " + userID,
			Value:   users[0],
		}
		return 200, succRes
	} else {
		errNotFound := model.BaseResponse{
			Message: "Error Not Found User By Id: " + userID,
			Status:  404,
		}
		return 404, errNotFound
	}
}
