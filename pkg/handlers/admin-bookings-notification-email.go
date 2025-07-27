package handlers

import (
	"fmt"
	"context"
	"os"
	"html/template"
	"bytes"
	da "rehearsal-bookings/pkg/data_access"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
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
	ctx := context.Background()

	sesClient, err := createSESClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to created SES client: %v", err)
	}

	from := os.Getenv("TRANSACTIONAL_FROM_ADDRESS")
	subject := "New Booking"
	to := os.Getenv("ADMIN_NOTIFICATION_EMAIL")
	
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
	
	input := &ses.SendEmailInput {
		Source: aws.String(from),
		Destination: &types.Destination {
			ToAddresses: []string{ to },
		},
		Message: &types.Message {
			Subject: &types.Content {
				Data: aws.String(subject),
			},
			Body: &types.Body {
				Text: &types.Content {
					Data: aws.String(emailContent.String()),
				},
			},
		},
	}

	// Send the email
	_, err = sesClient.SendEmail(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to send adming notification email: %v", err)
	}

	return nil
}
