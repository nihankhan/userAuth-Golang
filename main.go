package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"golang.org/x/crypto/bcrypt"
)

const (
	username = "root"
	password = "nihankhan"
	hostname = "127.0.0.1:3306"
	dbName   = "Nihan"
)

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

var (
	db    *sql.DB
	tmplt *template.Template
)

func main() {
	var err error
	tmplt, _ = template.ParseGlob("/home/nihan/Documents/GoSub/WEB/templates/*.html")

	db, err = sql.Open("mysql", dsn(""))

	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/signup", signUp)
	mux.HandleFunc("/login", logIn)

	mux.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates"))))

	fmt.Println("Server listening on 127.0.0.1:8000")

	http.ListenAndServe(":8000", mux)
}

func signUp(resp http.ResponseWriter, req *http.Request) {
	tmplt.ExecuteTemplate(resp, "signup##.html", nil)

	req.ParseForm()

	var err error

	name := req.FormValue("name")
	username := req.FormValue("username")
	password := req.FormValue("password")

	//check if username already exist
	stmt := "SELECT ID FROM Nihan.user WHERE username = ?"
	row := db.QueryRow(stmt, username)
	var uID string
	er := row.Scan(&uID)
	if er != sql.ErrNoRows {
		fmt.Println("username already exists, err:", er)

		return
	}

	var Password []byte

	Password, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		tmplt.Execute(resp, "there was a problem registering account!!")

		return
	}

	var insert *sql.Stmt

	insert, err = db.Prepare("INSERT INTO Nihan.user (name, username, password) VALUES (?, ?, ?);")

	if err != nil {
		tmplt.Execute(resp, "there was a problem registering account!!")

		return
	}

	fmt.Println("Name: ", name)
	fmt.Println("Username: ", username)
	fmt.Println("Password: ", string(Password))

	defer insert.Close()

	var result sql.Result

	result, _ = insert.Exec(username, Password)
	rowsAff, _ := result.RowsAffected()
	lastIns, _ := result.LastInsertId()
	fmt.Println("rowsAffect:", rowsAff)
	fmt.Println("lastInsert:", lastIns)

}

func home(resp http.ResponseWriter, req *http.Request) {
	tmplt.ExecuteTemplate(resp, "index.html", nil)
}

func logIn(resp http.ResponseWriter, req *http.Request) {
	tmplt.ExecuteTemplate(resp, "login.html", nil)

	fmt.Println("*****loginAuthHandler running*****")
	req.ParseForm()
	username := req.FormValue("username")
	password := req.FormValue("password")
	fmt.Println("username:", username) //"password:", password)
	// retrieve password from db to compare (hash) with user supplied password's hash
	var hash string
	stmt := "SELECT password FROM Nihan.user WHERE username = ?"
	row := db.QueryRow(stmt, username)
	er := row.Scan(&hash)
	fmt.Println("hash from db:", hash)
	if er != nil {
		fmt.Println("error selecting Hash in db by Username")
		//return
	}
	// func CompareHashAndPassword(hashedPassword, password []byte) error
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	// returns nill on succcess
	if err == nil {
		fmt.Fprintln(resp, "<h3>You have successfully logged in :)</h3>")
		fmt.Fprintf(resp, "<h3>\nHello %s, Welcome to your Profile.</h3>", username)
		fmt.Println("You have successfully logged in :)")

	} else {
		fmt.Fprintln(resp, "<h5>Incorrect Password. Please check your username and password!!</h5>")
	}
}