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

type AdminEmailTemplateData struct {
	CustomerName string
	RoomName string
	Day string
	StartTime string
	EndTime string
}

var adminBookingNotificationEmailTemplate = `
New Booking for: {{.CustomerName}}
Room: {{.RoomName}}
Day: {{.Day}}
Times: {{.StartTime}} to {{.EndTime}}
`

func SendAdminBookingNotificationEmail(booking *da.Booking) error {
	from := mail.NewEmail("Rehearsal Booking", os.Getenv("TRANSACTIONAL_FROM_ADDRESS"))
	subject := "New Booking"
	to := mail.NewEmail("Admin", os.Getenv("ADMIN_NOTIFICATION_EMAIL"))
	
	// Prepare the template data
	templateData := AdminEmailTemplateData {
		CustomerName: booking.CustomerName,
		RoomName: booking.RoomName,
		Day: booking.StartTime.Format("Monday, 2nd January, 2006"),
		StartTime: booking.StartTime.Format("3:04 PM"),
		EndTime: booking.EndTime.Format("3:04 PM"),
	}

	
	tmpl, err := template.New("notification-email").Parse(adminBookingNotificationEmailTemplate)
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
		return fmt.Errorf("Failed to send admin-notification email: %v", err)
	}
	
	// Log the response for monitoring purposes
	if response.StatusCode >= 400 {
		return fmt.Errorf("Failed to send admin-notification email. Status: %d, Body: %s", response.StatusCode, response.Body)
	}
	
	return nil
}
