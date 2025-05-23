package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"time"

	//	"github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
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

type TargetKey struct {
	ID uint `json:"id,string"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!\n")
	fmt.Println("Endpoint Hit: homePage")
}

func getDemoArticles(c int64) Articles {
	articles := Articles{}
	var i int64
	for i = 0; i < c; i++ {
		title := "Hello_%d"
		articles = append(articles, Article{Title: fmt.Sprintf(title, i), Desc: "Article Description", Content: "Article Content"})
	}
	return articles
}

func returnArticles(w http.ResponseWriter, r *http.Request) {
	articles := getDemoArticles(10)

	fmt.Println("Endpoint Hit: returnArticles")
	json.NewEncoder(w).Encode(articles)
}

func GetDBConn() *gorm.DB {
	//db, err := gorm.Open(GetDBConfigForSQLite(&gorm.Config{}))
	db, err := gorm.Open(GetDBConfigForMySQL(&gorm.Config{}))

	if err != nil {
		panic(err)
	}

	//	db.LogMode(true)
	return db
}

func GetDBConfigForMySQL(c *gorm.Config) (gorm.Dialector, *gorm.Config) {
	//	DBMS := "mysql"
	USER := "root"
	PASS := ""
	PROTOCOL := "tcp(0.0.0.0:3306)"
	DBNAME := "smplSvr_example"
	//Databaseは作成しておく必要あり Tableは未作成でも大丈夫
	QPTION := "charset=utf8&parseTime=True&loc=Local"

	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?" + QPTION

	return mysql.Open(CONNECT), c
}

func GetDBConfigForSQLite(c *gorm.Config) (gorm.Dialector, *gorm.Config) {
	return sqlite.Open("file::memory:?cache=shared"), c
}

func getTargetKey(r *http.Request) (TargetKey, int64) {
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)

	fmt.Println(body)
	var tgtKey TargetKey

	if len > 0 {
		if err := json.Unmarshal(body, &tgtKey); err != nil {
			fmt.Println(err)
			return TargetKey{}, len
		}
	}

	return tgtKey, len
}

func fetchArticles(w http.ResponseWriter, r *http.Request) {
	shwKey, len := getTargetKey(r)

	fmt.Printf("ID=%d\n", shwKey.ID)

	db := GetDBConn()

	if len > 0 {
		db.Find(&articles, shwKey.ID)
	} else {
		db.Find(&articles) //jsonデータ無しならすべて
	}

	fmt.Println(articles)
	profJson, _ := json.Marshal(articles)
	fmt.Fprintln(w, string(profJson))
	fmt.Println("Endpoint Hit: fetchArticles")
}

func writeArticles(w http.ResponseWriter, r *http.Request) {
	db := GetDBConn()

	articles := getDemoArticles(10)

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

func deleteArticles(w http.ResponseWriter, r *http.Request) {
	delKey, len := getTargetKey(r)
	fmt.Printf("ID=%d\n", delKey.ID)

	db := GetDBConn()

	db.Find(&articles)
	result := func(db gorm.DB, l int64, i uint) (tx *gorm.DB) {
		if l > 0 {
			return db.Delete(articles, i)
		} else {
			return db.Delete(articles) //jsonデータ無しなら全削除
		}
	}(*db, len, delKey.ID)
	fmt.Fprintf(w, "Delete %d records.\n", result.RowsAffected) //返り値もDBで、この場合 RowsAffected に削除したレコード数が入る
	fmt.Println("Endpoint Hit: deleteArticles")
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/articles", returnArticles)
	http.HandleFunc("/fetch", fetchArticles)
	http.HandleFunc("/write", writeArticles)
	http.HandleFunc("/postart", postArticles)
	http.HandleFunc("/delete", deleteArticles)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	db := GetDBConn()
	db.AutoMigrate(&Article{}) //SQLiteでは不要？

	handleRequests()
}
