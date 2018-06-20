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

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

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
	span.SetTag(ext.HTTPCode, "200")
	fmt.Printf("\n======\nParent span: %v\n\n======\n", span)
}

func fakeHandler(ctx context.Context) {
	span, ctx := tracer.StartSpanFromContext(ctx, "tracksomething", tracer.SpanType(ext.AppTypeWeb))
	defer span.Finish()
	rand.Seed(int64(time.Now().Nanosecond()))
	time.Sleep(time.Duration(1*rand.Int63n(2000)) * time.Millisecond)
	fakeFunc(ctx)
	fmt.Printf("\n======\nChild span: %v\n\n======\n", span)
}

func fakeFunc(ctx context.Context) {
	span, _ := tracer.StartSpanFromContext(ctx, "twolevelsdown", tracer.SpanType(ext.AppTypeCache))
	defer span.Finish()
	time.Sleep(time.Duration(1*rand.Int63n(2000)) * time.Millisecond)
	fmt.Printf("\n======\nGrandChild span: %v\n\n======\n", span)
}
