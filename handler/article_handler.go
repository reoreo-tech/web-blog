package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-tech-blog/model"
	"go-tech-blog/repository"

	"github.com/labstack/echo/v4"
)

// ArticleIndex
func ArticleIndex(c echo.Context) error {
	if c.Request().URL.Path == "/articles" {
		c.Redirect(http.StatusPermanentRedirect, "/")
	}

	articles, err := repository.ArticleListByCursor(0)

	if err != nil {

		c.Logger().Error(err.Error())

		return c.NoContent(http.StatusInternalServerError)
	}

	var cursor int
	if len(articles) != 0 {
		cursor = articles[len(articles)-1].ID
	}

	data := map[string]interface{}{
		"Articles": articles,
		"Cursor":   cursor,
	}

	return render(c, "article/index.html", data)
}

// ArticleNew
func ArticleNew(c echo.Context) error {
	data := map[string]interface{}{
		"Message": "Article New",
		"Now":     time.Now(),
	}

	return render(c, "article/new.html", data)
}

// ArticleShow ...
func ArticleShow(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("articleID"))

	// 記事データを取得
	article, err := repository.ArticleGetByID(id)

	if err != nil {
		// エラー内容をサーバーのログに出力
		c.Logger().Error(err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	data := map[string]interface{}{
		"Article": article,
	}

	return render(c, "article/show.html", data)
}

// ArticleEdit
func ArticleEdit(c echo.Context) error {
	// パスパラメータから記事IDを取得
	id, _ := strconv.Atoi(c.Param("articleID"))

	// 編集フォームの初期値として表示するために記事データを取得
	article, err := repository.ArticleGetByID(id)

	if err != nil {
		c.Logger().Error(err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	data := map[string]interface{}{
		"Article": article,
	}
	return render(c, "article/edit.html", data)
}

// ArticleCreateOutput
type ArticleCreateOutput struct {
	Article          *model.Article
	Message          string
	ValidationErrors []string
}

// ArticleCreate
func ArticleCreate(c echo.Context) error {
	var article model.Article
	var out ArticleCreateOutput

	if err := c.Bind(&article); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, out)
	}

	// バリデーションチェックを実行
	if err := c.Validate(&article); err != nil {
		c.Logger().Error(err.Error())
		// エラー内容を検査してカスタムエラーメッセージを取得
		out.ValidationErrors = article.ValidationErrors(err)
		return c.JSON(http.StatusUnprocessableEntity, out)
	}

	res, err := repository.ArticleCreate(&article)
	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, out)
	}

	id, _ := res.LastInsertId()
	article.ID = int(id)
	out.Article = &article
	return c.JSON(http.StatusOK, out)
}

// ArticleDelete
func ArticleDelete(c echo.Context) error {
	// パスパラメータから記事 ID を取得
	id, _ := strconv.Atoi(c.Param("articleID"))
	// repository の記事削除処理を呼び出し
	if err := repository.ArticleDelete(id); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("Article %d is deleted.", id))
}

// ArticleList
func ArticleList(c echo.Context) error {
	cursor, _ := strconv.Atoi(c.QueryParam("cursor"))
	// リポジトリの処理を呼び出して記事の一覧データを取得
	articles, err := repository.ArticleListByCursor(cursor)
	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, "")
	}
	return c.JSON(http.StatusOK, articles)
}

// ArticleUpdateOutput
type ArticleUpdateOutput struct {
	Article          *model.Article
	Message          string
	ValidationErrors []string
}

// ArticleUpdate
func ArticleUpdate(c echo.Context) error {
	// リクエスト送信元のパスを取得
	ref := c.Request().Referer()
	refID := strings.Split(ref, "/")[4]
	reqID := c.Param("articleID")
	if reqID != refID {
		return c.JSON(http.StatusBadRequest, "")
	}

	var article model.Article
	var out ArticleUpdateOutput

	if err := c.Bind(&article); err != nil {
		return c.JSON(http.StatusBadRequest, out)
	}

	// 入力値のバリデーションチェック
	if err := c.Validate(&article); err != nil {
		out.ValidationErrors = article.ValidationErrors(err)
		return c.JSON(http.StatusUnprocessableEntity, out)
	}

	articleID, _ := strconv.Atoi(reqID)
	article.ID = articleID
	_, err := repository.ArticleUpdate(&article)

	if err != nil {
		out.Message = err.Error()
		return c.JSON(http.StatusInternalServerError, out)
	}

	out.Article = &article
	return c.JSON(http.StatusOK, out)
}
