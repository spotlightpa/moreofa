// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"time"
)

type Comment struct {
	ID        int64
	Name      string
	Contact   string
	Message   string
	Ip        string
	UserAgent string
	Referrer  string
	HostPage  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
