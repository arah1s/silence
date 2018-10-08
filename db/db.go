package db

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"math/rand"
	"strconv"
	"time"
)

const (
/*	urlVPustotu                = "http://vpustotu.ru/story/"
	urlKillMePls               = "https://killpls.me/story/"
	urlPodslushano             = "https://ideer.ru/"
	urlNefart                  = "http://nefart.ru/"
	serviceNameVPustotuInDB    = "v pustotu"
	serviceNameKillMePlsInDB   = "kill me please"
	serviceNamePodslushanoInDB = "podslushano"
	serviceNameNefartInDB      = "nefart"*/

	dbHost     = "localhost"
	dbPort     = "5432"
	dbUser     = "god"
	dbPassword = "kartoshka"
	dbName     = "secretdb"
)

type Post struct {
	Id       int
	Content  string
	Likes    int
	Dislikes int
	Comments []Comment
	Resource string
}

type Comment struct {
	Id       int
	Content  string
	Likes    int
	Dislikes int
	Date     time.Time
}

func Connect()(*sql.DB) {
	dbConnectString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", dbConnectString)
	//defer db.Close()
	if err != nil {
		log.Printf("[ERROR] Database opening error -->%v\n", err)
		panic("Database error")
	}
	return db
}

//addPostInDB uppload posts in DB
func addPostInDB(db *sql.DB, posts []string, resourceName string) {
	//get resource ID from table resources
	query := fmt.Sprintf("select id from resources where name = '%s'", resourceName)
	resourceID, _ := strconv.Atoi(getStringFromDB(db, query))

	//get count rows in table posts before upload content
	query = fmt.Sprintf("select count(*) from posts where resource_id = %d", resourceID)
	countRowsBeforeUploadPosts, _ := strconv.Atoi(getStringFromDB(db, query))

	for i := 0; i < len(posts); i++ {
		query = "select max(id) from posts"
		maxID, _ := strconv.Atoi(getStringFromDB(db, query))

		hash := GetMD5Hash(posts[i])
		query = fmt.Sprintf("select hash from post_description where hash = '%s'", hash)
		hashInDB := getStringFromDB(db, query)
		//add only new posts
		if hashInDB == "" {
			query = fmt.Sprintf(`insert into posts(id, resource_id, active, likes, dislikes) values (%d, %d, true, 0, 0);`+
				`insert into post_description (post_id, text, hash) values (%d, '%s', '%s')`,
				maxID+1, resourceID, maxID+1, posts[i], hash)
			err := addRowInDB(db, query)
			if err != nil {
				fmt.Println("[ERROR]: Error insert row in DB :", err)
			}
		}
	}

	//get count rows in table posts after upload content
	query = fmt.Sprintf("select count(*) from posts where resource_id = %d", resourceID)
	countRowsAfterUploadPosts, _ := strconv.Atoi(getStringFromDB(db, query))

	fmt.Println("[LOG]: Load ", (countRowsAfterUploadPosts - countRowsBeforeUploadPosts), " posts from '"+resourceName+"'")
}

//addRowInDB added new row in database
func addRowInDB(db *sql.DB, query string) (err error){
	rows, err := db.Query(query)
	if err != nil {
		fmt.Printf("[ERROR]: Error insert row in DB with query '%s' error:%v", query, err)
	}
	defer rows.Close()
	return err
}

//getStringFromDB extract one string from db
func getStringFromDB(db *sql.DB, query string) (result string) {
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&result)
		if err != nil {
			log.Fatal(err)
		}
	}

	return result
}

//getRandomPostFromDB
func GetRandomPostFromDB(db *sql.DB) (post Post) {
	query := "select max(id) from posts"
	maxID, _ := strconv.Atoi(getStringFromDB(db, query))
	randomID := rand.Intn(maxID)
	fmt.Println(randomID)
	query = fmt.Sprintf("select posts.id, likes, dislikes, url, text from posts "+
		"join resources on posts.resource_id = resources.id "+
		"join post_description on posts.id = post_description.post_id "+
		"where posts.active = true and resources.active = true and posts.id >= %d limit 1;", randomID)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&post.Id, &post.Likes, &post.Dislikes, &post.Resource, &post.Content)
		if err != nil {
			log.Fatal(err)
		}
	}

	return post
}

//GetMD5Hash calculate hash for post
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}