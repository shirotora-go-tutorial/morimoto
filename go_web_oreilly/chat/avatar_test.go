package main
import (
	"testing"
	"path/filepath"
	"io/ioutil"
	"os"
	gomniaututest "github.com/stretchr/gomniauth/test"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	testUser := &gomniaututest.TestUser{}
	testUser.On("AvatardURL").Return("", ErrNoAvatarURL)
	testChatUser := &chatUser{User: testUser}
	url, err := authAvatar.GetAvatarURL(testChatUser)
	if err != ErrNoAvatarURL {
		t.Error("if it has value AvatardURL.GetAvatarURL mustnot return Error", err)
	}
	testUrl := "http://url-to-avatar/"
	testUser = &gomniaututest.TestUser{}
	testChatUser.User = testUser
	testUser.On("AvatarURL").Return(testUrl, nil)
	url, err = authAvatar.GetAvatarURL(testChatUser)
	if err != nil {
		t.Error("if value exist AvatardURL.GetAvatarURL must not return err", err)
	} else {
		if url != testUrl {
			t.Error("AvatardURL.GetAvatarURL must return correct url", url)
		}
	}
}


func TestGravatarAvatar(t *testing.T){
	var gravatarAvatar GravatarAvatar
	user := &chatUser{uniqueID: "abc"}
	url, err := gravatarAvatar.GetAvatarURL(user)
	if err != nil {
		t.Error("GravatarAvatar should not return error")
	}
	//if url != "//www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346" {
	if url != "//www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78f134" {
		t.Errorf("GravatarAvatar.GetAvatar return missed value %s", url)
	}
}

func TestFileSystemAvatar(t *testing.T){
	filename := filepath.Join("avatars", "abc.jpg")
	ioutil.WriteFile(filename, []byte{}, 0777)
	defer func() {
		os.Remove(filename)
	}()

	var fileSystemAvatar FileSystemAvatar
	user := &chatUser{uniqueID: "abc"}
	url, err := gravatarAvatar.GetAvatarURL(user)
	if err != nil{
		t.Error("FileSystemAvatar.GetAavtarURL should return error")
	}
	if url != "/avatars/abc.jpg" {
		t.Error("FileSystemAvatar.GetAvatarURL return %s that is incoreect", url)
	}
}