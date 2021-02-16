package main_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"example.com/social-gin/post"
	"example.com/social-gin/user"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var userHandler *user.Handler
var postHandler *post.Handler
var accessToken = "3b247d4b-071e-482a-ad4a-57a70a86bb0f"

func TestMain(m *testing.M) {
	dsn := "sqlserver://sa:1234@127.0.0.1:1433?database=social"
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&user.User{})

	client := redis.NewClient(&redis.Options{
		Addr:     "128.199.201.95:6379",
		Password: "GoLang789CodeD", // no password set
		DB:       0,                // use default DB
	})
	if _, err := client.Ping().Result(); err != nil {
		log.Fatal(err)
	}

	userHandler = &user.Handler{
		DB:          db,
		RedisClient: client,
	}
	postHandler = &post.Handler{
		DB: db,
	}

	os.Exit(m.Run())
}

func setupRouter() *gin.Engine {
	r := gin.New()

	g := r.Group("", userHandler.Authorize)

	r.POST("/login", userHandler.LogIn)
	r.GET("/users", userHandler.ListUser)
	r.GET("/users/:uid", userHandler.GetUser)
	r.POST("/users", userHandler.AddUser)
	g.PUT("/users/:uid", userHandler.UpdateUser)
	g.DELETE("/users/:uid", userHandler.DeleteUser)
	r.GET("/users/:uid/posts", postHandler.ListPost)
	r.GET("/users/:uid/posts/:pid", postHandler.GetPost)
	g.POST("/users/:uid/posts", postHandler.AddPost)
	g.PUT("/users/:uid/posts/:pid", postHandler.UpdatePost)
	g.DELETE("/users/:uid/posts/:pid", postHandler.DeletePost)

	return r
}

type Token struct {
	Token string `json:"token"`
}

func TestLogin(t *testing.T) {
	//Gin instance
	r := setupRouter()

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("u=blink&p=password"))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Error("status code is not ok", rec.Code)
		return
	}

	returnToken := Token{}
	if err := json.Unmarshal(rec.Body.Bytes(), &returnToken); err != nil {
		t.Error("can't unmarshal response", string(rec.Body.Bytes()))
		return
	}

	if len(returnToken.Token) <= 0 {
		t.Error("invalid return field")
	}

}

func TestListUsers(t *testing.T) {
	//Gin instance
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/users", strings.NewReader(""))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Error("status code is not ok", rec.Code)
		return
	}

	returnUser := []user.User{}
	if err := json.Unmarshal(rec.Body.Bytes(), &returnUser); err != nil {
		t.Error("can't unmarshal response", string(rec.Body.Bytes()))
		return
	}
}

func TestGetUser(t *testing.T) {
	//Gin instance
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/users/14", strings.NewReader(""))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Error("status code is not ok", rec.Code)
		return
	}

	returnUser := user.User{}
	if err := json.Unmarshal(rec.Body.Bytes(), &returnUser); err != nil {
		t.Error("can't unmarshal response", string(rec.Body.Bytes()))
		return
	}

	want := "blink"
	get := returnUser.Username
	if get != want {
		t.Error("given", 3, "want username", want, "but get", get)
	}

	want = "slil puangpoom"
	get = returnUser.Name
	if get != want {
		t.Error("given", 3, "want name", want, "but get", get)
	}

	want = "blink@email.com"
	get = returnUser.Email
	if get != want {
		t.Error("given", 3, "want email", want, "but get", get)
	}

}

func TestAddUsers(t *testing.T) {
	//Gin instance
	r := setupRouter()

	givenBytes, _ := json.Marshal(map[string]interface{}{
		"Username": "blink",
		"Password": "password",
		"Name":     "slil puangpoom",
		"Email":    "blink@email.com",
	})
	given := string(givenBytes)

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(given))
	req.Header.Add("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Error("status code is not ok", rec.Code)
		return
	}

	returnUser := user.User{}
	if err := json.Unmarshal(rec.Body.Bytes(), &returnUser); err != nil {
		t.Error("can't unmarshal response", string(rec.Body.Bytes()))
		return
	}

	want := "blink"
	get := returnUser.Username
	if get != want {
		t.Error("given", given, "want username", want, "but get", get)
	}

	want = "slil puangpoom"
	get = returnUser.Name
	if get != want {
		t.Error("given", given, "want name", want, "but get", get)
	}

	want = "blink@email.com"
	get = returnUser.Email
	if get != want {
		t.Error("given", given, "want email", want, "but get", get)
	}
}

