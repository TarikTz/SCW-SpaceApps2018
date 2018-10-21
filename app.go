package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// App : struktura
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Initialize : inacijalizacija baze
func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("%s:%s@tcp(127.0.0.1:8889)/%s", user, password, dbname)

	var err error
	a.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// Run : pokretanje servera
func (a *App) Run(addr string) {
	// allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With"})
	// allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	// allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	// log.Fatal(http.ListenAndServe(addr, handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(a.Router)))
	log.Println("Running server at localhost:" + addr)
	log.Fatal(http.ListenAndServe(addr, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(a.Router)))
}

func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	u := user{ID: id}
	if err := u.getUser(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "applications/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) getUsers(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("Start"))

	if count > 10 || count < 1 {
		count = 100
	}
	if start < 0 {
		start = 0
	}
	users, err := getUsers(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, users)
}

func (a *App) createUser(w http.ResponseWriter, r *http.Request) {
	var u user
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := u.checkEmail(a.DB); err == nil {
		respondWithError(w, http.StatusBadRequest, "This email is already in use")
		return
	}

	if err := u.createUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, u)
}

func (a *App) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var u user
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	u.ID = id

	if err := u.updateUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	u := user{ID: id}
	if err := u.deleteUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// Subjects
func (a *App) getSubject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Subject not found")
		return
	}

	s := subject{ID: id}
	if err := s.getSubject(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Subject not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, s)
}

func (a *App) createSubject(w http.ResponseWriter, r *http.Request) {
	var s subject
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&s); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := s.createSubject(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, s)
}

func (a *App) updateSubject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid subject ID")
		return
	}

	var s subject
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&s); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	s.ID = id

	if err := s.updateSubject(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, s)
}

func (a *App) deleteSubject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	s := subject{ID: id}
	if err := s.deleteSubject(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) getSubjects(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("Start"))

	if count > 10 || count < 1 {
		count = 30
	}
	if start < 0 {
		start = 0
	}
	subjects, err := getSubjects(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, subjects)
}

// user auth
func (a *App) authUser(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	var u user
	//	id, err := strconv.Atoi(vars["id"])
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// email := vars["email"]
	// password := vars["password"]

	// u := user{Email: email, Password: password}
	if err := u.authUser(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Please check your login information")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func jsonResponse(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprint(w, "message")
}
func (a *App) initializeRoutes() {

	a.Router.Handle("/", http.FileServer(http.Dir("./public")))

	// User endpoints
	a.Router.HandleFunc("/users", a.getUsers).Methods("GET")
	a.Router.HandleFunc("/user", a.createUser).Methods("POST", "OPTIONS")
	a.Router.HandleFunc("/user/{id:[0-9]+}", a.getUser).Methods("GET")
	a.Router.HandleFunc("/user/{id:[0-9]+}", a.updateUser).Methods("PUT")
	a.Router.HandleFunc("/user/{id:[0-9]+}", a.deleteUser).Methods("DELETE")

	// Auth endpoints
	a.Router.HandleFunc("/auth/user", a.authUser).Methods("POST")

	// Subject endpoints
	a.Router.HandleFunc("/subject", a.createSubject).Methods("POST", "OPTIONS")
	a.Router.HandleFunc("/subjects", a.getSubjects).Methods("GET")
	a.Router.HandleFunc("/subject/{id:[0-9]+}", a.getSubject).Methods("GET")
	a.Router.HandleFunc("/subject/{id:[0-9]+}", a.updateSubject).Methods("PUT")
	a.Router.HandleFunc("/subject/{id:[0-9]+}", a.deleteSubject).Methods("DELETE")
}
