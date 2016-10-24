package main

import (
	"net/http"
	"io/ioutil"
	"strings"
	"text/template"
	"log"
	"database/sql"
	"strconv"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type MyHandler struct {
}

type  CSSHandler struct{
}

type AddHandler struct {
}

type Context struct {
	Numbers []int
}



func (this *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//path := r.URL.Path[1:]
	path := "templates/index.html"
	data, err := ioutil.ReadFile(string(path))
	log.Println(path)

	if err == nil {
		var contentType string

		if strings.HasSuffix(path, ".css") {
			contentType = "text/css"
		} else if strings.HasSuffix(path, ".html") {
			contentType = "text/html"
		}  else {
			contentType = "text/plain"
		}
		w.Header().Add("Content-Type", contentType)
		tmpl, err := template.New("TestTemplate").Parse(string(data))
		if err == nil {
			db, error := sql.Open("mysql", "root:ADASARE@tcp(localhost:3306)/sphere?charset=utf8")
			checkErr(error)

			rows, error := db.Query("SELECT number FROM data")
			checkErr(error)

			var numbers []int

			for rows.Next() {
				var number int
				error = rows.Scan(&number)
				checkErr(error)
				numbers = append(numbers, number)
				//fmt.Println(number)
			}

			db.Close()

			context := Context{
				numbers,
			}

			tmpl.Execute(w, context)
		}
	} else {
		w.WriteHeader(404)
		w.Write([]byte("404 - " + http.StatusText(404)))
	}
}

func (this *AddHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//path := r.URL.Path[1:]
	r.ParseForm()
	path := "templates/index.html"
	data, err := ioutil.ReadFile(string(path))
	log.Println(path)

	if err == nil && r.Method == "POST" {
		var contentType string

		if strings.HasSuffix(path, ".css") {
			contentType = "text/css"
		} else if strings.HasSuffix(path, ".html") {
			contentType = "text/html"
		}  else {
			contentType = "text/plain"
		}

		w.Header().Add("Content-Type", contentType)
		tmpl, err := template.New("TestTemplate").Parse(string(data))
		if err == nil {

			db, error := sql.Open("mysql", "root:ADASARE@tcp(localhost:3306)/sphere?charset=utf8")
			checkErr(error)

			if dataf, errors := strconv.Atoi(r.Form["dataentry"][0]); errors == nil{
				stmt, error := db.Prepare("INSERT data SET number=?")
				checkErr(error)

				res, error := stmt.Exec(dataf)
				checkErr(error)

				id, error := res.LastInsertId()
				checkErr(error)

				fmt.Println(id)
			}

			rows, error := db.Query("SELECT number FROM data")
			checkErr(error)

			var numbers []int

			for rows.Next() {
				var number int
				error = rows.Scan(&number)
				checkErr(error)
				numbers = append(numbers, number)
				//fmt.Println(number)
			}

			db.Close()

			context := Context{
				numbers,
			}

			tmpl.Execute(w, context)
		}
	} else {
		w.WriteHeader(404)
		w.Write([]byte(""))
	}
}

func (this *CSSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	data, err := ioutil.ReadFile(string(path))
	log.Println(path)
	var contentType string = "text/css"
	w.Header().Add("Content-Type", contentType)
	if err == nil {
		w.Write(data)
	} else {
		w.WriteHeader(404)
		w.Write([]byte("404 - " + http.StatusText(404)))
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	http.Handle("/", new(MyHandler))
	http.Handle("/public/css/styles.css", new(CSSHandler))
	http.Handle("/add", new(AddHandler))
	http.ListenAndServe(":8080", nil)
}