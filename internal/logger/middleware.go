package logger

import (
	"context"
	"fmt"
	"log/slog"
)

// HandlerMiddleware adds context fields to every log record.
type HandlerMiddleware struct {
	next slog.Handler
}

// NewHandlerMiddleware wraps the given handler with additional context.
func NewHandlerMiddleware(next slog.Handler) *HandlerMiddleware {
	return &HandlerMiddleware{
		next: next,
	}
}

// Enabled forwards the Enabled check to the next handler.
func (h *HandlerMiddleware) Enabled(ctx context.Context, rec slog.Level) bool {
	return h.next.Enabled(ctx, rec)
}

// Handle enriches the record with context information before logging.
func (h *HandlerMiddleware) Handle(ctx context.Context, rec slog.Record) error {
	c, ok := ctx.Value(key).(logCtx)
	if !ok {
		return h.next.Handle(ctx, rec)
	}
	if c.Component != "" {
		rec.Add("component", c.Component)
	}
	if c.Method != "" {
		rec.Add("method", c.Method)
	}
	if c.Login != "" {
		rec.Add("login", c.Login)
	}
	if c.Password != "" {
		rec.Add("password", c.Password)
	}
	if c.IP != "" {
		rec.Add("ip", c.IP)
	}
	if c.Network != "" {
		rec.Add("network", c.Network)
	}
	return h.next.Handle(ctx, rec)
}

// WithAttrs returns a new handler with additional attributes.
func (h *HandlerMiddleware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &HandlerMiddleware{
		next: h.next.WithAttrs(attrs),
	}
}

// WithGroup returns a new handler with the given group name.
func (h *HandlerMiddleware) WithGroup(name string) slog.Handler {
	return &HandlerMiddleware{
		next: h.next.WithGroup(name),
	}
}

type logCtx struct {
	Component string
	Method    string
	Login     string
	Password  string
	IP        string
	Network   string
}

type keyType int

const key = keyType(0)

// WithLogComponent attaches a component name to the logging context.
func WithLogComponent(ctx context.Context, component string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.Component = component
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{
		Component: component,
	})
}

// WithLogMethod attaches a method name to the logging context.
func WithLogMethod(ctx context.Context, method string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.Method = method
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{
		Method: method,
	})
}

// WithLogLogin attaches a login to the logging context.
func WithLogLogin(ctx context.Context, login string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.Login = login
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{
		Login: login,
	})
}

// WithLogPassword attaches a password to the logging context.
func WithLogPassword(ctx context.Context, password string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.Password = password
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{
		Password: password,
	})
}

// WithLogIP attaches an IP address to the logging context.
func WithLogIP(ctx context.Context, ip string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.IP = ip
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{
		IP: ip,
	})
}

// WithLogNetwork attaches a network to the logging context.
func WithLogNetwork(ctx context.Context, network string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.Network = network
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{
		Network: network,
	})
}

// AddPrefix adds a prefix to the error message based on the context.
func AddPrefix(ctx context.Context, err error) error {
	var prefix string

	// Extract component from context
	c := logCtx{}
	// Extract method from context and append to prefix if present
	if x, ok := ctx.Value(key).(logCtx); ok {
		c = x
	}

	if c.Component != "" {
		prefix = c.Component
	}
	if c.Method != "" {
		if prefix != "" {
			prefix += "." + c.Method
		} else {
			prefix = c.Method
		}
	}

	// Wrap the error with prefix if available
	if prefix != "" {
		err = fmt.Errorf("%s: %w", prefix, err)
	}

	return err
}

// type errorWithCtx struct {
// 	next error
// 	ctx  logCtx
// }

// func (e *errorWithCtx) Unwrap() error {
// 	return e.next
// }

// func (e *errorWithCtx) Error() string {
// 	return e.next.Error()
// }

// // ErrorCtx extracts logging context from a wrapped error.
// func ErrorCtx(ctx context.Context, err error) context.Context {
// 	var errWithCtx *errorWithCtx
// 	if errors.As(err, &errWithCtx) {
// 		return context.WithValue(ctx, key, errWithCtx.ctx)
// 	}
// 	return ctx
// }
