package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTextComment(t *testing.T) {
	tc := NewTextComment("hello world", "red")
	assert.Equal(t, "hello world", tc.Content())
	assert.Equal(t, map[string]string{"color": "red"}, tc.Attributes())
	assert.Equal(t, "text", tc.Type())
}

func TestNewTextCommentFromMap(t *testing.T) {
	tc := NewTextComment("hello world", "red")
	tc2, ok := NewTextCommentFromMap(map[string]string{
		"text":  "hello world",
		"color": "red",
	})
	if !ok {
		assert.Fail(t, "fail to new")
	}
	assert.Equal(t, tc, tc2)
}
