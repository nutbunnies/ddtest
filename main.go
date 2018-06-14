package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	fmt.Println("Starting tracer")
	tracer.Start(
		tracer.WithAgentAddr("localhost:8126"),
		tracer.WithDebugMode(true),
		tracer.WithServiceName("test_app"),
		tracer.WithGlobalTag("env", "local_dev"),
		tracer.WithGlobalTag("runtime", "go"),
		tracer.WithGlobalTag("version", "1.0123"),
	)

	for i := 0; i < 10; i++ {
		fmt.Println("Request: ", i)
		fakeMiddleware()
	}

	fmt.Println("stopping tracer")
	tracer.Stop()
}

func fakeMiddleware() {
	//Using context.Context to mimic request context
	ctx := context.Background()
	opts := []ddtrace.StartSpanOption{
		tracer.SpanType(ext.AppTypeWeb),
		tracer.ServiceName("testapp"),
		tracer.Tag(ext.HTTPMethod, "GET"),
		tracer.Tag(ext.HTTPURL, "/some/path"),
	}
	span, ctx := tracer.StartSpanFromContext(ctx, "foo.request", opts...)
	defer span.Finish()
	//would be passed thru middleware here

	fakeHandler(ctx)

	//return path thru middleware
	span.SetTag(ext.ResourceName, "/some/{var}") // using chi router so don't have this info until after return from servehttp
	span.SetTag("request_id", "random-guid")
	span.SetTag(ext.HTTPCode, "200")
}

func fakeHandler(ctx context.Context) {
	span, _ := tracer.StartSpanFromContext(ctx, "tracksomething", tracer.SpanType(ext.AppTypeWeb))
	defer span.Finish()
	rand.Seed(int64(time.Now().Nanosecond()))
	time.Sleep(time.Duration(1*rand.Int63n(5000)) * time.Millisecond)
}
