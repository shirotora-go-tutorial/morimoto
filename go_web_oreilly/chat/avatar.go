package main

import (
	"errors"
	"io/ioutil"
	"path/filepath"
)

var ErrNoAvatarURL = errors.New("chat: can t get avatar url")

type Avatar interface {
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}
var UseAuthAvatar AuthAvatar
func (_ AuthAvatar) GetAvatarURL(u ChatUser) (string, error){
	url := u.AvatarURL()
	if url != ""{
		return url, nil
	}
	return "", ErrNoAvatarURL
}

type GravatarAvatar struct {}
var UseGravatar GravatarAvatar
func (_ GravatarAvatar) GetAvatarURL(u ChatUser) (string, error){
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

type FileSystemAvatar struct{}
var UseFileSystemAvatar FileSystemAvatar
func (_ FileSystemAvatar) GetAvatarURL(c *client)(string, error){
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok{
			if files, err := ioutil.ReadDir("avatars"); err == nil{
				for _, file := range files{
					if file.IsDir(){
						continue
					}
					if match, _ := filepath.Match(useridStr + "*", file.Name());
						match {
							return "/avatars/" + file.Name(), nil
					}
				}
			}
		}
	}
	return "", ErrNoAvatarURL
}

