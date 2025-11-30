package lib

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"

	"github.com/sodefrin/PP/server/db"
)

type TraceHandler struct {
	slog.Handler
}

func NewTraceHandler(h slog.Handler) *TraceHandler {
	return &TraceHandler{Handler: h}
}

func (h *TraceHandler) Handle(ctx context.Context, r slog.Record) error {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		traceID := span.SpanContext().TraceID().String()
		r.AddAttrs(slog.String("trace_id", traceID))
	}
	return h.Handler.Handle(ctx, r)
}

type UserHandler struct {
	slog.Handler
}

func NewUserHandler(h slog.Handler) *UserHandler {
	return &UserHandler{Handler: h}
}

func (h *UserHandler) Handle(ctx context.Context, r slog.Record) error {
	if user, ok := ctx.Value(UserContextKey).(db.User); ok {
		r.AddAttrs(slog.Int64("user_id", user.ID))
	}
	return h.Handler.Handle(ctx, r)
}
