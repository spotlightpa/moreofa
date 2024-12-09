package akismet_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/carlmjohnson/be"
	"github.com/carlmjohnson/requests/reqtest"
	"github.com/spotlightpa/moreofa/internal/akismet"
)

func TestClient(t *testing.T) {
	tb := be.Relaxed(t)
	ctx := context.Background()
	cl := &http.Client{
		Transport: reqtest.Caching(nil, "testdata"),
	}
	akcl := akismet.New("https://www.spotlightpa.org", "abc123", cl)
	ok, err := akcl.Verify(ctx)
	be.NilErr(tb, err)
	be.True(tb, ok)

	kind, err := akcl.Check(ctx, akismet.Comment{
		Content: "I have a tip about Pennsylvannia political corruption.",
		Context: []string{"Pennsylvannia"},
		Type:    akismet.TypeContactForm,
	})
	be.NilErr(tb, err)
	be.Equal(tb, akismet.HamKind, kind)

	kind, err = akcl.Check(ctx, akismet.Comment{
		Content:  `I hope this message finds you well. At Risk-FreewithMoney-BackGuarantee.com, we understand that navigating the financial landscape can be challenging, and finding opportunities with real potential is no small feat. That’s why we’re excited to present you with a unique investment opportunity that combines transparency, security, and unparalleled growth potential.`,
		Context:  []string{"Pennsylvannia"},
		Type:     akismet.TypeContactForm,
		Honeypot: "true",
	})
	be.NilErr(tb, err)
	be.Equal(tb, akismet.SpamKind, kind)

	kind, err = akcl.Check(ctx, akismet.Comment{
		IsTest: true,
	})
	be.NilErr(tb, err)
	be.Equal(tb, akismet.HamKind, kind)

	err = akcl.SubmitHam(ctx, akismet.Comment{
		IsTest: true,
	})
	be.NilErr(tb, err)

	err = akcl.SubmitSpam(ctx, akismet.Comment{
		IsTest: true,
	})
	be.NilErr(tb, err)
}
