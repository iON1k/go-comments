package models

// Полная модель комментария
type Comment struct {
	ID      int    `json:"id"`       // идентификатор комментария
	News    int    `json:"-"`        // идентификатор новости
	Parent  int    `json:"-"`        // идентификатор родительского комментария
	Content string `json:"content"`  // содержание комментария
	PubTime int64  `json:"pub_time"` // время комментария
}

type NewsComments struct {
	Comments    []Comment         `json:"сomments"`    // основные комментарии
	Subcomments map[int][]Comment `json:"subcomments"` // подкомментарии
}

type AddCommentRequest struct {
	News    *int   `json:"news_id"` // идентификатор новости
	Parent  *int   `json:"parent"`  // идентификатор родительского комментария
	Content string `json:"content"` // содержание комментария
}
