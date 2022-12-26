package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	Routers()
}

func Routers() {
	InitDB()
	defer db.Close()
	log.Println("Starting the HTTP server on port 9080")
	router := mux.NewRouter()
	router.HandleFunc("/quests", GetQuests).Methods("GET")
	router.HandleFunc("/quests", CreateQuest).Methods("POST")
	router.HandleFunc("/quests/{id}", GetQuest).Methods("GET")
	router.HandleFunc("/quests/{id}", UpdateQuest).Methods("PUT")
	router.HandleFunc("/quests/{id}", DeleteQuest).Methods("DELETE")
	router.HandleFunc("/login", Login).Methods("POST")
	router.HandleFunc("/register", Register).Methods("POST")
	http.ListenAndServe(":9080", &CORSRouterDecorator{router})
}

/***************************************************/

// Get all quest
func GetQuests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var users []Quest
	result, err := db.Query("SELECT id, questname, description, score, value, impact, deadline FROM quest")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var user Quest
		err := result.Scan(&user.ID, &user.Questname, &user.Description, &user.Score, &user.Value, &user.Impact, &user.Deadline)
		if err != nil {
			panic(err.Error())
		}
		users = append(users, user)
	}
	json.NewEncoder(w).Encode(users)
}

// Create quest
func CreateQuest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	stmt, err := db.Prepare("INSERT INTO quest(questname, description, score, value, impact, deadline) VALUES(?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	questname := keyVal["questname"]
	description := keyVal["description"]
	score := keyVal["score"]
	value := keyVal["value"]
	impact := keyVal["impact"]
	deadline := keyVal["deadline"]
	_, err = stmt.Exec(questname, description, score, value, impact, deadline)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "New quest was created")
}

// Get quest by ID
func GetQuest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)
	result, err := db.Query("SELECT id, questname, description, score, value, impact, deadline FROM quest WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var user Quest
	for result.Next() {
		err := result.Scan(&user.ID, &user.Questname, &user.Description, &user.Score, &user.Value, &user.Impact, &user.Deadline)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(user)
}

// Update quest
func UpdateQuest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)
	stmt, err := db.Prepare("UPDATE quest SET questname = ?, description = ?, score = ?, value = ?, impact = ?, deadline = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	questname := keyVal["questname"]
	description := keyVal["description"]
	score := keyVal["score"]
	value := keyVal["value"]
	impact := keyVal["impact"]
	deadline := keyVal["deadline"]
	_, err = stmt.Exec(questname, description, score, value, impact, deadline, params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Quest with ID = %s was updated",
		params["id"])
}

// Delete quest
func DeleteQuest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM quest WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Quest with ID = %s was deleted",
		params["id"])
}

// Register
func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	stmt, err := db.Prepare("INSERT INTO user(name, username, email, password) VALUES(?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	name := keyVal["name"]
	username := keyVal["username"]
	email := keyVal["email"]
	password := keyVal["password"]
	_, err = stmt.Exec(name, username, email, password)
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "New account was created")
}

// Login
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	username := keyVal["username"]
	password := keyVal["password"]
	result, err := db.Query("SELECT id, name, username, email, password FROM user WHERE username = ? and password = ?", username, password)
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var user User
	for result.Next() {
		err := result.Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.Password)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(user)
}

/*============== FOR QUEST LIST ================*/
type Quest struct {
	ID          string `json:"id"`
	Questname   string `json:"questname"`
	Description string `json:"description"`
	Score       int    `json:"score"`
	Value       int    `json:"value"`
	Impact      string `json:"impact"`
	Deadline    string `json:"deadline"`
}

/*============== FOR LOGIN REGISTER ===============*/
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Db configuration
var db *sql.DB
var err error

func InitDB() {
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/questroom")
	if err != nil {
		panic(err.Error())
	}
}

/***************************************************/

// CORSRouterDecorator applies CORS headers to a mux.Router
type CORSRouterDecorator struct {
	R *mux.Router
}

func (c *CORSRouterDecorator) ServeHTTP(rw http.ResponseWriter,
	req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods",
			"POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers",
			"Accept, Accept-Language,"+
				" Content-Type, YourOwnHeader")
	}
	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}

	c.R.ServeHTTP(rw, req)
}
