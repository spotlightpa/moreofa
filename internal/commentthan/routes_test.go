package commentthan

import (
	"context"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/carlmjohnson/be"
	"github.com/carlmjohnson/be/testfile"
	"github.com/carlmjohnson/requests"
	"github.com/carlmjohnson/requests/reqtest"
	"github.com/spotlightpa/moreofa/internal/db"
)

func clearComment(c *db.Comment) {
	c.ID = 0
	c.Ip = "1.2.3.4"
	c.CreatedAt = time.Date(1969, 7, 20, 20, 17, 40, 0, time.UTC)
	c.ModifiedAt = time.Date(1969, 7, 20, 20, 17, 40, 0, time.UTC)
}

func TestPostComment(t *testing.T) {
	dir := t.TempDir()
	dbname := filepath.Join(dir, "test.db")
	ctx := context.Background()
	app := appEnv{
		dbname: dbname,
	}
	be.NilErr(t, app.configureService())
	defer app.closeService()

	h := app.router()
	srv := httptest.NewServer(h)
	defer srv.Close()

	err := requests.
		New(reqtest.Server(srv)).
		Path("/comment").
		BodyForm(nil).
		Fetch(ctx)
	be.NilErr(t, err)
	comments, err := app.svc.q.ListComments(ctx, db.ListCommentsParams{
		Limit:  1,
		Offset: 0,
	})
	for i := range comments {
		clearComment(&comments[i])
	}
	be.NilErr(t, err)
	testfile.EqualJSON(t, "testdata/comments.json", comments)
}
