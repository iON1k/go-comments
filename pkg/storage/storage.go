package storage

import "comments/pkg/models"

// Хранилище данных
type Store interface {
	NewsComments(newsId int) (models.NewsComments, error) // Получение комментариев к новости
	AddComment(comment models.Comment) error              // Добавление нового комментария к новости
	AddSubcomment(subcomment models.Comment) error        // Добавление нового подкомментария
}
