// Package slogdiscard provides a no-op implementation of slog.Handler
// that discards all log output. Useful for testing or disabling logging.
package slogdiscard

import (
	"context"
	"log/slog"
)

// NewDiscardLogger creates a new Logger that discards all log output.
func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

// DiscardHandler implements slog.Handler that discards all log records.
type DiscardHandler struct{}

// NewDiscardHandler creates a new DiscardHandler instance.
func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

// Handle discards the log record (no-op implementation).
func (h *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	// Silently discard the log record
	return nil
}

// WithAttrs returns the same handler since we don't store attributes.
func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	// No attributes to store, return same handler
	return h
}

// WithGroup returns the same handler since we don't support groups.
func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	// No group support needed, return same handler
	return h
}

// Enabled always returns false since we discard all logs.
func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	// Logging is never enabled for this handler
	return false
}
