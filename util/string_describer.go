package util

import "context"

type StringDescriber struct{ Value string }

func (d StringDescriber) MarkdownDescription(context.Context) string {
	return d.Value
}

func (d StringDescriber) Description(ctx context.Context) string {
	return d.MarkdownDescription(ctx)
}
