package controllers

import (
	"context"

	"github.com/rs/xid"
	"google.golang.org/grpc"
)

// WrappedServerStream is a grpc.ServerStream wrapper.
type WrappedServerStream struct {
	grpc.ServerStream
	context context.Context
}

// Context override Context function of grpc.ServerStream.
func (wss *WrappedServerStream) Context() context.Context {
	return wss.context
}

// ClientIDSetter middleware sets a client id on the incoming stream,
// saves in stream context.
func ClientIDSetter(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := ss.Context()

	// Because the handler function receives a grpc.ServerStream interface,
	// so we can use our own struct to wrap original grpc.ServerStream,
	// to replace context with our own. Out own context contain a client id,
	// we can extract this value in the controller.
	//
	// You can find another usage like this on GitHub:
	// https://github.com/grpc-ecosystem/go-grpc-middleware/blob/master/auth/auth.go
	wrapped := &WrappedServerStream{
		ServerStream: ss,
		context:      context.WithValue(ctx, "client_id", xid.New().String()),
	}
	return handler(srv, wrapped)
}

// ExtractClientID find client id from context.
func ExtractClientID(ctx context.Context) string {
	return ctx.Value("client_id").(string)
}
