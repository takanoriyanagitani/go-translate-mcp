package otranslate

import (
	"bufio"
	"context"
	_ "embed"
	"io"
	"os"
	"strings"
	"text/template"

	oapi "github.com/ollama/ollama/api"
	eif "github.com/takanoriyanagitani/go-translate-mcp"
	disp "golang.org/x/text/language/display"
)

//go:embed prompt.tmpl
var tmplStrDefault string

var tmplDefault *template.Template = template.
	Must(
		template.New("default").Parse(tmplStrDefault),
	).
	Option("missingkey=error")

type Client struct{ *oapi.Client }

type Config struct {
	// Filepath to the optional template file.
	TemplatePath string

	// Name of the optional template.
	TemplateName string

	// Max size of the template.
	TemplateLimit int64

	Model string
}

func (c Config) ToTemplate() (*template.Template, error) {
	if 0 == len(c.TemplatePath) {
		return tmplDefault, nil
	}

	return path2template(c.TemplateName, c.TemplatePath, c.TemplateLimit)
}

func (c Config) ToUserPrompt(input eif.Input) (string, error) {
	tmpl, e := c.ToTemplate()
	if nil != e {
		return "", e
	}

	return InputToUserPrompt(tmpl, input)
}

func (c Config) ToChatRequest(input eif.Input) (oapi.ChatRequest, error) {
	prompt, e := c.ToUserPrompt(input)

	return PromptToChatRequest(
		prompt,
		c.Model,
	), e
}

func (c Client) ToEngine(cfg Config) eif.Engine {
	return func(ctx context.Context, i eif.Input) (eif.Output, error) {
		req, e := cfg.ToChatRequest(i)
		if nil != e {
			return eif.Output{}, e
		}

		translated, e := c.ReqToRes(ctx, &req)
		return eif.Output{
			Text: translated,
		}, e
	}
}

func (c Client) ReqToRes(ctx context.Context, req *oapi.ChatRequest) (string, error) {
	var bldr strings.Builder
	e := c.Client.Chat(
		ctx,
		req,
		func(res oapi.ChatResponse) error {
			var content string = res.Message.Content
			bldr.WriteString(content) // always nil error
			return nil
		},
	)

	return bldr.String(), e
}

func str2tmpl(name string, tmpl string) (*template.Template, error) {
	return template.New(name).Parse(tmpl)
}

func reader2string(rdr io.Reader, limit int64) (string, error) {
	var br io.Reader = bufio.NewReader(rdr)
	taken := &io.LimitedReader{R: br, N: limit}

	var buf strings.Builder

	_, e := io.Copy(&buf, taken)
	if nil != e {
		return "", e
	}

	return buf.String(), nil
}

func path2string(trustedPath string, limit int64) (string, error) {
	f, e := os.Open(trustedPath) //nolint:gosec
	if nil != e {
		return "", e
	}
	defer f.Close() //nolint:errcheck

	return reader2string(bufio.NewReader(f), limit)
}

func path2template(name string, trustedPath string, limit int64) (*template.Template, error) {
	s, e := path2string(trustedPath, limit)
	if nil != e {
		return nil, e
	}

	return str2tmpl(name, s)
}

func InputToUserPrompt(
	tmpl *template.Template,
	input eif.Input,
) (string, error) {
	var buf strings.Builder
	e := tmpl.Execute(&buf, map[string]string{
		"SOURCE_LANG": disp.English.Languages().Name(input.Source),
		"SOURCE_CODE": input.Source.String(),

		"TARGET_LANG": disp.English.Languages().Name(input.Target),
		"TARGET_CODE": input.Target.String(),

		"TEXT": input.Text,
	})
	return buf.String(), e
}

func PromptToChatRequest(
	prompt string,
	model string,
) oapi.ChatRequest {
	var strm bool = false
	return oapi.ChatRequest{
		Model: model,
		Messages: []oapi.Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: &strm,
	}
}
