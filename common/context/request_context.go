package requestcontext

import "context"

// Custom key untuk Request ID.
type contextKey string

const requestIDKey contextKey = "request_id"

// WithRequestID menyimpan Request ID ke Go Context.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// RequestID mengambil Request ID dari Go Context.
func RequestID(ctx context.Context) string {
	requestID, _ := ctx.Value(requestIDKey).(string)
	return requestID
}
