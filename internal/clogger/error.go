package clogger

import "context"

func LogErr(ctx context.Context, err error) {
	FromContext(ctx).ErrorContext(ctx, "error", "error", err)
}
