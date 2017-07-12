package service

import (
	"log"

	"go-gellato-membership/utility"

	"strconv"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendNotifPoint(fullName, emailTo string, point int) {

	pointLeft := utility.Configuration.BonusPoint - point

	from := mail.NewEmail("no-reply@saintropic.com", "no-reply@saintropic.com")
	subject := strconv.Itoa(pointLeft) + " Point Left To get Free 1 Gellato From Saintropic"
	to := mail.NewEmail(fullName, emailTo)
	plainTextContent := "Hey! It's just " + strconv.Itoa(pointLeft) + " to get free 1 your own choice of Gellato from Saintropic"
	htmlContent := `
		<h3>Hey! It's just ` + strconv.Itoa(pointLeft) + ` to get free 1 of your own choice Gellato from Saintropic</h3>
		<p>
			Your current point : <strong>` + strconv.Itoa(point) + `</strong><br/>
			Thank you for stay being our loyal members. We will inform our promotions period and our new products.
		</p>
		<p>
			<small><a href="https://saintropic.com">Saintropic.com</a></small>
		</p>
	`
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(utility.Configuration.SendgridKey)
	_, err := client.Send(message)
	if err != nil {
		log.Println(err)
	}
	// else {
	// 	fmt.Println(response.StatusCode)
	// 	fmt.Println(response.Body)
	// 	fmt.Println(response.Headers)
	// }
}

func SendConfirmationMail(fullName, emailTo, token string) {

	from := mail.NewEmail("no-reply@saintropic.com", "no-reply@saintropic.com")
	subject := "Register Confirmation Saintropic"
	to := mail.NewEmail(fullName, emailTo)
	plainTextContent := "Thank you for registering at saintropic.com. Please confirm your account!"
	htmlContent := `
		<h3>Thank you for registering at saintropic.com</h3>
		<p>
			You're just one step closer to be our member. 
			Please click this following button to confirm your account.<br/>
			Token : <strong>` + token + `</strong> <br/>
			<button><a class="btn btn-default" href="#` + token + `" >Confirm Account</a></button>
		</p>
		<p>
			<small><a href="https://saintropic.com">Saintropic.com</a></small>
		</p>
	`
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(utility.Configuration.SendgridKey)
	_, err := client.Send(message)
	if err != nil {
		log.Println(err)
	}
	// else {
	// 	fmt.Println(response.StatusCode)
	// 	fmt.Println(response.Body)
	// 	fmt.Println(response.Headers)
	// }
}

func SendPasswordConfirmation(fullName, emailTo, token string) {

	from := mail.NewEmail("no-reply@saintropic.com", "no-reply@saintropic.com")
	subject := "Change Password Confirmation at Saintropic"
	to := mail.NewEmail(fullName, emailTo)
	plainTextContent := "You had requested for password change!"
	htmlContent := `
		<h3>Change Password Request</h3>
		<p>
			Ignore this message if you never attempt to do this action.<br/>
			Please click the button below to complete your password change request.<br/>
			Token : <strong>` + token + `</strong> <br/>
			<button><a class="btn btn-default" href="#` + token + `" >Change Password</a></button>
		</p>
		<footer>
			<p>
				<small><a href="https://saintropic.com">Saintropic.com</a></small>
			</p>
		</footer>
	`
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(utility.Configuration.SendgridKey)
	_, err := client.Send(message)
	if err != nil {
		log.Println(err)
	}
	// else {
	// 	fmt.Println(response.StatusCode)
	// 	fmt.Println(response.Body)
	// 	fmt.Println(response.Headers)
	// }
}
