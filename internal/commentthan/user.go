package commentthan

import (
	"context"
	"encoding/gob"
	"net/http"

	"github.com/spotlightpa/moreofa/internal/clogger"
)

type User struct {
	Role []string
	// Email: The user's email address.
	Email string `json:"email,omitempty"`
	// FamilyName: The user's last name.
	FamilyName string `json:"family_name,omitempty"`
	// Gender: The user's gender.
	Gender string `json:"gender,omitempty"`
	// GivenName: The user's first name.
	GivenName string `json:"given_name,omitempty"`
	// Hd: The hosted domain e.g. example.com if the user is Google apps user.
	Hd string `json:"hd,omitempty"`
	// Id: The obfuscated ID of the user.
	Id string `json:"id,omitempty"`
	// Link: URL of the profile page.
	Link string `json:"link,omitempty"`
	// Locale: The user's preferred locale.
	Locale string `json:"locale,omitempty"`
	// Name: The user's full name.
	Name string `json:"name,omitempty"`
	// Picture: URL of the user's picture image.
	Picture string `json:"picture,omitempty"`
}

func init() {
	gob.Register(User{})
}

type userCtxKey int

func ReqWithUser(r *http.Request, user User) *http.Request {
	ctx := context.WithValue(r.Context(), userCtxKey(0), user)
	l := clogger.FromContext(ctx).
		With("user.name", user.Name, "user.email", user.Email)
	ctx = clogger.NewContext(ctx, l)
	r2 := r.WithContext(ctx)
	return r2
}

func UserFromReq(r *http.Request) User {
	ctx := r.Context()
	user, ok := ctx.Value(userCtxKey(0)).(User)
	if !ok {
		panic("session middleware misconfigured")
	}
	return user
}
