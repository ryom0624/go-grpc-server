package repository

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"go-grpc-server/article/pb"
)

type Repository interface {
	InsertArticle(ctx context.Context, input *pb.ArticleInput) (int64, error)
	SelectArticleByID(ctx context.Context, id int64) (*pb.Article, error)
	UpdateArticle(ctx context.Context, id int64, input *pb.ArticleInput) error
	DeleteArticle(ctx context.Context, id int64) error
	SelectAllArticles() (*sql.Rows, error)
}

type sqliteRepo struct {
	db *sql.DB
}

var _ Repository = (*sqliteRepo)(nil)

func NewsqliteRepo() (Repository, error) {
	db, err := sql.Open("sqlite3", "./article/article.sql")
	if err != nil {
		return nil, err
	}

	cmd := `CREATE TABLE IF NOT EXISTS articles(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		author STRING,
		title STRING,
		content STRING)`

	_, err = db.Exec(cmd)
	if err != nil {
		return nil, err
	}

	return sqliteRepo{db: db}, nil
}

func (s sqliteRepo) InsertArticle(ctx context.Context, input *pb.ArticleInput) (int64, error) {
	cmd := `INSERT INTO articles(author, title, content) VALUES (?,?,?)`
	result, err := s.db.Exec(cmd, input.Author, input.Title, input.Content)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s sqliteRepo) SelectArticleByID(ctx context.Context, id int64) (*pb.Article, error) {
	cmd := `SELECT id, author, title, content FROM articles WHERE id = ?`
	row := s.db.QueryRow(cmd, id)

	var a pb.Article

	err := row.Scan(&a.Id, &a.Author, &a.Title, &a.Content)
	if err != nil {
		return nil, err
	}

	return &pb.Article{
		Id:      a.Id,
		Author:  a.Author,
		Title:   a.Title,
		Content: a.Content,
	}, nil

}

func (s sqliteRepo) UpdateArticle(ctx context.Context, id int64, input *pb.ArticleInput) error {
	cmd := `UPDATE articles SET author= ?, title= ?, content= ? WHERE id = ?`
	_, err := s.db.Exec(cmd, input.Author, input.Title, input.Content, id)
	if err != nil {
		return err
	}
	return nil
}

func (s sqliteRepo) DeleteArticle(ctx context.Context, id int64) error {
	cmd := `DELETE FROM articles WHERE id = ?`
	_, err := s.db.Exec(cmd, id)
	if err != nil {
		return err
	}
	return nil
}

func (s sqliteRepo) SelectAllArticles() (*sql.Rows, error) {
	cmd := `SELECT * FROM articles`
	rows, err := s.db.Query(cmd)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
