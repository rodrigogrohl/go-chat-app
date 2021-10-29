package chat

import (
	"github.com/rodrigogrohl/go-chat-app/configs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar Avatar = AuthAvatar{}
	client := new(Client)
	avatarURL, err := authAvatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL {
		t.Error("Avatar.getAvatarURL should return ErrNoAvatarURL when has no value.")
	}

	testAvatarURL := "https://avatars.githubusercontent.com/u/5113606?v=4"
	client.UserData = map[string]interface{}{
		"avatar_url": testAvatarURL,
	}
	avatarURL, err = authAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("Avatar.getAvatarURL should not return any error")
	}
	if avatarURL != testAvatarURL {
		t.Error("Avatar.GetAvatarURL should return correct URL")
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	client := new(Client)
	client.UserData = map[string]interface{}{"user_id": "69c32ba7a0b6c26c1f110d9003ed3a49"}

	url, err := gravatarAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarUrl should not return an error")
	}
	if url != "//www.gravatar.com/avatar/69c32ba7a0b6c26c1f110d9003ed3a49" {
		t.Errorf("GravatarAvatar.GetAvatarURL wrongly returned '%s'", url)
	}
}

func TestFileSystemAvatar(t *testing.T) {
	filename := filepath.Join(configs.StoragePathAvatars, "abc.jpg")
	_ = ioutil.WriteFile(filename, []byte{}, 0777)
	defer os.Remove(filename)
	var fileSystemAvatar FileSystemAvatar
	client := new(Client)
	client.UserData = map[string]interface{}{"user_id": "abc"}
	url, err := fileSystemAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("FileSystemAvatar.GetAvatarURL should not return an error")
	}
	if url != "/avatars/abc.jpg" {
		t.Errorf("FileSystemAvatar.GetAvatarURL wrongly returned %s", url)
	}
}