// Package fcm sends push notifications via Firebase Cloud Messaging.
package fcm

import (
	"context"
	"log/slog"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"

	"github.com/fark-tee/fark-tee-backend/internal/config"
)

// Client sends push notifications to a single device via FCM.
type Client interface {
	// SendNudge notifies token that fromDisplayName is waiting on the
	// recipient in meetupID. Android gets a data-only, high-priority message
	// so the client's background handler can build its own full-screen local
	// notification; iOS gets a normal displayable alert, since a terminated
	// iOS app can't be woken the same way.
	SendNudge(ctx context.Context, token, meetupID, fromUserID, fromDisplayName string) error

	// SendCheckInRequest notifies token that fromDisplayName wants to know
	// whether the recipient (currently heading home) is okay, so the client
	// can show a full-screen prompt with "okay"/"not okay" actions. Same
	// delivery shape as SendNudge.
	SendCheckInRequest(ctx context.Context, token, meetupID, fromUserID, fromDisplayName string) error

	// SendCheckInEmergencyAlert notifies token that fromDisplayName answered
	// "not okay" to a check-in, carrying their emergency contact so the
	// recipient's full-screen alert can offer a tap-to-call action. Sent to
	// every other accepted party member except fromUserID.
	SendCheckInEmergencyAlert(ctx context.Context, token, meetupID, fromUserID, fromDisplayName, emergencyContactName, emergencyContactPhone string) error
}

type clientImpl struct {
	// messaging is nil when no Firebase credentials are configured yet, in
	// which case SendNudge logs and no-ops instead of sending.
	messaging *messaging.Client
}

// @WireSet("Infrastructure")
func New(ctx context.Context, cfg *config.Config) (Client, error) {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, err
	}

	msgClient, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}

	return &clientImpl{messaging: msgClient}, nil
}

func (c *clientImpl) SendNudge(ctx context.Context, token, meetupID, fromUserID, fromDisplayName string) error {
	if c.messaging == nil {
		slog.Warn("FCM not configured, skipping nudge notification", slog.String("token", token))

		return nil
	}

	message := &messaging.Message{
		Token: token,
		Data: map[string]string{
			"type":            "nudge",
			"meetupId":        meetupID,
			"fromUserId":      fromUserID,
			"fromDisplayName": fromDisplayName,
		},
		Android: &messaging.AndroidConfig{
			Priority: "high",
		},
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{"apns-priority": "10"},
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Alert: &messaging.ApsAlert{
						Title: "รีบมาได้แล้ว!",
						Body:  fromDisplayName + " ตามคุณแล้ว!!",
					},
					Sound:            "default",
					ContentAvailable: true,
				},
			},
		},
	}

	_, err := c.messaging.Send(ctx, message)

	return err
}

func (c *clientImpl) SendCheckInEmergencyAlert(ctx context.Context, token, meetupID, fromUserID, fromDisplayName, emergencyContactName, emergencyContactPhone string) error {
	if c.messaging == nil {
		slog.Warn("FCM not configured, skipping check-in emergency alert", slog.String("token", token))

		return nil
	}

	message := &messaging.Message{
		Token: token,
		Data: map[string]string{
			"type":                  "checkin_emergency",
			"meetupId":              meetupID,
			"fromUserId":            fromUserID,
			"fromDisplayName":       fromDisplayName,
			"emergencyContactName":  emergencyContactName,
			"emergencyContactPhone": emergencyContactPhone,
		},
		Android: &messaging.AndroidConfig{
			Priority: "high",
		},
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{"apns-priority": "10"},
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Alert: &messaging.ApsAlert{
						Title: fromDisplayName + " ตอบว่าไม่โอเค!",
						Body:  "แตะเพื่อดูเบอร์ติดต่อฉุกเฉิน",
					},
					Sound:            "default",
					ContentAvailable: true,
				},
			},
		},
	}

	_, err := c.messaging.Send(ctx, message)

	return err
}

func (c *clientImpl) SendCheckInRequest(ctx context.Context, token, meetupID, fromUserID, fromDisplayName string) error {
	if c.messaging == nil {
		slog.Warn("FCM not configured, skipping check-in request notification", slog.String("token", token))

		return nil
	}

	message := &messaging.Message{
		Token: token,
		Data: map[string]string{
			"type":            "checkin_request",
			"meetupId":        meetupID,
			"fromUserId":      fromUserID,
			"fromDisplayName": fromDisplayName,
		},
		Android: &messaging.AndroidConfig{
			Priority: "high",
		},
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{"apns-priority": "10"},
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Alert: &messaging.ApsAlert{
						Title: "โอเคนะคะ?",
						Body:  fromDisplayName + " อยากรู้ว่าคุณโอเคไหม",
					},
					Sound:            "default",
					ContentAvailable: true,
				},
			},
		},
	}

	_, err := c.messaging.Send(ctx, message)

	return err
}
