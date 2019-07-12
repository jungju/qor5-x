package vuetify

import (
	"context"
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VListBuilder struct {
	tag *h.HTMLTagBuilder
}

func VList(children ...h.HTMLComponent) (r *VListBuilder) {
	r = &VListBuilder{
		tag: h.Tag("v-list").Children(children...),
	}
	return
}

func (b *VListBuilder) Dark(v bool) (r *VListBuilder) {
	b.tag.Attr(":dark", fmt.Sprint(v))
	return b
}

func (b *VListBuilder) Dense(v bool) (r *VListBuilder) {
	b.tag.Attr(":dense", fmt.Sprint(v))
	return b
}

func (b *VListBuilder) Expand(v bool) (r *VListBuilder) {
	b.tag.Attr(":expand", fmt.Sprint(v))
	return b
}

func (b *VListBuilder) Light(v bool) (r *VListBuilder) {
	b.tag.Attr(":light", fmt.Sprint(v))
	return b
}

func (b *VListBuilder) Subheader(v bool) (r *VListBuilder) {
	b.tag.Attr(":subheader", fmt.Sprint(v))
	return b
}

func (b *VListBuilder) ThreeLine(v bool) (r *VListBuilder) {
	b.tag.Attr(":three-line", fmt.Sprint(v))
	return b
}

func (b *VListBuilder) TwoLine(v bool) (r *VListBuilder) {
	b.tag.Attr(":two-line", fmt.Sprint(v))
	return b
}

func (b *VListBuilder) MarshalHTML(ctx context.Context) (r []byte, err error) {
	return b.tag.MarshalHTML(ctx)
}
