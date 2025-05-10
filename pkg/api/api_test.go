package api

import (
	"comments/pkg/models"
	"comments/pkg/storage/memdb"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type TestContext struct {
	api *API
	db  *memdb.Store
}

func setup(_ *testing.T) TestContext {
	db := memdb.New()
	db.AddComment(models.Comment{ID: 1, Content: "Comment1"})
	db.AddComment(models.Comment{ID: 2, Content: "Comment2"})
	db.AddSubcomment(models.Comment{ID: 3, Parent: 1, Content: "Subcomment1"})
	db.AddSubcomment(models.Comment{ID: 4, Parent: 2, Content: "Subcomment2"})
	api := New(db)

	return TestContext{api, db}
}

func TestAddComment(t *testing.T) {
	ctx := setup(t)

	req_json := `
	{"news_id": 1, "content":"Test"}
	`
	req := httptest.NewRequest(http.MethodPost, "/comments", strings.NewReader(req_json))
	resp := httptest.NewRecorder()
	ctx.api.router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatal("Wrong status code")
	}

	news_c := ctx.db.RawNewsComments()
	if len(news_c.Comments) != 3 || news_c.Comments[2].News != 1 || news_c.Comments[2].Content != "Test" {
		t.Fatal("Wrong data in DB")
	}
}

func TestAddSubcomment(t *testing.T) {
	ctx := setup(t)

	req_json := `
	{"parent": 3, "content":"Test"}
	`
	req := httptest.NewRequest(http.MethodPost, "/comments", strings.NewReader(req_json))
	resp := httptest.NewRecorder()
	ctx.api.router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatal("Wrong status code")
	}

	news_c := ctx.db.RawNewsComments()
	if len(news_c.Subcomments) != 3 || news_c.Subcomments[3][0].Content != "Test" {
		t.Fatal("Wrong data in DB")
	}
}

func TestGetComments(t *testing.T) {
	ctx := setup(t)

	req := httptest.NewRequest(http.MethodGet, "/comments?news_id=1", nil)
	resp := httptest.NewRecorder()
	ctx.api.router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatal("Wrong status code")
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var news_c models.NewsComments
	err = json.Unmarshal(b, &news_c)
	if err != nil {
		t.Fatal(err)
	}

	if len(news_c.Comments) != 2 || len(news_c.Subcomments) != 2 ||
		news_c.Comments[0].Content != "Comment1" || news_c.Comments[1].Content != "Comment2" ||
		len(news_c.Subcomments[1]) != 1 || news_c.Subcomments[1][0].Content != "Subcomment1" ||
		len(news_c.Subcomments[2]) != 1 || news_c.Subcomments[2][0].Content != "Subcomment2" {
		t.Fatalf("Wrong data in response")
	}
}
