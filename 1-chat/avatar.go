package main

import (
	"errors"
)

// ErrNoAvatarURL is an error when no avatar URL is available
var ErrNoAvatarURL = errors.New("chat: Unable to get avatar URL")

// Avatar is an interface to get avatar URL
type Avatar interface {
	GetAvatarURL(c *client) (string, error)
}

// AuthAvatar represent an avatar that we get from oAuth provider
type AuthAvatar struct{}

// UseAuthAvatar is a zero initialization for AuthAvatar
var UseAuthAvatar AuthAvatar

// GetAvatarURL return avatar URL
func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	url, ok := c.userData["avatar_url"]
	if !ok {
		return "", ErrNoAvatarURL
	}
	urlStr, ok := url.(string)
	if !ok {
		return "", ErrNoAvatarURL
	}
	return urlStr, nil
}
