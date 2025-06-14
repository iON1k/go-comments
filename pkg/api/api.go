package api

import (
	"comments/pkg/models"
	"comments/pkg/storage"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Программный интерфейс сервера
type API struct {
	store  storage.Store
	router *mux.Router
}

// Конструктор объекта API
func New(store storage.Store) *API {
	api := API{
		store: store,
	}
	api.router = mux.NewRouter()
	api.endpoints()
	return &api
}

// Маршрутизатор запросов.
func (api *API) Router() *mux.Router {
	return api.router
}

func (api *API) endpoints() {
	api.router.Methods(http.MethodGet).Path("/comments").HandlerFunc(api.commentsList)
	api.router.Methods(http.MethodPost).Path("/comments").HandlerFunc(api.addComment)
}

func (api *API) commentsList(w http.ResponseWriter, r *http.Request) {
	newsIdStr := r.URL.Query().Get("news_id")
	newsId, err := strconv.Atoi(newsIdStr)
	if err != nil {
		http.Error(w, "News id expected", http.StatusBadRequest)
		return
	}

	news_c, err := api.store.NewsComments(newsId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(news_c)
}

func (api *API) addComment(w http.ResponseWriter, r *http.Request) {
	var req models.AddCommentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Body decoding error", http.StatusBadRequest)
		return
	}

	if req.News == 0 && req.Parent == 0 {
		http.Error(w, "Either news id or parent expected", http.StatusBadRequest)
		return
	}

	if req.Parent != 0 {
		api.store.AddSubcomment(models.Comment{Parent: req.Parent, Content: req.Content})
	} else {
		api.store.AddComment(models.Comment{News: req.News, Content: req.Content})
	}
}
