package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

type Todo struct {
	Id       int    `json:"Id,sting"`
	Title    string `json:"Title"`
	Category string `json:"Category"`
	State    string `json:"State"`
}

type Todos []Todo

type Server struct {
	db *sql.DB
}

func main() {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang_todo_dev")
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxIdleConns(100)
	defer db.Close()

	server := &Server{db: db}
	router := httprouter.New()

	router.GET("/todos/", server.todoIndex)
	router.POST("/todos/", server.todoCreate)
	router.GET("/todos/[0-9]", server.todoShow)
	router.PUT("/todos/[0-9]", server.todoUpdate)
	router.DELETE("/todos/[0-9]", server.todoDelete)

	//reHandler.HandleFunc(".*.[js|css|png|eof|svg|ttf|woff]", "GET", server.assets)
	//reHandler.HandleFunc("/", "GET", server.homepage)

	fmt.Println("Starting server on port 3000")
	http.ListenAndServe(":3000", router)
}

//Todo CRUD

func (s *Server) todoIndex(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	var todos []*Todo

	rows, err := s.db.Query("SELECT Id, Title, Category, State FROM Todo")
	error_check(res, err)
	for rows.Next() {
		todo := &Todo{}
		rows.Scan(&todo.Id, &todo.Title, &todo.Category, &todo.State)
		todos = append(todos, todo)
	}
	rows.Close()
	jsonResponse(res, todos)
}

func (s *Server) todoCreate(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	todo := &Todo{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&todo)
	if err != nil {
		fmt.Println("ERROR decoding JSON file", err)
		return
	}
	defer req.Body.Close()

	result, err := s.db.Exec("INSERT INTO Todo(Title, Category, State) VALUES(?,?,?)", todo.Title, todo.Category, todo.State)
	if err != nil {
		fmt.Println("ERROR saving to db -", err)
	}

	Id64, err := result.LastInsertId()
	Id := int(Id64)
	todo = &Todo{Id: Id}

	s.db.QueryRow("SELECT State, Title, Category FROM Todo WHERE Id=?", todo.Id).Scan(&todo.State, &todo.Title, &todo.Category)

	jsonResponse(res, todo)

}

func (s *Server) todoShow(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	fmt.Println("Render todo json")
}

func (s *Server) todoUpdate(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	todoParams := &Todo{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&todoParams)
	if err != nil {
		fmt.Println("ERROR decoding JSON -", err)
		return
	}

	_, err = s.db.Exec("UPDATE Todo SET Title = ?, Category = ?, State = ? WHERE Id = ?", todoParams.Title, todoParams.Category, todoParams.State, todoParams.Id)

	if err != nil {
		fmt.Println("ERROR saving db - ", err)
	}

	todo := &Todo{Id: todoParams.Id}
	err = s.db.QueryRow("SELECT State,Title,Category FROM Todo WHERE Id = ?", todo.Id).Scan(&todo.State, &todo.Title, &todo.Category)
	if err != nil {
		fmt.Println("ERROR reading db -", err)
	}

	jsonResponse(res, todo)

}

func (s *Server) todoDelete(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	r, _ := regexp.Compile(`¥d+$`)
	Id := r.FindString(req.URL.Path)
	s.db.Exec("DELETE FROM Todo WHERE Id = ?", Id)
	res.WriteHeader(200)
}

func jsonResponse(res http.ResponseWriter, data interface{}) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(http.StatusOK)

	//payload, err := json.Marshal(data)
	//if error_check(res, err) {
	//	return
	//}

	if err := json.NewEncoder(res).Encode(data); err != nil {
		panic(err)
	}
	//fmt.Fprintf(res, string(payload))
}

func error_check(res http.ResponseWriter, err error) bool {
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return true
	}
	return false

}
