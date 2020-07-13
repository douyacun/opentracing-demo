package main

import (
	"context"
	greet "demo/zipkin/proto"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"google.golang.org/grpc"
	"log"
	"net"
)

type handler struct{}

func (h *handler) Say(ctx context.Context, req *greet.SayRequest) (*greet.SayResponse, error) {
	res := &greet.SayResponse{
		Msg:  "success",
		Code: 0,
	}
	log.Printf("%s say %s", req.Name, req.Msg)

	return res, nil
}

func main() {
	// reporter
	reporter := zipkinhttp.NewReporter("http://localhost:9411/api/v2/spans")
	defer reporter.Close()
	// endpoint
	endpoint, err := zipkin.NewEndpoint("opentracing-demo", "0.0.0.0")
	if err != nil {
		log.Fatalf("zipkin new endpoint error: %v", err)
	}
	// tracer
	nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		log.Fatalf("zipkin new tracer error: %v", err)
	}

	tracer := zipkinot.Wrap(nativeTracer)
	opentracing.SetGlobalTracer(tracer)
	grpcServer := grpc.NewServer(
		// otgrpc.LogPayloads 是否记录 入参和出参
		// otgrpc.SpanDecorator 装饰器，回调函数
		// otgrpc.IncludingSpans 是否记录
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(otgrpc.OpenTracingServerInterceptor(tracer,
			otgrpc.LogPayloads(),
			// IncludingSpans是请求前回调
			otgrpc.IncludingSpans(func(parentSpanCtx opentracing.SpanContext, method string, req, resp interface{}) bool {
				log.Printf("method: %s", method)
				log.Printf("req: %+v", req)
				log.Printf("resp: %+v", resp)
				return true
			}),
			// SpanDecorator是请求后回调
			otgrpc.SpanDecorator(func(span opentracing.Span, method string, req, resp interface{}, grpcError error) {
				log.Printf("method: %s", method)
				log.Printf("req: %+v", req)
				log.Printf("resp: %+v", resp)
				log.Printf("grpcError: %+v", grpcError)
			}),
		))),
	)

	greet.RegisterServiceServer(grpcServer, &handler{})
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("net listen error: %v", err.Error())
	}
	grpcServer.Serve(listener)
}
