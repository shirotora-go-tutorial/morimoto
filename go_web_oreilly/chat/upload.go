package main

import (
	"net/http"
	"io"
	"io/ioutil"
	"path/filepath"
)

func uploaderHandle(w http.ResponseWriter, req *http.Request){
	userId := req.FormValue("userid")
	file, header, err := req.FormFile("avatarFile")
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	filename := filepath.Join("avatars", userId + filepath.Ext(header.Filename))
	err = ioutil.WriteFile(filename, data, 0600)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	io.WriteString(w, "success")
}