func TestUpdateUsers(t *testing.T) {
	//Gin instance
	r := setupRouter()

	givenBytes, _ := json.Marshal(map[string]interface{}{
		"Name": "blink update name",
	})
	given := string(givenBytes)

	req := httptest.NewRequest(http.MethodPut, "/users/14", strings.NewReader(given))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Error("status code is not ok", rec.Code)
		return
	}

	returnUser := user.User{}
	if err := json.Unmarshal(rec.Body.Bytes(), &returnUser); err != nil {
		t.Error("can't unmarshal response", string(rec.Body.Bytes()))
		return
	}

	want := "blink"
	get := returnUser.Username
	if get != want {
		t.Error("given", given, "want username", want, "but get", get)
	}

	want = "blink update name"
	get = returnUser.Name
	if get != want {
		t.Error("given", given, "want name", want, "but get", get)
	}

	want = "blink@email.com"
	get = returnUser.Email
	if get != want {
		t.Error("given", given, "want email", want, "but get", get)
	}
}

func TestDeleteUsers(t *testing.T) {
	//Gin instance
	r := setupRouter()

	req := httptest.NewRequest(http.MethodDelete, "/users/14", strings.NewReader(""))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Error("status code is not ok", rec.Code)
		return
	}

	returnUser := user.User{}
	if err := json.Unmarshal(rec.Body.Bytes(), &returnUser); err != nil {
		t.Error("can't unmarshal response", string(rec.Body.Bytes()))
		return
	}

	want := "blink"
	get := returnUser.Username
	if get != want {
		t.Error("given", 5, "want username", want, "but get", get)
	}

}

func TestListPosts(t *testing.T) {
	//Gin instance
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/users/14/posts", strings.NewReader(""))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Error("status code is not ok", rec.Code)
		return
	}

	returnPost := []post.Post{}
	if err := json.Unmarshal(rec.Body.Bytes(), &returnPost); err != nil {
		t.Error("can't unmarshal response", string(rec.Body.Bytes()))
		return
	}
}

func TestGetPosts(t *testing.T) {
	//Gin instance
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/users/14/posts/3", strings.NewReader(""))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Error("status code is not ok", rec.Code)
		return
	}

	returnPost := post.Post{}
	if err := json.Unmarshal(rec.Body.Bytes(), &returnPost); err != nil {
		t.Error("can't unmarshal response", string(rec.Body.Bytes()))
		return
	}

	want := "TESTING"
	get := returnPost.Content
	if get != want {
		t.Error("given", 2, "want content", want, "but get", get)
	}

}

func TestAddPost(t *testing.T) {
	//Gin instance
	r := setupRouter()

	givenBytes, _ := json.Marshal(map[string]interface{}{
		"Content": "TESTING",
	})
	given := string(givenBytes)

	req := httptest.NewRequest(http.MethodPost, "/users/14/posts", strings.NewReader(given))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Error("status code is not ok", rec.Code)
		return
	}

	returnPost := post.Post{}
	if err := json.Unmarshal(rec.Body.Bytes(), &returnPost); err != nil {
		t.Error("can't unmarshal response", string(rec.Body.Bytes()))
		return
	}

	want := "TESTING"
	get := returnPost.Content
	if get != want {
		t.Error("given", given, "want content", want, "but get", get)
	}
}

func TestUpdatePost(t *testing.T) {
	//Gin instance
	r := setupRouter()

	givenBytes, _ := json.Marshal(map[string]interface{}{
		"Content": "TESTING UPDATE",
		"Likes":   1,
	})
	given := string(givenBytes)

	req := httptest.NewRequest(http.MethodPut, "/users/14/posts/3", strings.NewReader(given))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Error("status code is not ok", rec.Code)
		return
	}

	returnPost := post.Post{}
	if err := json.Unmarshal(rec.Body.Bytes(), &returnPost); err != nil {
		t.Error("can't unmarshal response", string(rec.Body.Bytes()))
		return
	}

	want := "TESTING UPDATE"
	get := returnPost.Content
	if get != want {
		t.Error("given", given, "want content", want, "but get", get)
	}

	wantInt := 1
	getInt := returnPost.Likes
	if getInt != wantInt {
		t.Error("given", given, "want like", wantInt, "but get", getInt)
	}

}

func TestDeletePost(t *testing.T) {
	//Gin instance
	r := setupRouter()

	req := httptest.NewRequest(http.MethodDelete, "/users/14/posts/3", strings.NewReader(""))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Error("status code is not ok", rec.Code)
		return
	}

	returnUser := post.Post{}
	if err := json.Unmarshal(rec.Body.Bytes(), &returnUser); err != nil {
		t.Error("can't unmarshal response", string(rec.Body.Bytes()))
		return
	}

}
