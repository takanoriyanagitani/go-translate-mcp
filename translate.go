package translate

import (
	"context"
	"errors"

	langparse "github.com/takanoriyanagitani/go-translate-mcp/language"
	lang "golang.org/x/text/language"
)

type Input struct {
	Source lang.Tag
	Target lang.Tag
	Text   string
}

type Output struct {
	Text string
}

type Engine func(context.Context, Input) (Output, error)

type RawInput struct {
	Source string
	Target string
	Text   string
}

func (r RawInput) ToSourceTag() (lang.Tag, error) { return langparse.Parse(r.Source) }
func (r RawInput) ToTargetTag() (lang.Tag, error) { return langparse.Parse(r.Target) }

func (r RawInput) ToInput() (Input, error) {
	ts, es := r.ToSourceTag()
	tt, et := r.ToTargetTag()
	return Input{
		Source: ts,
		Target: tt,
		Text:   r.Text,
	}, errors.Join(es, et)
}

type Translate func(context.Context, RawInput) (Output, error)

func (eng Engine) ToTranslate() Translate {
	return func(ctx context.Context, raw RawInput) (Output, error) {
		inp, e := raw.ToInput()
		if nil != e {
			return Output{}, e
		}

		return eng(ctx, inp)
	}
}
