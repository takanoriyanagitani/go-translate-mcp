package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	oapi "github.com/ollama/ollama/api"
	eif "github.com/takanoriyanagitani/go-translate-mcp"
	oeng "github.com/takanoriyanagitani/go-translate-mcp/engine/ai/ollama"
)

const (
	defaultPort         = 12041
	readTimeoutSeconds  = 10
	writeTimeoutSeconds = 10
	maxHeaderExponent   = 20
)

var (
	port = flag.Int("port", defaultPort, "port to listen")
	mdl  = flag.String("model", "translategemma:4b-it-q4_K_M", "translation model")
)

func main() {
	flag.Parse()

	cli, e := oapi.ClientFromEnvironment()
	if nil != e {
		log.Fatalf("unable to create ollama client: %v", e)
	}

	oc := oeng.Client{Client: cli}
	var eng eif.Engine = oc.ToEngine(oeng.Config{
		TemplatePath:  "", // use default
		TemplateName:  "default",
		TemplateLimit: 1048576,
		Model:         *mdl,
	})

	var tra eif.Translate = eng.ToTranslate()

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "go-translate",
		Version: "v0.1.0",
		Title:   "Go Translate",
	}, nil)

	translateTool := func(ctx context.Context, req *mcp.CallToolRequest, input eif.RawInput) (
		*mcp.CallToolResult,
		eif.Output,
		error,
	) {
		output, err := tra(ctx, input)
		return nil, output, err
	}

	mcp.AddTool(server, &mcp.Tool{
		Name:         "translate",
		Title:        "Translate",
		Description:  "Translate text from a source language to a target language.",
		InputSchema:  nil, // Inferred by AddTool
		OutputSchema: nil, // Inferred by AddTool
	}, translateTool)

	address := fmt.Sprintf(":%d", *port)

	mcpHandler := mcp.NewStreamableHTTPHandler(
		func(req *http.Request) *mcp.Server { return server },
		&mcp.StreamableHTTPOptions{Stateless: true},
	)

	httpServer := &http.Server{
		Addr:           address,
		Handler:        mcpHandler,
		ReadTimeout:    readTimeoutSeconds * time.Second,
		WriteTimeout:   writeTimeoutSeconds * time.Second,
		MaxHeaderBytes: 1 << maxHeaderExponent,
	}

	log.Printf("Ready to start HTTP MCP server. Listening on %s\n", address)
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Printf("Failed to listen and serve: %v\n", err)
		return
	}
}
