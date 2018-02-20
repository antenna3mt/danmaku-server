// Copyright 2018 Yi Jin. All rights reserved.
// license that can be found in the LICENSE file.

package main

import "sync"

const (
	CommentStatusInitial   int = iota
	CommentStatusPending
	CommentStatusApproved
	CommentStatusDenied
	CommentStatusDisplayed
)

const (
	QueueDefaultLength = 1000
)

// comment with labelled id and status
type LabelComment struct {
	Id         int
	Status     int
	Type       string
	Content    string
	Attributes map[string]string
}

// BasicActivity struct
type BasicActivity struct {
	mutex          sync.Mutex
	CommentMap     map[int]*LabelComment
	InitialQueue   []*LabelComment
	ApprovedQueue  []*LabelComment
	TotalCount     int
	ApprovedCount  int
	DeniedCount    int
	DisplayedCount int
}

// add a comment, initialized with an unique id and Initial status
func (act *BasicActivity) Add(c Comment) *LabelComment {
	act.mutex.Lock()
	defer act.mutex.Unlock()

	act.TotalCount++
	id := act.TotalCount
	lc := &LabelComment{Id: id, Type: c.Type(), Content: c.Content(), Attributes: c.Attributes(), Status: CommentStatusInitial}
	act.CommentMap[id] = lc
	act.InitialQueue = append(act.InitialQueue, lc)
	return lc
}

// get comments with Initial status for reviewing, and then change their status to Pending
func (act *BasicActivity) Review() (r []*LabelComment) {
	act.mutex.Lock()
	defer act.mutex.Unlock()

	r = act.InitialQueue
	act.InitialQueue = make([]*LabelComment, 0, QueueDefaultLength)
	for _, d := range r {
		d.Status = CommentStatusPending
	}
	return
}

// approve comments, that is, change their status to Approved
func (act *BasicActivity) Approve(lcs []*LabelComment) {
	act.mutex.Lock()
	defer act.mutex.Unlock()

	act.ApprovedQueue = append(act.ApprovedQueue, lcs...)

	for _, c := range lcs {
		c.Status = CommentStatusApproved
	}
	act.ApprovedCount += int(len(lcs))
}

// deny comments, that is, change their status to Denied
func (act *BasicActivity) Deny(lcs []*LabelComment) {
	act.mutex.Lock()
	defer act.mutex.Unlock()

	for _, c := range lcs {
		c.Status = CommentStatusDenied
	}
	act.DeniedCount += len(lcs)
}

// get comments with Approved status for displaying, then change their status to Displayed
func (act *BasicActivity) Display() (r []*LabelComment) {
	act.mutex.Lock()
	defer act.mutex.Unlock()

	r = act.ApprovedQueue
	act.ApprovedQueue = make([]*LabelComment, 0, QueueDefaultLength)
	for _, c := range r {
		c.Status = CommentStatusDisplayed
	}
	act.DisplayedCount += len(r)
	return
}

// get comments by their ids
func (act *BasicActivity) Fetch(ids []int) (r []*LabelComment) {
	act.mutex.Lock()
	defer act.mutex.Unlock()

	r = make([]*LabelComment, 0, len(ids))
	for _, id := range ids {
		if d, ok := act.CommentMap[id]; ok {
			r = append(r, d)
		}
	}
	return
}

// reset the activity
func (act *BasicActivity) Reset() {
	act.mutex.Lock()
	defer act.mutex.Unlock()

	act.TotalCount = 0
	act.ApprovedCount = 0
	act.DeniedCount = 0
	act.DisplayedCount = 0

	act.CommentMap = make(map[int]*LabelComment)
	act.InitialQueue = make([]*LabelComment, 0, QueueDefaultLength)
	act.ApprovedQueue = make([]*LabelComment, 0, QueueDefaultLength)
}
