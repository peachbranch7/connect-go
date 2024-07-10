package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/bufbuild/connect-go"
	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	greetv1 "connect-go/gen/greet/v1"
	"connect-go/gen/greet/v1/greetv1connect"
)

type GreetServer struct{}

func (s *GreetServer) Greet(ctx context.Context, req *connect.Request[greetv1.GreetRequest]) (*connect.Response[greetv1.GreetResponse], error) {
	log.Println("Request headers: ", req.Header())

	if req.Msg.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("name is required."))
	}

	greetResp := &greetv1.GreetResponse{
		Greeting: fmt.Sprintf("Hello, %s!", req.Msg.Name),
	}
	resp := connect.NewResponse(greetResp)
	resp.Header().Set("Greet-Version", "v1")
	return resp, nil
}

// リフレクション
func newServeMuxWithReflection() *http.ServeMux {
	mux := http.NewServeMux()
	reflector := grpcreflect.NewStaticReflector(
		"greet.v1.GreetService",
	)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
	return mux
}

// インターセプタ
func newInterCeptors() connect.Option {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			req.Header().Set("hoge", "fuga")
			return next(ctx, req)
		})
	}
	return connect.WithInterceptors(connect.UnaryInterceptorFunc(interceptor))
}

func main() {
	greetServer := &GreetServer{}

	mux := newServeMuxWithReflection()
	interceptor := newInterCeptors()
	path, handler := greetv1connect.NewGreetServiceHandler(greetServer, interceptor)
	mux.Handle(path, handler)
	http.ListenAndServe(":8080", h2c.NewHandler(mux, &http2.Server{})) // Use h2c so we can serve HTTP/2 without TLS.
}