// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"time"
)

type Comment struct {
	ID         int64
	Name       string
	Contact    string
	Subject    string
	Cc         string
	Message    string
	Ip         string
	UserAgent  string
	Referrer   string
	HostPage   string
	CreatedAt  time.Time
	ModifiedAt time.Time
}

type Session struct {
	Token  string
	Data   []byte
	Expiry float64
}
