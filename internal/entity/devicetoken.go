package entity

import "time"

// DevicePlatform identifies the OS of a device a push notification token was
// registered from.
type DevicePlatform string

const (
	DevicePlatformAndroid DevicePlatform = "ANDROID"
	DevicePlatformIOS     DevicePlatform = "IOS"
)

// DeviceToken is an FCM registration token for a user's device. A user can
// have multiple tokens (one per device); a token is unique across users,
// since the same physical device's token can move to a different account
// (e.g. logout/login as someone else).
type DeviceToken struct {
	ID        string
	UserID    string
	Token     string
	Platform  DevicePlatform
	UpdatedAt time.Time
}
