// Copyright 2018 Yi Jin. All rights reserved.
// license that can be found in the LICENSE file.

package main

/*
Comment interface
*/

func NewComment(tp string, attr map[string]string) (Comment, bool) {
	switch tp {
	case "text":
		return NewTextCommentFromMap(attr)
	default:
		return nil, false
	}
}

type Comment interface {
	Type() string
	Content() string
	Attributes() map[string]string
}

/*
Text Comment
*/

func NewTextCommentFromMap(attr map[string]string) (*TextComment, bool) {
	text, ok := attr["text"]
	if !ok {
		return nil, false
	}
	color, ok := attr["color"]
	if !ok {
		return nil, false
	}
	if len(text) == 0 {
		return nil, false
	}
	return NewTextComment(text, color), true
}

func NewTextComment(text string, color string) *TextComment {
	return &TextComment{
		Text:  text,
		Color: color,
	}
}

type TextComment struct {
	Text  string
	Color string
}

func (c *TextComment) Type() string {
	return "text"
}

func (c *TextComment) Content() string {
	return c.Text
}

func (c *TextComment) Attributes() map[string]string {
	return map[string]string{"color": c.Color}
}

/*
Picture Comment
*/

// TODO
