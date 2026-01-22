package interceptors

import (
	"community-backend/pkg/xcode"
	"context"

	"google.golang.org/grpc"
)

func ClientErrorInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {

		}
	}
}
