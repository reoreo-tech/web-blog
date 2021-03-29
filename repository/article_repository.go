package repository

import (
	"database/sql"
	"math"
	"time"

	"go-tech-blog/model"
)

// ArticleCreate
func ArticleCreate(article *model.Article) (sql.Result, error) {
	now := time.Now()

	article.Created = now
	article.Updated = now

	query := `INSERT INTO articles (title, body, created, updated)
  VALUES (:title, :body, :created, :updated);`

	// トランザクションを開始
	tx := db.MustBegin()

	res, err := tx.NamedExec(query, article)
	if err != nil {
		tx.Rollback()

		return nil, err
	}

	tx.Commit()

	return res, nil
}

// ArticleListByCursor
func ArticleListByCursor(cursor int) ([]*model.Article, error) {
	if cursor <= 0 {
		cursor = math.MaxInt32
	}

	query := `SELECT *
	FROM articles
	WHERE id < ?
	ORDER BY id desc
	LIMIT 10`

	articles := make([]*model.Article, 0, 10)

	// クエリ結果を格納する変数、クエリ文字列、パラメータを指定してクエリを実行
	if err := db.Select(&articles, query, cursor); err != nil {
		return nil, err
	}

	return articles, nil
}

// ArticleDelete
func ArticleDelete(id int) error {
	query := "DELETE FROM articles WHERE id = ?"
	tx := db.MustBegin()
	if _, err := tx.Exec(query, id); err != nil {
		tx.Rollback()

		return err
	}

	return tx.Commit()
}

// ArticleGetByID
func ArticleGetByID(id int) (*model.Article, error) {
	query := `SELECT *
	FROM articles
	WHERE id = ?;`

	var article model.Article

	if err := db.Get(&article, query, id); err != nil {
		return nil, err
	}

	return &article, nil
}

// ArticleUpdate
func ArticleUpdate(article *model.Article) (sql.Result, error) {
	now := time.Now()

	article.Updated = now

	query := `UPDATE articles
	SET title = :title,
			body = :body,
			updated = :updated
	WHERE id = :id;`

	tx := db.MustBegin()

	res, err := tx.NamedExec(query, article)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return res, nil
}

// ---------------
// Appendix start
// ---------------

// // ArticleGetWithWriterName ...
// func ArticleGetWithWriterName(id int) (*model.Article, error) {
// 	// クエリ文字列を生成します。
// 	// 取得カラムは AS 句でリネームします。
// 	// リネーム後の名称は Article 構造体の db タグで指定した名称とします。
// 	// Null の可能性のあるカラムは COALESCE 関数を使って初期値を指定すると Go でのエラーを回避できます。
// 	query := `SELECT
// 		articles.id AS id,
// 		articles.title AS title,
// 		COALESCE(writers.name, '') AS writer_name
// 	FROM articles
// 	INNER JOIN writers ON writers.id = articles.writer_id
// 	WHERE articles.id = ? AND articles.writer_id IS NOT NULL;`

// 	var article model.Article
// 	if err := db.Get(&article, query, id); err != nil {
// 		return nil, err
// 	}
// 	return &article, nil
// }

// // ArticleGetWithWriter ...
// func ArticleGetWithWriter(id int) (*model.Article, error) {
// 	// 構造体を階層化した状態でデータを取得する場合は、
// 	// AS 句でのリネームでドット繋ぎの名称にします。
// 	// Article 構造体の db タグで指定した `writer` にドットで続けて、
// 	// Writer 構造体の db タグで指定した `id` と `name` を指定します。
// 	query := `SELECT
// 		articles.id AS id,
// 		articles.title AS title,
// 		writers.id AS 'writer.id',
// 		writers.name AS 'writer.name'
// 	FROM articles
// 	INNER JOIN writers ON writers.id = articles.writer_id
// 	WHERE articles.id = ?;`

// 	var article model.Article
// 	if err := db.Get(&article, query, id); err != nil {
// 		return nil, err
// 	}
// 	return &article, nil
// }

// // ArticleListByWriterID ...
// func ArticleListByWriterID(writerID int) ([]*model.Article, error) {
// 	query := `SELECT * FROM articles WHERE writer_id = ?;`
// 	var articles []*model.Article
// 	if err := db.Select(&articles, query, writerID); err != nil {
// 		return nil, err
// 	}
// 	return articles, nil
// }

// // ArticleGetWithTags ...
// func ArticleGetWithTags(id int) (*model.Article, error) {
// 	// 記事データを取得します。
// 	article, err := ArticleGetByID(id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// タグデータを取得します。
// 	tags, err := TagListByArticleID(id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 記事の構造体にタグ情報を格納します。
// 	article.Tags = tags

// 	return article, nil
// }

// // ArticleListWithTags ...
// func ArticleListWithTags() ([]*model.Article, error) {
// 	// 記事の一覧データを取得します。
// 	q1 := `SELECT id, title FROM articles;`

// 	var articles []*model.Article
// 	if err := db.Select(&articles, q1); err != nil {
// 		return nil, err
// 	}

// 	// 取得できた記事データ一覧から記事 ID を抽出します。
// 	articleIDs := make([]int, len(articles))
// 	for i, article := range articles {
// 		articleIDs[i] = article.ID
// 	}

// 	// タグ情報を map で取得します。
// 	tagListMap, err := TagListMapByArticleIDs(articleIDs)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 記事の一覧データにタグ情報を格納します。
// 	for _, article := range articles {
// 		article.Tags = tagListMap[article.ID]
// 	}

// 	return articles, nil
// }

// ---------------
// Appendix end
// ---------------
