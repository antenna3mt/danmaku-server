package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TokenMatch(e *Engine, act *Activity, token string) bool {
	if a, ok := e.ActivityByToken(token); ok {
		return act == a
	} else {
		return false
	}
}

func TestEngine_Admin(t *testing.T) {
	e := NewEngine()

	act1, err := e.NewActivity(e.AdminToken, "First")
	assert.Nil(t, err)
	assert.Equal(t, e.IdCount, act1.Id)
	act2, err := e.NewActivity(e.AdminToken, "Second")
	assert.Nil(t, err)
	assert.Equal(t, e.IdCount, act2.Id)

	assert.True(t, TokenMatch(e, act1, act1.DisplayToken))
	assert.True(t, TokenMatch(e, act1, act1.ReviewToken))
	assert.True(t, TokenMatch(e, act1, act1.CommentToken))
	assert.True(t, TokenMatch(e, act2, act2.DisplayToken))
	assert.True(t, TokenMatch(e, act2, act2.ReviewToken))
	assert.True(t, TokenMatch(e, act2, act2.CommentToken))

	_, err = e.Activities("")
	assert.Error(t, err)
	acts, err := e.Activities(e.AdminToken)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(acts))

	e.DelActivity(e.AdminToken, act1.Id)
	acts, _ = e.Activities(e.AdminToken)
	assert.Equal(t, 1, len(acts))

	err = e.RenameActivity(e.AdminToken, act1.Id, "Hello")
	assert.Error(t, err)

	err = e.RenameActivity(e.AdminToken, act2.Id, "Hello")
	assert.Nil(t, err)
	assert.Equal(t, "Hello", act2.Name)
}

func TestEngine_Activity(t *testing.T) {
	e := NewEngine()
	act, err := e.NewActivity(e.AdminToken, "Hello")
	assert.Nil(t, err)

	lc, err := e.Push(act.CommentToken, "text", map[string]string{"text": "Hello", "color": "red"})

	assert.Nil(t, err)
	assert.Equal(t, "text", lc.Type)
	assert.Equal(t, "Hello", lc.Content)
	assert.Equal(t, map[string]string{"color": "red"}, lc.Attributes)
	assert.Equal(t, CommentStatusInitial, lc.Status)
	assert.Equal(t, act.TotalCount, lc.Id)

	rcs, err := e.Review(act.ReviewToken)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(rcs))

	err = e.Approve(act.ReviewToken, []int{rcs[0].Id})
	assert.Nil(t, err)

	dcs, err := e.Display(act.DisplayToken)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(dcs))
}
