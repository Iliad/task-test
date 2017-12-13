package main

import (
	"encoding/json"
	"net/http"
	"reflect"
	"log"
	"math"
	"strconv"
	"fmt"
	"flag"
	"sync"

	"github.com/Iliad/task-test/model"
	"github.com/gorilla/mux"

)

var (
	appName = "Task service" // название сервиса
	version = "1.0" // версия
	date    = "2017-12-13" // дата сборки
	host    = ":8080" // адрес сервера и порт
)

var (
	auth map[string]string
	serviceMutex sync.Mutex
)

//Безопасное получение и установка пароля
func setPassword(name string, password string) {
	serviceMutex.Lock()
	defer serviceMutex.Unlock()
	auth[name] = password
}

func getPassword(name string) string {
	serviceMutex.Lock()
	defer serviceMutex.Unlock()
	return auth[name]
}

func main() {
	flag.StringVar(&host, "host", host, "Main server host name")
	flag.Parse()

	r := mux.NewRouter()

	r.HandleFunc("/", mainPage)
	r.HandleFunc("/login", login)
	r.HandleFunc("/password", changePassword)
	r.HandleFunc("/task", doWork)

	http.Handle("/", r)

	model.GormInit();
	defer model.GormClose();

	log.Printf("%s %s (%s) is starting on host: %s", appName, version, date, host)

	if err := http.ListenAndServe(host, r); err!=nil {
		log.Fatal("Can't start server. Error: ", err)
	}
}

func init() {
	auth = make(map[string]string)
}

//Хостим главную страницу
func mainPage(w http.ResponseWriter, r *http.Request) {
	if r.Method!="GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		w.Write([]byte(`<!DOCTYPE html>
		<html>
		<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<meta name="theme-color" content="#375EAB">

			<title>main page</title>
		</head>
		<body>
			Page body and some more content
		</body>
		</html>`))
	}
}

//Проверка логина пользователя
func checkLogin (login string, password string) error {
	if getPassword(login) == password {
		return nil
	}
	user := &model.User{}
	err := user.Get(login, password)
	if err == nil {
		setPassword(login, password)
		return nil
	} else {
		return err
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method!="POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		login := r.FormValue("login");
		password := r.FormValue("pass");
		err := checkLogin(login , password)
		if err != nil {
			log.Printf("Login error (User: %s): %s", login, err)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		log.Printf("Login succesfull (User: %s)", login)
		w.WriteHeader(http.StatusOK)
	}
}

func changePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method!="POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		login := r.FormValue("login")
		password := r.FormValue("pass")
		newPassword := r.FormValue("newPass")

		if newPassword != "" {
			err := checkLogin(login , password)
			if err != nil {
				log.Printf("Login error (User: %s): %s", login, err)
				w.WriteHeader(http.StatusForbidden)
				return
			}
			log.Printf("Login succesfull (User: %s)", login)

			user := &model.User{}
			err = user.ChangePassword(login, newPassword)
			if err != nil {
				log.Printf("Error changing password (User: %s): %s", login, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			setPassword(login, newPassword)

			log.Printf("Password changed (User: %s)", login)
			w.WriteHeader(http.StatusOK)
		} else {
			log.Printf("Error changing password (User: %s): %s", login, "empty new password")
			w.WriteHeader(http.StatusBadRequest)
		}

	}
}

type Values struct {
	Values []interface{} `json:"Values"`
}

func doWork(w http.ResponseWriter, r *http.Request) {
	if r.Method!="POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		login := r.FormValue("login")
		password := r.FormValue("pass")

		err := checkLogin(login, password)
		if err != nil {
			log.Printf("Login error (User: %s): %s", login, err)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		log.Printf("Login succesfull (User: %s)", login)

		var values Values

		err = json.Unmarshal([]byte(r.FormValue("value")), &values)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		for _, value := range values.Values {
			w.Write(reverse(value))
		}
		log.Printf("Work complete (User: %s)", login)
	}
}

//Разворот строк и
func reverse(value interface{}) []byte {
	switch reflect.TypeOf(value).String() {
	//Все числа после декодирования становятся float64
	case "float64":
		s := strconv.FormatFloat(math.MaxFloat64-value.(float64), 'E', -1, 64)
		return []byte(fmt.Sprintln(s))
	//Все строки после декодирования становятся string
	case "string":
		runes := []rune(value.(string))
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return []byte(fmt.Sprintln(string(runes)))
	//Все остальные типы данных не обрабатываются
	default:
		log.Println("Unsupported data type")
	}
	return nil
}
