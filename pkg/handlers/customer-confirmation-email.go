package handlers

import (
	"fmt"
	"os"
	da "rehearsal-bookings/pkg/data_access"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"html/template"
	"bytes"
)

type CustomerEmailTemplateData struct {
	RoomName string
	Day string
	StartTime string
	EndTime string
}

var customerConfirmationEmailTemplate = `
Just to confirm you've booked {{.RoomName}} at Bad Habit Studios on the {{.Day}} from
{{.StartTime}} to {{.EndTime}}.

If you need to cancel or modify the booking reach out to badhabitstudioseb@gmail.com

This is an automated email so if you reply to it, we might not be able to get back to you.
`

func SendCustomerBookingConfirmationEmail(booking *da.Booking) error {
	from := mail.NewEmail("Rehearsal Booking", os.Getenv("TRANSACTIONAL_FROM_ADDRESS"))
	subject := "Bad Habit Studios: Rehearsal Booking Confirmation"
	to := mail.NewEmail(booking.CustomerName, booking.CustomerEmail)
	
	// Prepare the template data
	templateData := CustomerEmailTemplateData {
		RoomName: booking.RoomName,
		Day: booking.StartTime.Format("Monday, 2nd January, 2006"),
		StartTime: booking.StartTime.Format("3:04 PM"),
		EndTime: booking.EndTime.Format("3:04 PM"),
	}

	
	tmpl, err := template.New("confirmation-email").Parse(customerConfirmationEmailTemplate)
	if err != nil {
		return fmt.Errorf("Failed to parse email template: %v", err)
	}
	
	// Execute the template with our data
	var emailContent bytes.Buffer
	if err := tmpl.Execute(&emailContent, templateData); err != nil {
		return fmt.Errorf("Failed to execute email template: %v", err)
	}
	
	message := mail.NewSingleEmail(from, subject, to, emailContent.String(), "")
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
