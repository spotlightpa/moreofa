package akismet

import (
	"context"
	"errors"
	"maps"
	"net/http"
	"net/url"
	"time"

	"github.com/carlmjohnson/requests"
)

type Client struct {
	site, key string
	cl        *http.Client
}

func New(site, key string, cl *http.Client) *Client {
	return &Client{
		site,
		key,
		cl,
	}
}

func (cl Client) new(rb *requests.Builder) {
	rb.BaseURL("https://rest.akismet.com/1.1/").
		Client(cl.cl)
}

func (cl Client) insertAuth(values url.Values) {
	values.Set("api_key", cl.key)
	values.Set("blog", cl.site)
}

func (cl Client) Verify(ctx context.Context) (bool, error) {
	values := make(url.Values)
	cl.insertAuth(values)
	var body string
	err := requests.New(cl.new).
		Path("verify-key").
		BodyForm(values).
		ToString(&body).
		Fetch(ctx)
	if err != nil {
		return false, err
	}
	return body == "valid", nil
}

type Comment struct {
	// Homepage URL of the website being commented on.
	// Defaults to the site used for key verification.
	Site string

	// IP address of the commenter.
	UserIP string

	// User agent string of the commenter's browser.
	UserAgent string

	// The HTTP_REFERER header.
	Referer string

	// URL of the page being commented on.
	Page string

	// Publish date/time of the page being commented on.
	PageTimestamp time.Time

	// Name of the commenter.
	Author string

	// Email address of the commenter.
	AuthorEmail string

	// Website of the commenter.
	AuthorSite string

	// Content type, e.g. "comment", "forum-post".
	// See https://blog.akismet.com/2012/06/19/pro-tip-tell-us-your-comment_type/
	// for more examples.
	Type string

	// Content of the comment. May contain HTML.
	Content string

	// Publish date/time of the comment. Akismet uses
	// the current time if one is not specified.
	Timestamp time.Time

	// Comma-separated list of languages in use on the
	// website being commented on, e.g. "en, fr_ca".
	SiteLanguage string

	// Character encoding for the website being commented
	// on, e.g. "UTF-8".
	SiteCharset string

	Context       []string
	IsTest        bool
	Honeypot      string
	RecheckReason string
}

func (c Comment) Values() url.Values {
	var timestamp, pageTimestamp string
	if !c.Timestamp.IsZero() {
		timestamp = c.Timestamp.UTC().Format(time.RFC3339)
	}
	if !c.PageTimestamp.IsZero() {
		pageTimestamp = c.PageTimestamp.UTC().Format(time.RFC3339)
	}

	values := url.Values{
		"comment_context[]":         c.Context,
		"recheck_reason":            []string{c.RecheckReason},
		"blog":                      []string{c.Site},
		"blog_charset":              []string{c.SiteCharset},
		"blog_lang":                 []string{c.SiteLanguage},
		"comment_author":            []string{c.Author},
		"comment_author_email":      []string{c.AuthorEmail},
		"comment_author_url":        []string{c.AuthorSite},
		"comment_content":           []string{c.Content},
		"comment_date_gmt":          []string{timestamp},
		"comment_post_modified_gmt": []string{pageTimestamp},
		"comment_type":              []string{c.Type},
		"permalink":                 []string{c.Page},
		"referrer":                  []string{c.Referer},
		"user_agent":                []string{c.UserAgent},
		"user_ip":                   []string{c.UserIP},
	}
	if c.IsTest {
		values.Set("user_role", "administrator")
		values.Set("is_test", "true")
	}
	if c.Honeypot != "" {
		values.Set("honeypot_field_name", "hidden_honeypot_field")
		values.Set("hidden_honeypot_field", c.Honeypot)
	}

	maps.DeleteFunc(values, func(key string, value []string) bool {
		return values.Get(key) == ""
	})
	return values
}

//go:generate stringer -type=CommentKind
type CommentKind int8

const (
	UnknownKind CommentKind = iota
	HamKind
	SpamKind
	TrashKind
)

func (cl Client) Check(ctx context.Context, c Comment) (CommentKind, error) {
	values := c.Values()
	cl.insertAuth(values)
	var body string
	headers := make(url.Values)
	err := requests.
		New(cl.new).
		Path("comment-check").
		BodyForm(values).
		ToString(&body).
		CopyHeaders(headers).
		CheckStatus(http.StatusOK).
		Fetch(ctx)
	if err != nil {
		return UnknownKind, err
	}
	if body != "true" {
		return HamKind, nil
	}
	if headers.Get("X-Akismet-Pro-Tip") == "discard" {
		return TrashKind, nil
	}
	return SpamKind, nil
}

const (
	TypeComment     = "comment"
	TypeForumPost   = "forum‑post"
	TypeReply       = "reply"
	TypeBlogPost    = "blog‑post"
	TypeContactForm = "contact‑form"
	TypeSignup      = "signup"
	TypeMessage     = "message"
)

func (cl Client) SubmitHam(ctx context.Context, c Comment) error {
	values := c.Values()
	cl.insertAuth(values)
	var body string
	err := requests.
		New(cl.new).
		Path("submit-ham").
		BodyForm(values).
		ToString(&body).
		CheckStatus(http.StatusOK).
		Fetch(ctx)
	if err != nil {
		return err
	}
	if body != "Thanks for making the web a better place." {
		return errors.New("unexpected response")
	}
	return nil
}

func (cl Client) SubmitSpam(ctx context.Context, c Comment) error {
	values := c.Values()
	cl.insertAuth(values)
	var body string
	err := requests.
		New(cl.new).
		Path("submit-spam").
		BodyForm(values).
		ToString(&body).
		CheckStatus(http.StatusOK).
		Fetch(ctx)
	if err != nil {
		return err
	}
	if body != "Thanks for making the web a better place." {
		return errors.New("unexpected response")
	}
	return nil
}
