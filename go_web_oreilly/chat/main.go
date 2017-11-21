package main

import (
	"../trace"
	"encoding/json"
	"flag"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"
)

var avatars Avatar = TryAvatars{
	UseFileSystemAvatar,
	UseAuthAvatar,
	UseGravatar,
}

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(
			template.ParseFiles(filepath.Join(
				"templates",
				t.filename,
			)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

type Authkey struct {
	Key          string `json:"key"`
	Clientid     string `json:"clientid"`
	Clientsecret string `json:"clientsecret"`
	Redirecturl  string `json:"redirecturl"`
}

func main() {
	var addr = flag.String("addr", ":8080", "アドレス")
	flag.Parse()

	keyfile, err := ioutil.ReadFile("./key.json")
	if err != nil {
		log.Fatal("keyfile cannot not found")
	}

	var authkey []Authkey
	err = json.Unmarshal(keyfile, &authkey)
	if err != nil {
		log.Fatal("Format Error: ", err)
	}

	gomniauth.SetSecurityKey(authkey[0].Key)
	gomniauth.WithProviders(
		google.New(
			authkey[0].Clientid,
			authkey[0].Clientsecret,
			authkey[0].Redirecturl,
		),
	)

	//r := newRoom(UseAuthAvatar)
	//r := newRoom(UseGravatar)
	r := newRoom(UseFileSystemAvatar)
	r.tracer = trace.New(os.Stdout)
	http.Handle("/", MustAuth(&templateHandler{filename: "./chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/uploader", uploaderHandle)
	http.Handle("/avatars/",
		http.StripPrefix("/avatars/",
			http.FileServer(http.Dir("./avatars"))))

	go r.run()

	log.Println("Start Web  server Port:", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServer", err)
	}
}
