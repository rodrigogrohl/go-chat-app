package chat

import (
	"errors"
	"github.com/rodrigogrohl/go-chat-app/configs"
	"io/ioutil"
	"path"
)

var ErrNoAvatarURL = errors.New("chat: unable to get an avatar URL")

type Avatar interface {
	GetAvatarURL(c *Client) (string, error)
}

// AuthAvatar retrieve information from Authentication provider
type AuthAvatar struct {}
var UseAuthAvatar AuthAvatar

func (AuthAvatar) GetAvatarURL(c *Client) (string, error) {
	if url, ok := c.UserData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

// GravatarAvatar retrieve user picture from gravatar.com
type GravatarAvatar struct {}
var UseGravatarAvatar GravatarAvatar

func (GravatarAvatar) GetAvatarURL(c *Client) (string, error) {
	if email, ok := c.UserData["user_id"]; ok {
		if userIdStr, ok := email.(string); ok {
			return "//www.gravatar.com/avatar/" + userIdStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

type FileSystemAvatar struct {}
var UseFileSystemAvatar FileSystemAvatar

func (FileSystemAvatar) GetAvatarURL(c *Client) (string, error) {
	if userId, ok := c.UserData["user_id"]; ok {
		if userIdStr, ok := userId.(string); ok {
			files, err := ioutil.ReadDir(configs.StoragePathAvatars)
			if err != nil {
				return "", ErrNoAvatarURL
			}
			for _, file := range files {
				if file.IsDir() {
					continue
				}
				if match, _ := path.Match(userIdStr + "*", file.Name()); match {
					return configs.HttpAvatarGet + file.Name(), nil
				}
			}
			return configs.HttpAvatarGet + userIdStr + ".png", nil
		}
	}
	return "", ErrNoAvatarURL
}
