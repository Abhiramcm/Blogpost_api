package main

import (
//"strconv"
"fmt"
"database/sql"
_ "github.com/mattn/go-sqlite3"
"net/http"
"io/ioutil"
"encoding/json"
"log"
"github.com/gorilla/mux"
)

type Blog struct {
    Id      int `json:Id`    
    Title   string `json:title` 
    Body    string `json:body` 
    Author string `json:author_name` 
    Email string `json:email` 
}

var database *sql.DB
var err error

func homePage(w http.ResponseWriter, r *http.Request){
	fmt.Println("Endpoint Hit: homePage")
	fmt.Fprintf(w, "Welcome to the HomePage!")
	}

func createBlog(w http.ResponseWriter, r *http.Request){
	reqBody,_ :=ioutil.ReadAll(r.Body)
	var blog Blog	
	json.Unmarshal([]byte(reqBody), &blog)
	database.Exec("INSERT INTO blog(title, body, author_name, email) VALUES (?,?,?,?)",
					blog.Title, blog.Body, blog.Author, blog.Email)
	fmt.Fprintf(w, "New post was created")
}

func allBlogs(w http.ResponseWriter, r *http.Request){
	fmt.Println("Endpoint Hit: AllBlogs")
	rows,_ := database.Query("select Id,title, body, author_name, email from blog where length(title)>1")
	defer rows.Close()
	for rows.Next() {
		var blog Blog
     		rows.Scan(&blog.Id, &blog.Title, &blog.Body, &blog.Author, &blog.Email)    
		fmt.Fprintf(w,"%s\n%s\n%s\n%s\n\n",blog.Title, blog.Body, blog.Author, blog.Email)  		
    	}
}

func getBlog(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	rows,_ := database.Query("select Id,title, body, author_name, email from blog where Id=?",params["id"])
	defer rows.Close()
	for rows.Next() {
		var blog Blog
     		rows.Scan(&blog.Id, &blog.Title, &blog.Body, &blog.Author, &blog.Email)    
		fmt.Fprintf(w,"%s\n%s\n%s\n%s\n\n",blog.Title, blog.Body, blog.Author, blog.Email)  		
    	}
	
}

func deletePost(w http.ResponseWriter, r *http.Request) {
  	params := mux.Vars(r)
  	statement,_ := database.Prepare("DELETE FROM blog WHERE Id = ?")
	statement.Exec(params["id"])
  	fmt.Fprintf(w, "Post with ID = %s was deleted", params["id"])
}

func BasicAuthMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		fmt.Println("username: ", user)
		fmt.Println("password: ", pass)
		if !ok || !checkUsernameAndPassword(user, pass) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password for this site"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}
		handler(w, r)
	}
}

func checkUsernameAndPassword(username, password string) bool {
	return username == "abc" && password == "123"
}

func handleRequest(){
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/",BasicAuthMiddleware(http.HandlerFunc(homePage)))
	myRouter.HandleFunc("/all",BasicAuthMiddleware(http.HandlerFunc(allBlogs)))
	myRouter.HandleFunc("/new",BasicAuthMiddleware(http.HandlerFunc(createBlog))).Methods("POST")
	myRouter.HandleFunc("/new/{id}",BasicAuthMiddleware(http.HandlerFunc(deletePost))).Methods("DELETE")
	myRouter.HandleFunc("/new/{id}",BasicAuthMiddleware(http.HandlerFunc(getBlog)))
	log.Fatal(http.ListenAndServe(":8000", myRouter))
}

func main(){
database,_ = sql.Open("sqlite3", "./blog.db")
statement,_ := database.Prepare("CREATE TABLE IF NOT EXISTS blog(Id INTEGER PRIMARY KEY,title TEXT, body TEXT, author_name TEXT, email TEXT)")	
statement.Exec()	
handleRequest()	
}
