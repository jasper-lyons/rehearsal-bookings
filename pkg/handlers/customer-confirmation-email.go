package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	da "rehearsal-bookings/pkg/data_access"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type CustomerEmailTemplateData struct {
	CustomerName  string
	Type          string
	RoomName      string
	Day           string
	StartTime     string
	EndTime       string
	CustomerPhone string
}

var customerConfirmationEmailHTMLTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Booking Confirmation - Bad Habit Studios</title>
</head>
<body style="font-family: Arial; line-height: 1.6; color: #333;">
    <p>Hey {{ .CustomerName }},</p>
    
    <p>Thanks for booking a rehearsal at <strong>Bad Habit Studios</strong>!</p>
    
    <p>Just to confirm, you've booked a <strong>{{ .Type }}</strong> rehearsal in <strong>{{ .RoomName }}</strong> on <strong>{{ .Day }}</strong> from <strong>{{ .StartTime }}</strong> to <strong>{{ .EndTime }}</strong>.</p>
    
    <p>If you need to cancel or modify your booking, please reach out to us at <a href="tel:07496983488">07496 983 488</a> or <a href="mailto:hello@badhabitstudios.co.uk">hello@badhabitstudios.co.uk</a>.</p>
    
    <p>We are a self-service studio, which means there won't always be a member of staff on-site during your session.</p>
    
    <p>On the day of your rehearsal, we will text door codes to <strong>{{ .CustomerPhone }}</strong>. For security reasons, we generally keep the code-lock enabled even when the studio is staffed, so please share these with your bandmates and make a note of them when you go in and out of the building.</p>
    
    <p>Further information about location and room specs can be found on our website at <a href="https://badhabitstudios.co.uk" target="_blank">badhabitstudios.co.uk</a>.</p>
    
	<p>Please note that we have a 24 hour cancellation policy. More details can be found <a href="https://badhabitstudios.co.uk/cancellation-policy" target="_blank">here</a>.</p>
    <p>Any questions or issues, please reach out!</p>
    
    <p>Cheers,<br>
    <strong>Bad Habit Studios</strong></p>
    
    <table>
        <tr>
            <td>
                <img src="https://badhabiteastbourne.co.uk/wp-content/uploads/2024/09/bad-habit-logo-no-background.webp" width="80" height="80" alt="Bad Habit Studios Logo">
            </td>
            <td style="padding-left: 10px;">
                <p><a href="tel:07496983488">07496 983 488</a></p> 
                <p><a href="mailto:hello@badhabitstudios.co.uk">hello@badhabitstudios.co.uk</a></p>
                <p>
                    <a style="color: transparent;" href="https://facebook.com/badhabiteastbourne" target="_blank">
                        <img src="https://cdn-icons-png.flaticon.com/24/733/733547.png" alt="Facebook" width="24" height="24">
                    </a>
                    <a style="color: transparent;" href="https://instagram.com/badhabiteastbourne" target="_blank">
                        <img src="https://cdn-icons-png.flaticon.com/24/733/733558.png" alt="Instagram" width="24" height="24">
                    </a>
                    <a style="color: transparent;" href="https://wa.me/7496983488" target="_blank">
                        <img src="https://cdn-icons-png.flaticon.com/24/733/733585.png" alt="WhatsApp" width="24" height="24">
                    </a>
                </p>
            </td>
        </tr>
    </table>
    
    <p style="font-size: 12px; color: #777;">This is an automated email, so if you reply to it, we might not be able to get back to you.</p>
</body>
</html>
`

var customerConfirmationEmailTemplatePlainText = `
Hey {{ .CustomerName }},

Thanks for booking a rehearsal at Bad Habit Studios!

Just to confirm you've booked a {{ .Type}} rehearsal {{ .RoomName }} on {{ .Day }} from {{ .StartTime }} to {{ .EndTime }}. 

If you need to cancel or modify your booking, please reach out to us on 07496 983 488 or hello@badhabitstudios.co.uk.

We are a self-service studio, which means there won't always be a member of staff on-site during your session.

On the day of your rehearsal, we will text door codes to {{ .CustomerPhone }}. For security reasons we generally keep the code-lock enabled even when the studio is staffed, so please share these with your bandmates and make a note of them when you go in and out of the building.

Further information about location and room specs can be found on our website at https://badhabitstudios.co.uk
Please note that we have a 24 hour cancellation policy. More details can be found at https://badhabitstudios.co.uk/cancellation-policy.
Any questions or issues, please reach out!

Cheers,
Bad Habit Studios

This is an automated email so if you reply to it, we might not be able to get back to you.
`

func SendCustomerBookingConfirmationEmail(booking *da.Booking) error {
	from := mail.NewEmail("Rehearsal Booking", os.Getenv("TRANSACTIONAL_FROM_ADDRESS"))
	subject := "Bad Habit Studios: Rehearsal Booking Confirmation"
	to := mail.NewEmail(booking.CustomerName, booking.CustomerEmail)

	// Prepare the template data
	templateData := CustomerEmailTemplateData{
		CustomerName:  booking.CustomerName,
		Type:          booking.Type,
		RoomName:      booking.RoomName,
		Day:           booking.StartTime.Format("Monday, 2 January, 2006"),
		StartTime:     booking.StartTime.Format("3:04 PM"),
		EndTime:       booking.EndTime.Format("3:04 PM"),
		CustomerPhone: booking.CustomerPhone,
	}

	// Attempt to generate the HTML content
	var htmlContent bytes.Buffer
	htmlTmpl, err := template.New("confirmation-email-html").Parse(customerConfirmationEmailHTMLTemplate)
	if err == nil {
		// Only execute the template if parsing succeeded
		err = htmlTmpl.Execute(&htmlContent, templateData)
	}
	if err != nil {
		fmt.Printf("Failed to generate HTML email content: %v. Falling back to plain text.\n", err)
	}

	// Execute the plain text template with our data
	var plainTextContent bytes.Buffer
	plainTextTmpl, errPlain := template.New("confirmation-email-plain-text").Parse(customerConfirmationEmailTemplatePlainText)
	if errPlain != nil {
		return fmt.Errorf("Failed to parse email template: %v", errPlain)
	}

	if err := plainTextTmpl.Execute(&plainTextContent, templateData); err != nil {
		return fmt.Errorf("Failed to execute email template: %v", err)
	}

	message := mail.NewSingleEmail(from, subject, to, plainTextContent.String(), htmlContent.String())
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("Failed to send confirmation email: %v", err)
	}

	// Log the response for monitoring purposes
	if response.StatusCode >= 400 {
		return fmt.Errorf("Failed to send confirmation email. Status: %d, Body: %s", response.StatusCode, response.Body)
	}

	return nil
}
