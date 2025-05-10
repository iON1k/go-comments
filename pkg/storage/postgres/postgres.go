package postgres

import (
	"comments/pkg/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Хранилище данных.
type Store struct {
	db *pgxpool.Pool
}

// Конструктор хранилища с URL для коннекта к БД
func New(conn string) (*Store, error) {
	db, err := pgxpool.New(context.Background(), conn)

	if err != nil {
		return nil, err
	}

	return NewFromPGX(db), nil
}

// Конструктор хранилища с готовым коннектом к БД
func NewFromPGX(db *pgxpool.Pool) *Store {
	return &Store{db}
}

func (s *Store) Close() {
	s.db.Close()
}

// Получение комментариев к новости
func (s *Store) NewsComments(newsId int) (models.NewsComments, error) {
	rows, err := s.db.Query(
		context.Background(),
		`
		WITH RECURSIVE news_comments AS (
			SELECT id, news, parent, content, pub_time
			FROM comments
			WHERE news = $1
			
			UNION ALL
			
			SELECT c.id, c.news, c.parent, c.content, c.pub_time
			FROM comments c
			INNER JOIN news_comments nc ON c.parent = nc.id
		)
		SELECT id, parent, content, pub_time 
		FROM news_comments
		ORDER BY pub_time ASC;
		`,
		newsId,
	)

	if err != nil {
		return models.NewsComments{}, err
	}

	all_c := []models.Comment{}
	for rows.Next() {
		var c models.Comment
		err := rows.Scan(
			&c.ID,
			&c.Parent,
			&c.Content,
			&c.PubTime,
		)

		if err != nil {
			return models.NewsComments{}, err
		}
		all_c = append(all_c, c)
	}

	var comments []models.Comment
	subcomments := map[int][]models.Comment{}
	for _, c := range all_c {
		if c.Parent == 0 {
			comments = append(comments, c)
		} else {
			sub_c := subcomments[c.Parent]
			sub_c = append(sub_c, c)
			subcomments[c.Parent] = sub_c
		}
	}

	return models.NewsComments{Comments: comments, Subcomments: subcomments}, rows.Err()
}

// Добавление нового комментария к новости
func (s *Store) AddComment(comment models.Comment) error {
	ctx := context.Background()
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		context.Background(),
		`
		INSERT INTO comments (news, content) 
		VALUES ($1, $2);
		`,
		comment.News,
		comment.Content,
	)

	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

// Добавление нового подкомментария
func (s *Store) AddSubcomment(subcomment models.Comment) error {
	ctx := context.Background()
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		context.Background(),
		`
		INSERT INTO comments (parent, content) 
		VALUES ($1, $2);
		`,
		subcomment.Parent,
		subcomment.Content,
	)

	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
