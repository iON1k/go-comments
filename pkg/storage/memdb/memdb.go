package memdb

import (
	"comments/pkg/models"
	"sync"
)

// Хранилище данных в памяти
type Store struct {
	c   models.NewsComments
	mut sync.Mutex
}

// Конструктор хранилища.
func New() *Store {
	c := models.NewsComments{Subcomments: map[int][]models.Comment{}}
	return &Store{c, sync.Mutex{}}
}

// Получение комментариев к новости
func (s *Store) NewsComments(newsId int) (models.NewsComments, error) {
	s.mut.Lock()
	defer s.mut.Unlock()
	return s.c, nil
}

// Добавление нового комментария к новости
func (s *Store) AddComment(comment models.Comment) error {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.c.Comments = append(s.c.Comments, comment)
	return nil
}

// Добавление нового подкомментария
func (s *Store) AddSubcomment(subcomment models.Comment) error {
	s.mut.Lock()
	defer s.mut.Unlock()
	subcomments := s.c.Subcomments[subcomment.Parent]
	subcomments = append(subcomments, subcomment)
	s.c.Subcomments[subcomment.Parent] = subcomments
	return nil
}

func (s *Store) RawNewsComments() models.NewsComments {
	return s.c
}
