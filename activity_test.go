package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestBasicActivity(t *testing.T) {
	act := &BasicActivity{
		CommentMap:    make(map[int]*LabelComment),
		InitialQueue:  make([]*LabelComment, 0, QueueDefaultLength),
		ApprovedQueue: make([]*LabelComment, 0, QueueDefaultLength),
	}

	bc1 := NewTextComment("content", "red")
	bc2 := NewTextComment("content", "green")
	c1 := act.Add(bc1)

	assert.Equal(t, bc1.Content(), c1.Content)
	assert.Equal(t, bc1.Attributes(), c1.Attributes)
	assert.Equal(t, "text", c1.Type)
	assert.Equal(t, CommentStatusInitial, c1.Status)
	assert.Equal(t, act.TotalCount, c1.Id)
	assert.Equal(t, 1, act.TotalCount)
	assert.Equal(t, 0, act.ApprovedCount)
	assert.Equal(t, 0, act.DeniedCount)
	assert.Equal(t, 0, act.DisplayedCount)

	c2 := act.Add(bc2)
	assert.Equal(t, act.TotalCount, c2.Id)
	assert.Equal(t, 2, act.TotalCount)

	rcs := act.Review()
	assert.Equal(t, 2, len(rcs))
	assert.Equal(t, 2, act.TotalCount)
	assert.Equal(t, 0, act.ApprovedCount)
	assert.Equal(t, 0, act.DeniedCount)
	assert.Equal(t, 0, act.DisplayedCount)

	for _, rc := range rcs {
		assert.Equal(t, "content", rc.Content)
	}

	act.Approve(rcs)
	assert.Equal(t, 2, act.TotalCount)
	assert.Equal(t, 2, act.ApprovedCount)
	assert.Equal(t, 0, act.DeniedCount)
	assert.Equal(t, 0, act.DisplayedCount)

	c3 := act.Add(NewTextComment("content2", "red"))
	c4 := act.Add(NewTextComment("content2", "red"))
	act.Deny([]*LabelComment{c3, c4})

	assert.Equal(t, 4, act.TotalCount)
	assert.Equal(t, 2, act.ApprovedCount)
	assert.Equal(t, 2, act.DeniedCount)
	assert.Equal(t, 0, act.DisplayedCount)

	dcs := act.Display()
	assert.Equal(t, 2, len(dcs))
	assert.Equal(t, 4, act.TotalCount)
	assert.Equal(t, 2, act.ApprovedCount)
	assert.Equal(t, 2, act.DeniedCount)
	assert.Equal(t, 2, act.DisplayedCount)

	for _, rc := range dcs {
		assert.Equal(t, "content", rc.Content)
	}

	act.Reset()
	assert.Equal(t, 0, act.TotalCount)
	assert.Equal(t, 0, act.ApprovedCount)
	assert.Equal(t, 0, act.DeniedCount)
	assert.Equal(t, 0, act.DisplayedCount)
	assert.Equal(t, 0, len(act.InitialQueue))
	assert.Equal(t, 0, len(act.ApprovedQueue))
}
