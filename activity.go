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

type LabelComment struct {
	Id      int
	Status  int
	Comment Comment
}

// create new activity with customized name
func NewActivity(name string) *Activity {
	return &Activity{
		Name:          name,
		Dict:          make(map[int]*LabelComment),
		InitialQueue:  make([]*LabelComment, 0, QueueDefaultLength),
		ApprovedQueue: make([]*LabelComment, 0, QueueDefaultLength),
	}
}

// Activity struct
type Activity struct {
	mutex          sync.Mutex
	Name           string
	Dict           map[int]*LabelComment
	InitialQueue   []*LabelComment
	ApprovedQueue  []*LabelComment
	TotalCount     int
	ApprovedCount  int
	DeniedCount    int
	DisplayedCount int
}

// add a comment, initialized with an unique id and Initial status
func (act *Activity) Add(c Comment) *LabelComment {
	act.mutex.Lock()
	defer act.mutex.Unlock()

	act.TotalCount++
	id := act.TotalCount
	lc := &LabelComment{Id: id, Comment: c, Status: CommentStatusInitial}
	act.Dict[id] = lc
	act.InitialQueue = append(act.InitialQueue, lc)
	return lc
}

// get comments with Initial status for reviewing, and then change their status to Pending
func (act *Activity) Review() (r []*LabelComment) {
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
func (act *Activity) Approve(lcs []*LabelComment) {
	act.mutex.Lock()
	defer act.mutex.Unlock()

	act.ApprovedQueue = append(act.ApprovedQueue, lcs...)

	for _, c := range lcs {
		c.Status = CommentStatusApproved
	}
	act.ApprovedCount += int(len(lcs))
}

// deny comments, that is, change their status to Denied
func (act *Activity) Deny(lcs []*LabelComment) {
	act.mutex.Lock()
	defer act.mutex.Unlock()

	for _, c := range lcs {
		c.Status = CommentStatusDenied
	}
	act.DeniedCount += len(lcs)
}

// get comments with Approved status for displaying, then change their status to Displayed
func (act *Activity) Display() (r []*LabelComment) {
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
func (act *Activity) Fetch(ids []int) (r []*LabelComment) {
	act.mutex.Lock()
	defer act.mutex.Unlock()

	r = make([]*LabelComment, 0, len(ids))
	for _, id := range ids {
		if d, ok := act.Dict[id]; ok {
			r = append(r, d)
		}
	}
	return
}

// reset the activity
func (act *Activity) Reset() {
	act.mutex.Lock()
	defer act.mutex.Unlock()

	act.TotalCount = 0
	act.ApprovedCount = 0
	act.DeniedCount = 0
	act.DisplayedCount = 0

	act.Dict = make(map[int]*LabelComment)
	act.InitialQueue = make([]*LabelComment, 0, QueueDefaultLength)
	act.ApprovedQueue = make([]*LabelComment, 0, QueueDefaultLength)
}
