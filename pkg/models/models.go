package models

// Модель комментария
type Comment struct {
	ID      int    `json:"id"`       // идентификатор комментария
	News    int    `json:"-"`        // идентификатор новости
	Parent  int    `json:"-"`        // идентификатор родительского комментария
	Content string `json:"content"`  // содержание комментария
	PubTime int64  `json:"pub_time"` // время комментария
}

// Модель коллекции комментариев к новости
type NewsComments struct {
	Comments    []Comment         `json:"сomments"`    // основные комментарии
	Subcomments map[int][]Comment `json:"subcomments"` // подкомментарии
}

// Модель запроса на создание комментария
type AddCommentRequest struct {
	News    int    `json:"news_id"` // идентификатор новости
	Parent  int    `json:"parent"`  // идентификатор родительского комментария
	Content string `json:"content"` // содержание комментария
}
