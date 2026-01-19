package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	oapi "github.com/ollama/ollama/api"
	eif "github.com/takanoriyanagitani/go-translate-mcp"
	oeng "github.com/takanoriyanagitani/go-translate-mcp/engine/ai/ollama"
)

func check(e error) {
	if nil == e {
		return
	}
	log.Fatalf("unable to run: %v", e)
}

func tryReadStdin(stdin io.Reader) (string, error) {
	var builder strings.Builder
	_, e := io.Copy(&builder, stdin)
	return builder.String(), e
}

func main() {
	var src string
	var tgt string
	var txt string
	var mdl string
	var tmplPath string
	var tmplName string
	var tmplLimit int64

	flag.StringVar(&src, "source", "en", "source language")
	flag.StringVar(&src, "s", "en", "source language (shorthand)")
	flag.StringVar(&tgt, "target", "ja", "target language")
	flag.StringVar(&tgt, "t", "ja", "target language (shorthand)")
	flag.StringVar(&txt, "text", "", "text to translate")
	flag.StringVar(&txt, "x", "", "text to translate (shorthand)")
	flag.StringVar(&mdl, "model", "translategemma:4b-it-q4_K_M", "translation model")
	flag.StringVar(&mdl, "m", "translategemma:4b-it-q4_K_M", "translation model (shorthand)")
	flag.StringVar(&tmplPath, "template-path", "", "custom prompt template path")
	flag.StringVar(&tmplPath, "p", "", "custom prompt template path (shorthand)")
	flag.StringVar(&tmplName, "template-name", "default", "custom prompt template name")
	flag.StringVar(&tmplName, "n", "default", "custom prompt template name (shorthand)")
	flag.Int64Var(&tmplLimit, "template-limit", 1048576, "custom prompt template limit")
	flag.Int64Var(&tmplLimit, "l", 1048576, "custom prompt template limit (shorthand)")
	flag.Parse()

	if len(txt) < 1 {
		stdin, e := tryReadStdin(os.Stdin)
		check(e)
		txt = stdin
	}

	cli, e := oapi.ClientFromEnvironment()
	check(e)

	oc := oeng.Client{Client: cli}
	var eng eif.Engine = oc.ToEngine(oeng.Config{
		TemplatePath:  tmplPath,
		TemplateName:  tmplName,
		TemplateLimit: tmplLimit,
		Model:         mdl,
	})

	var tra eif.Translate = eng.ToTranslate()

	input := eif.RawInput{
		Source: src,
		Target: tgt,
		Text:   txt,
	}

	out, e := tra(context.Background(), input)
	check(e)

	fmt.Printf("%s", out.Text)
}
