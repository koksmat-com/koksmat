package service

import (
	"context"
	stdlog "log" // Renamed log to stdlog
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log" // Renamed to otlog
	jaegerconfig "github.com/uber/jaeger-client-go/config"
)

func director(req *http.Request) {
	stdlog.Println("Request received")
	// You can inspect or modify the request here before it's sent to the target server
}

func modifyResponse(res *http.Response) error {
	// You can inspect or modify the response here before it's sent back to the client
	return nil
}

func handleRequestAndRedirect(w http.ResponseWriter, req *http.Request) {
	// Replace "https://example.com" with the URL of the HTTPS service you want to proxy to
	targetURL, err := url.Parse("https://example.com")
	if err != nil {
		stdlog.Fatal(err) // Use stdlog here
	}

	// Create a custom reverse proxy with Director and ModifyResponse functions
	proxy := &httputil.ReverseProxy{
		Director:       director,
		ModifyResponse: modifyResponse,
	}

	// Update the headers to allow for SSL redirection
	req.URL.Host = targetURL.Host
	req.URL.Scheme = targetURL.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = targetURL.Host

	// Note: this is not strictly necessary, but it is good practice
	req.Header.Set("Accept-Encoding", "")

	// Start tracing
	span := opentracing.StartSpan("proxy")
	defer span.Finish()

	// Inject the span context into the HTTP headers
	opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))

	// Serve the request by proxying to the target URL
	proxy.ServeHTTP(w, req)

	// Log any errors encountered during the proxying process
	if err := req.Context().Err(); err != nil && err != context.Canceled {
		span.SetTag("error", true)
		span.LogFields(otlog.String("event", "error"), otlog.String("message", err.Error())) // Use otlog here
	}
}

func Serve() {
	// Initialize Jaeger tracer
	cfg, err := jaegerconfig.FromEnv()
	if err != nil {
		stdlog.Fatal(err) // Use stdlog here
	}
	cfg.ServiceName = "proxy"
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		stdlog.Fatal(err) // Use stdlog here
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	// Start the HTTP server
	http.HandleFunc("/", handleRequestAndRedirect)
	exitCode := http.ListenAndServe(":8080", nil)
	stdlog.Fatal(exitCode) // Use stdlog here
}
