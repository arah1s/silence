package main

import (
	"fmt"
	"html/template"
	"net/http"
	"database/sql"
	"silence/db"
	"strings"
)



const (
	resourceName = "Silence"
)

var (
	post [1]db.Post
	dbConnect *sql.DB
)

func main() {

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	http.HandleFunc("/", postHandler)
	http.HandleFunc("/write", writeHandler)
	http.HandleFunc("/contacts", contacts)
	http.HandleFunc("/SavePost", savePostHandler)

	dbConnect = db.Connect()

	fmt.Println("Listening on port:3000")
	http.ListenAndServe(":3000", nil)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	GetRandomPost := db.GetRandomPostFromDB(dbConnect)
	GetRandomPost.Content = NormaliseString(GetRandomPost.Content)
	post = newPost(GetRandomPost.Content)

	t.ExecuteTemplate(w, "index", post)
}

func contacts(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/contacts.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "write", nil)
}

func writeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/write.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "write", nil)
}

func savePostHandler(w http.ResponseWriter, r *http.Request) {
	content := r.FormValue("content")
	post = newPost(content)
	//save post in DB

	http.Redirect(w, r, "/", 302)
}

func newPost(content string) [1]db.Post {
	var post [1]db.Post
	post[0].Content = content
	post[0].Resource = resourceName

	return post
}

func NormaliseString(str string)(string){
	str = strings.Replace(str, `&#34;`, `"`, -1)
	str = strings.Replace(str, `&#39;`, `'`, -1)
	return str
}