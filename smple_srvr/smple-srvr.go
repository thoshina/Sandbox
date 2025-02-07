package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"time"

	//	"github.com/jinzhu/gorm"
	//	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Model struct {
	ID        uint `gotm:"primary_key" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index" json:"-"`
}
type Article struct {
	Model
	Title   string `json:"Title"`
	Desc    string `json:"Description"`
	Content string `json:"Content"`
}

type Articles []Article

var articles Articles

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!\n")
	fmt.Println("Endpoint Hit: homePage")
}

func returnArticles(w http.ResponseWriter, r *http.Request) {
	articles := Articles{}
	for i := 0; i < 10; i++ {
		title := "Hello_%d"
		articles = append(articles, Article{Title: fmt.Sprintf(title, i), Desc: "Article Description", Content: "Article Content"})
	}
	fmt.Println("Endpoint Hit: returnArticles")
	json.NewEncoder(w).Encode(articles)
}

func GetDBConn() *gorm.DB {
	db, err := gorm.Open(GetDBConfigForSQLite(&gorm.Config{}))

	if err != nil {
		panic(err)
	}

	//	db.LogMode(true)
	return db
}

func GetDBConfigForMySQL() (string, string) {
	DBMS := "mysql"
	USER := "root"
	PASS := ""
	PROTOCOL := ""
	DBNAME := "gorm-example"
	QPTION := "charset=utf8&parseTime=True&loc=Local"

	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?" + QPTION

	return DBMS, CONNECT
}

func GetDBConfigForSQLite(c *gorm.Config) (gorm.Dialector, *gorm.Config) {
	return sqlite.Open("file::memory:?cache=shared"), c
}

func fetchArticles(w http.ResponseWriter, r *http.Request) {
	db := GetDBConn()

	db.Find(&articles)
	fmt.Println(articles)
	profJson, _ := json.Marshal(articles)
	fmt.Fprintln(w, string(profJson))
	fmt.Println("Endpoint Hit: fetchArticles")
}

func writeArticles(w http.ResponseWriter, r *http.Request) {
	db := GetDBConn()

	articles := Articles{}
	for i := 0; i < 10; i++ {
		title := "Hello_%d"
		articles = append(articles, Article{Title: fmt.Sprintf(title, i), Desc: "Article Description", Content: "Article Content"})
	}

	db.Save(articles)
	fmt.Println(articles)
	profJson, _ := json.Marshal(articles)
	fmt.Fprintln(w, string(profJson))
	fmt.Println("Endpoint Hit: writeArticles")
}

func postArticles(w http.ResponseWriter, r *http.Request) {
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)

	fmt.Println(body)
	db := GetDBConn()

	var article Article

	if err := json.Unmarshal(body, &article); err != nil {
		fmt.Println(err)
		return
	}
	articles := Articles{}
	articles = append(articles, article)

	db.Save(articles)
	fmt.Println(articles)
	profJson, _ := json.Marshal(articles)
	fmt.Fprintln(w, string(profJson))
	fmt.Println("Endpoint Hit: postArticles")
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/articles", returnArticles)
	http.HandleFunc("/fetch", fetchArticles)
	http.HandleFunc("/write", writeArticles)
	http.HandleFunc("/postart", postArticles)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	db := GetDBConn()
	db.AutoMigrate(&Article{}) //SQLiteでは不要？

	handleRequests()
}
