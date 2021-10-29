package http_server

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
)

func UploaderHandler(w http.ResponseWriter, req *http.Request) {
	userId := req.FormValue("user_id")
	if userId == "" {
		http.Error(w, "required user_id", http.StatusBadRequest)
		return
	}
	file, header, err := req.FormFile("avatarFile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	filename := path.Join("web/avatars", userId + path.Ext(header.Filename))
	err = ioutil.WriteFile(filename, data, 0777)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = io.WriteString(w, "Successful")
}
