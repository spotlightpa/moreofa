package commentthan

import (
	"context"
	"net/http/httptest"
	"net/url"
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
	c.Ip = "127.0.0.1"
	c.CreatedAt = time.Date(1969, 7, 20, 20, 17, 40, 0, time.UTC)
	c.ModifiedAt = time.Date(1969, 7, 20, 20, 17, 40, 0, time.UTC)
}

func testService(t *testing.T) *service {
	dir := t.TempDir()
	dbname := filepath.Join(dir, "test.db")
	app := appEnv{
		dbname: dbname,
	}
	svc, err := app.configureService()
	be.NilErr(t, err)
	t.Cleanup(func() {
		svc.closeService()
	})
	return svc
}

func (svc *service) testRouter() *httptest.Server {
	h := svc.router()
	srv := httptest.NewServer(h)
	srv.Client().CheckRedirect = requests.NoFollow
	return srv
}

func TestHealthcheck(t *testing.T) {
	t.Parallel()
	svc := testService(t)
	ctx := context.Background()

	srv := svc.testRouter()
	defer srv.Close()

	var body string
	rb := requests.
		New(reqtest.Server(srv)).
		Path("/api/healthcheck").
		ToString(&body)

	be.NilErr(t, rb.Fetch(ctx))
	be.Equal(t, "OK", body)
}

func TestPostComment(t *testing.T) {
	t.Parallel()
	svc := testService(t)
	ctx := context.Background()

	srv := svc.testRouter()
	defer srv.Close()

	rb := requests.
		New(reqtest.Server(srv)).
		Path("/comment")

	be.NilErr(t, rb.Clone().
		BodyForm(url.Values{
			"bot-field": []string{},
			"host_page": []string{"host_page1"},
			"name":      []string{"name1"},
			"email":     []string{"email1"},
			"CC":        []string{"CC1"},
			"subject":   []string{"subject1"},
			"anonymous": []string{},
			"comment":   []string{"comment1"},
		}).
		CheckStatus(303).
		Fetch(ctx))

	be.NilErr(t, rb.Clone().
		BodyForm(url.Values{
			"bot-field": []string{},
			"host_page": []string{"host_page2"},
			"name":      []string{"name2"},
			"email":     []string{"email2"},
			"CC":        []string{"CC2"},
			"subject":   []string{"subject2"},
			"anonymous": []string{"1"},
			"comment":   []string{"comment2"},
		}).
		CheckStatus(303).
		Fetch(ctx))

	be.NilErr(t, rb.Clone().
		BodyForm(url.Values{
			"anonymous": []string{"XXX"},
		}).
		CheckStatus(400).
		Fetch(ctx))

	comments, err := svc.q.ListComments(ctx, db.ListCommentsParams{
		Limit:  3,
		Offset: 0,
	})
	be.NilErr(t, err)
	for i := range comments {
		clearComment(&comments[i])
	}
	testfile.EqualJSON(t, "testdata/comments.json", comments)
}
