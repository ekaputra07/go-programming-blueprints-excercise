package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"strings"
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

// GravatarAvatar represent a way to get avatar URL from Gravatar service
type GravatarAvatar struct{}

// UseGravatar is a zero initialization for GravatarAvatar
var UseGravatar GravatarAvatar

// GetAvatarURL return avatar URL from gravatar
func (GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	email, ok := c.userData["email"]
	if !ok {
		return "", ErrNoAvatarURL
	}
	emailStr, ok := email.(string)
	if !ok {
		return "", ErrNoAvatarURL
	}
	m := md5.New()
	io.WriteString(m, strings.ToLower(emailStr))
	return fmt.Sprintf("//www.gravatar.com/avatar/%x", m.Sum(nil)), nil
}
