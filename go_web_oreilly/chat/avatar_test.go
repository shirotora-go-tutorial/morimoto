package main
import (
	"testing"
	"path/filepath"
	"io/ioutil"
	"os"
)

func TestAutuAvatar(t *testing.T){
	var authAvatar AuthAvatar
	client := new(client)
	url, err := authAvatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL{
		t.Error("if dont have value AuthAvatar.GetAvatarURL should return ErrNoAvatarURL")
	}

	testURL  := "http://url-to-avatar/"
	client.userData = map[string]interface{}{"avatar_url": testURL}
	url, err = authAvatar.GetAvatarURL(client)
	if err != nil{
		t.Error("if it have value AuthAvatar.GetAvatarURL should not return error", err, url)
	}else {
		if url != testURL {
			t.Error("AuthAvatar.GetAvatarURL must return correct URL", err)
		}
	}
}

func TestGravatarAvatar(t *testing.T){
	var gravatarAvatar GravatarAvatar
	client := new(client)
	//client.userData = map[string]interface{}{"email" : "MyEmailAddress@example.com"}
	client.userData = map[string]interface{}{"userid" : "0bc83cb571cd1c50ba6f3e8a78f134"}
	url, err := gravatarAvatar.GetAvatarURL(client)
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
	client := new(client)
	client.userData = map[string]interface{}{"userid" : "abc"}
	url, err := fileSystemAvatar.GetAvatarURL(client)
	if err != nil{
		t.Error("FileSystemAvatar.GetAavtarURL should return error")
	}
	if url != "/avatars/abc.jpg" {
		t.Error("FileSystemAvatar.GetAvatarURL return %s that is incoreect", url)
	}
}