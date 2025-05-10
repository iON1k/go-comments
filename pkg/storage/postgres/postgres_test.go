package postgres

import (
	"comments/pkg/models"
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type TestContext struct {
	s *Store
}

func setup(t *testing.T) TestContext {
	ctx := TestContext{makeStorage(t)}

	comments := []models.Comment{
		{News: 1, Content: "Comment1"},
		{News: 2, Content: "Comment2"},
		{News: 1, Content: "Comment3"},
	}

	subcomments := []models.Comment{
		{Parent: 1, Content: "Subcomment1"},
		{Parent: 1, Content: "Subcomment2"},
		{Parent: 3, Content: "Subcomment3"},
		{Parent: 3, Content: "Subcomment4"},
		{Parent: 2, Content: "Subcomment5"},
		{Parent: 2, Content: "Subcomment6"},
		{Parent: 4, Content: "Subcomment7"},
		{Parent: 5, Content: "Subcomment8"},
		{Parent: 6, Content: "Subcomment9"},
		{Parent: 7, Content: "Subcomment10"},
	}

	for _, c := range comments {
		err := ctx.s.AddComment(c)
		if err != nil {
			t.Fatal(err)
		}
	}

	for _, c := range subcomments {
		err := ctx.s.AddSubcomment(c)
		if err != nil {
			t.Fatal(err)
		}
	}

	return ctx
}

func (c TestContext) teardown() {
	if c.s != nil {
		c.s.Close()
	}
}

func TestNewsComments(t *testing.T) {
	ctx := setup(t)
	defer ctx.teardown()

	news_c, err := ctx.s.NewsComments(1)
	if err != nil {
		t.Fatal(err)
	}

	if len(news_c.Comments) != 2 || news_c.Comments[0].Content != "Comment1" || news_c.Comments[1].Content != "Comment3" {
		t.Fatal("Got wrong comments from DB")
	}

	if len(news_c.Subcomments) != 6 || len(news_c.Subcomments[3]) != 2 ||
		news_c.Subcomments[3][0].Content != "Subcomment3" || news_c.Subcomments[3][1].Content != "Subcomment4" ||
		news_c.Subcomments[7][0].Content != "Subcomment10" || news_c.Subcomments[2] != nil {
		t.Fatal("Got wrong subcomments from DB")
	}
}

func makeStorage(t *testing.T) *Store {
	err := godotenv.Load()
	if err != nil {
		t.Fatal(err)
	}

	db_conn := os.Getenv("TEST_DB")
	if db_conn == "" {
		t.Fatal("No environment for DB")
	}

	db, err := pgxpool.New(context.Background(), db_conn)
	if err != nil {
		t.Fatal(err)
	}

	bytes, err := os.ReadFile("schema.sql")
	if err != nil {
		db.Close()
		t.Fatal(err)
	}

	_, err = db.Exec(context.Background(), string(bytes))
	if err != nil {
		db.Close()
		t.Fatal(err)
	}

	return NewFromPGX(db)
}
