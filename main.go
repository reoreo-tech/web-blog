package main

import (
	"log"
	"os"

	"go-tech-blog/handler"
	"go-tech-blog/repository"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/go-playground/validator.v9"
)

var authUser = os.Getenv("AUTH_USER")
var authPassword = os.Getenv("AUTH_PASSWORD")
var db *sqlx.DB
var e = createMux()

func main() {
	db = connectDB()
	repository.SetDB(db)

	auth := e.Group("")
	auth.Use(basicAuth())

	e.GET("/", handler.ArticleIndex)

	e.GET("/articles", handler.ArticleIndex)                   // 一覧画面
	auth.GET("/articles/new", handler.ArticleNew)              // 新規作成画面
	e.GET("/articles/:articleID", handler.ArticleShow)         // 詳細画面
	auth.GET("/articles/:articleID/edit", handler.ArticleEdit) // 編集画面

	e.GET("/api/articles", handler.ArticleList)                    // 一覧
	auth.POST("/api/articles", handler.ArticleCreate)              // 作成
	auth.DELETE("/api/articles/:articleID", handler.ArticleDelete) // 削除
	auth.PATCH("/api/articles/:articleID", handler.ArticleUpdate)  // 更新

	// [Appendix start]
	// e.GET("/test", handler.Test)
	// [Appendix end]

	port := os.Getenv("PORT")
	if port == "" {
		e.Logger.Fatal("$PORT must be set")
	}
	e.Logger.Fatal(e.Start(":" + port))
}

func createMux() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.Gzip())
	e.Use(middleware.CSRF())

	e.Static("/css", "src/css")
	e.Static("/js", "src/js")

	e.Validator = &CustomValidator{validator: validator.New()}

	return e
}

func connectDB() *sqlx.DB {
	dsn := os.Getenv("DSN")
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		e.Logger.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		e.Logger.Fatal(err)
	}
	log.Println("db connection succeeded")
	return db
}

// CustomValidator
type CustomValidator struct {
	validator *validator.Validate
}

// Validate
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func basicAuth() echo.MiddlewareFunc {
	return middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == authUser && password == authPassword {
			return true, nil
		}
		return false, nil
	})
}
