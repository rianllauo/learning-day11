package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"personal-web/connection"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()

	connection.DatabaseConnect()

	//route untuk menginisialisai folder public
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/", home).Methods("GET")
	// route.HandleFunc("/formblog", formBlog).Methods("GET")
	route.HandleFunc("/blog-detail/{id}", blogDetail).Methods("GET")
	route.HandleFunc("/blog", blog).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/add-blog", addBlog).Methods("POST")
	route.HandleFunc("/add-blog", formAddBlog).Methods("GET")
	route.HandleFunc("/delete-blog/{index}", deleteBlog).Methods("GET")

	fmt.Println("Server berjalan pada port 5000")
	http.ListenAndServe("localhost:5000", route)
}

type Blog struct {
	ID        int
	Title     string
	Content   string
	Image     string
	Post_date time.Time
	Author    string
}

// var blogs = []
var blogs = []Blog{
	{
		Title:   "Samsul Rijal",
		Content: "Hallo Dumbways",
		// Post_date: "24 November 2022",
		Author: "Samsul Rijal",
	},
}

func addBlog(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	var newBlog = Blog{
		Title:   title,
		Content: content,
		// Post_date: "24 November 2022",
		Author: "Samsul Rijal",
	}

	// blogs.push(newBlog)
	blogs = append(blogs, newBlog)

	http.Redirect(w, r, "/blog", http.StatusMovedPermanently)
}

func blog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/blog.html")

	// if condition
	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	dataBlog, errQuery := connection.Conn.Query(context.Background(), "SELECT id, title, content, post_date  FROM tb_blog")
	if errQuery != nil {
		fmt.Println("Message : " + errQuery.Error())
		return
	}

	var result []Blog

	for dataBlog.Next() {
		var each = Blog{}

		err := dataBlog.Scan(&each.ID, &each.Title, &each.Content, &each.Post_date)
		if err != nil {
			fmt.Println("Message : " + err.Error())
			return
		}

		result = append(result, each)
	}

	fmt.Println(result)
	resData := map[string]interface{}{
		"Blogs": result,
	}

	tmpt.Execute(w, resData)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/index.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	tmpt.Execute(w, nil)
}

func formAddBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/add-blog.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	tmpt.Execute(w, nil)
}

func blogDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/blog-detail.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// object golang
	// data := map[string]interface{}{
	// 	"Title":   "Pasar Coding Dari Dumbways",
	// 	"Content": "REPUBLIKA.CO.ID, JAKARTA -- Ketimpangan sumber daya manusia (SDM) di sektor digital masih menjadi isu yang belum terpecahkan.",
	// 	"Id":      id,
	// }
	var BlogDetail = Blog{}

	for index, data := range blogs {
		if index == id {
			BlogDetail = Blog{
				Title:     data.Title,
				Content:   data.Content,
				Post_date: data.Post_date,
				Author:    data.Author,
			}
		}
	}

	fmt.Println(BlogDetail)

	dataDetail := map[string]interface{}{
		"Blog": BlogDetail,
	}

	tmpt.Execute(w, dataDetail)
}

func deleteBlog(w http.ResponseWriter, r *http.Request) {

	index, _ := strconv.Atoi(mux.Vars(r)["index"])
	// fmt.Println(index)

	blogs = append(blogs[:index], blogs[index+1:]...)

	http.Redirect(w, r, "/blog", http.StatusFound)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/contact.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	tmpt.Execute(w, nil)
}
