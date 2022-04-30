package web

import (
	"entry_task/mysqldb"
	"entry_task/tcp"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
)

func getpath(tag string) string {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if strings.Contains(pwd, tag) == false {
		pwd = pwd + "/" + tag
	}
	return pwd
}

var web_dir = getpath("web")
var login_path = web_dir + "/login_page.html"
var profile_path = web_dir + "/profile_page.html"
var photo_path = web_dir + "/../photos"

var userMap = make(map[string]*mysqldb.User)
var userMapLock = &sync.RWMutex{}

func StartWEBServer() {
	// if photo dir is start with '/' like /photos/123.jpeg, use this
	// http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./"))))

	http.Handle("/profile_page/photos/", http.StripPrefix("/profile_page/photos/", http.FileServer(http.Dir("./photos"))))
	http.HandleFunc("/login_page/", MakeHandler(LoginPageHandler))
	http.HandleFunc("/profile_page/", MakeHandler(ProfilePageHandler))
	http.HandleFunc("/login/", MakeHandler(LoginHandler))
	http.HandleFunc("/nickname/", MakeHandler(NicknameHandler))
	http.HandleFunc("/photo/", MakeHandler(PhotoHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request, Acc string) {
	renderTemplate(w, "login_page", nil)
}

func ProfilePageHandler(w http.ResponseWriter, r *http.Request, Acc string) {
	userMapLock.Lock()
	defer userMapLock.Unlock()

	if user, ok := userMap[Acc]; ok {
		renderTemplate(w, "profile_page", user)
		return
	}
	http.Redirect(w, r, "/login_page/", http.StatusFound)
}

func LoginHandler(w http.ResponseWriter, r *http.Request, Acc string) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	user := &mysqldb.User{
		Acc:      username,
		Pwd:      password,
		Nickname: "",
		Photo:    "photos/init.jpeg",
		Id:       -1,
	}
	user, err := tcp.LoginRPC(user)
	// login fails
	if err != nil {
		fmt.Println("[WEB Server LoginHandler] Login error: ", err)
		http.Redirect(w, r, "/login_page/", http.StatusFound)
		return
	}
	// record this user
	userMapLock.Lock()
	defer userMapLock.Unlock()
	userMap[user.Acc] = user
	http.Redirect(w, r, "/profile_page/"+user.Acc, http.StatusFound)
}

func NicknameHandler(w http.ResponseWriter, r *http.Request, Acc string) {
	if user, ok := userMap[Acc]; ok {
		oldNickname := user.Nickname
		user.Nickname = r.FormValue("nickname")
		err := tcp.NicknameRPC(user)
		if err != nil {
			fmt.Println("[WEB Server NicknameHandler] Update nickname error: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			user.Nickname = oldNickname
			return
		}
		userMap[Acc] = user
		http.Redirect(w, r, "/profile_page/"+user.Acc, http.StatusFound)
		return
	}
	http.Redirect(w, r, "/login_page/", http.StatusFound)
}

func PhotoHandler(w http.ResponseWriter, r *http.Request, Acc string) {
	if user, ok := userMap[Acc]; ok {
		oldPhoto := user.Photo
		file, _, err := r.FormFile("photo")
		if err != nil {
			fmt.Println("[WEB Server PhotoHandler] Error Retrieving the File: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			// http.Redirect(w, r, "/profile_page/"+user.Acc, http.StatusFound)
			return
		}
		defer file.Close()

		// save new photo and delete the old one
		tempFile, err := ioutil.TempFile(photo_path, user.Acc+"-*.jpeg")
		if err != nil {
			fmt.Println("[WEB Server PhotoHandler] Error opening tmp file: ", err)
		}
		defer tempFile.Close()

		tmp_path_arr := strings.Split(tempFile.Name(), "/")
		path_len := len(tmp_path_arr)
		user.Photo = tmp_path_arr[path_len-2] + "/" + tmp_path_arr[path_len-1]
		// user.Photo = tempFile.Name()
		err = tcp.PhotoRPC(user)
		if err != nil {
			fmt.Println("[WEB Server PhotoHandler] Update photo error: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			user.Photo = oldPhoto
			return
		}

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println("[WEB Server PhotoHandler] Read file error: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			user.Photo = oldPhoto
			return
		}
		tempFile.Write(fileBytes)

		if oldPhoto != "" && oldPhoto != "photos/init.jpeg" {
			err = os.Remove(oldPhoto)
			if err != nil {
				fmt.Println("[WEB Server PhotoHandler] Delete old photo error: ", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				user.Photo = oldPhoto
				return
			}
		}

		userMap[Acc] = user
		http.Redirect(w, r, "/profile_page/"+user.Acc, http.StatusFound)
		return
	}
	http.Redirect(w, r, "/login_page/", http.StatusFound)
}

var validPath = regexp.MustCompile("^(/(login_page|login)/$)|^(/(profile_page|nickname|photo)/([a-zA-Z0-9]+)$)")

func MakeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		// if r.URL.Path == "/profile_page/photos/init.jpeg" {
		//     r.URL.Path = "photos/init.jpeg"
		// }
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[5])
	}
}

var templates = template.Must(template.ParseFiles(login_path, profile_path))

func renderTemplate(w http.ResponseWriter, tmpl string, user *mysqldb.User) {
	err := templates.ExecuteTemplate(w, tmpl+".html", user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
