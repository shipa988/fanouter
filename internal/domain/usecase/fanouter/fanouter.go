package fanouter

import (
	"context"
)

// Fanouter is abstract object receiving incoming feed id and transmitting multi queries to external urls.
type Fanouter interface {
	Fanout(ctx context.Context, id string) error
	Init(ctx context.Context) error
}
