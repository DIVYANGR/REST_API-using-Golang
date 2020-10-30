package main

import (
	"database/sql"

	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type Employee struct {
	Id         int    `json:id`
	Name       string `json:name`
	City       string `json:city`
	Department string `json:department`
	Salary     int    `json:salary`
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "divyangk1998"
	dbName := "goblog"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM Employee ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	res := []Employee{}
	for selDB.Next() {
		var id, salary int
		var name, city, department string
		err = selDB.Scan(&id, &name, &city, &department, &salary)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.City = city
		emp.Department = department
		emp.Salary = salary
		res = append(res, emp)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Employee WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	for selDB.Next() {
		var id, salary int
		var name, city, department string
		err = selDB.Scan(&id, &name, &city, &department, &salary)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.City = city
		emp.Department = department
		emp.Salary = salary
	}
	tmpl.ExecuteTemplate(w, "Show", emp)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Employee WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	for selDB.Next() {
		var id, salary int
		var name, city, department string
		err = selDB.Scan(&id, &name, &city, &department, &salary)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.City = city
		emp.Department = department
		emp.Salary = salary
	}
	tmpl.ExecuteTemplate(w, "Edit", emp)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {

	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		city := r.FormValue("city")
		department := r.FormValue("department")
		salary := r.FormValue("salary")
		insForm, err := db.Prepare("INSERT INTO Employee(name, city, department, salary) VALUES(?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, city, department, salary)
		log.Println("INSERT: Name: " + name + " | City: " + city + " | Department:" + department + " | salary:" + salary)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Insert1(name1 string, city1 string, department1 string, salary1 int) {

	db := dbConn()
	name := name1
	city := city1
	department := department1
	salary := salary1
	insForm, err := db.Prepare("INSERT INTO Employee(name, city, department, salary) VALUES(?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	insForm.Exec(name, city, department, salary)
	//log.Println("INSERT: Name: " + name + " | City: " + city + " | Department:" + department + " | salary:" + salary)

	defer db.Close()

}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		city := r.FormValue("city")
		department := r.FormValue("department")
		salary := r.FormValue("salary")
		id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE Employee SET name=?, city=?, department=?, salary=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, city, department, salary, id)
		log.Println("UPDATE: Name: " + name + " | City: " + city + " | Department:" + department + " | salary:" + salary)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	emp := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM Employee WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {
	log.Println("Server started on: http://localhost:8080")

	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	//http.HandleFunc("/insert1", Insert1)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":8080", nil)
}
