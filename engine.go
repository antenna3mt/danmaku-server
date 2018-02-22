// Copyright 2018 Yi Jin. All rights reserved.
// license that can be found in the LICENSE file.

package main

import (
	"sync"
	"fmt"
)

const (
	ActivityTokenLength = 8
	AdminTokenLength    = 16
)

var (
	NotAuthorizedError = fmt.Errorf("not authorized")
	NotExistError      = fmt.Errorf("not exist")
	IllFormatError     = fmt.Errorf("ill format")
)

// activity extend BasicActivity
type Activity struct {
	BasicActivity
	Id           int
	Name         string
	CommentToken string
	ReviewToken  string
	DisplayToken string
	ReviewOn     bool
}

func NewEngine() *Engine {
	return &Engine{
		AdminToken:  NewAuthToken(AdminTokenLength),
		ActivityMap: make(map[int]*Activity),
		TokenMap:    make(map[string]*Activity),
	}
}

// Engine struct
type Engine struct {
	AdminToken  string
	ActivityMap map[int]*Activity
	TokenMap    map[string]*Activity
	IdCount     int
	mutex       sync.Mutex
}

// generate a unique token
func (e *Engine) newToken() string {
	for {
		token := NewAuthToken(ActivityTokenLength)
		if _, ok := e.TokenMap[token]; !ok {
			return token
		}
	}
}

// login
func (e *Engine) Login(auth_token string) (string, error) {
	if auth_token == e.AdminToken {
		return "admin", nil
	}

	act, ok := e.ActivityByToken(auth_token)
	if ok {
		switch auth_token {
		case act.CommentToken:
			return "comment", nil
		case act.ReviewToken:
			return "review", nil
		case act.DisplayToken:
			return "display", nil
		}
	}
	return "", NotAuthorizedError
}

// create a activity with name and add it to engine
func (e *Engine) NewActivity(auth_token string, name string) (*Activity, error) {
	if !IsOneOf(auth_token, e.AdminToken) {
		return nil, NotAuthorizedError
	}

	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.IdCount++
	id := e.IdCount
	act := &Activity{
		BasicActivity: BasicActivity{
			CommentMap:    make(map[int]*LabelComment),
			InitialQueue:  make([]*LabelComment, 0, QueueDefaultLength),
			ApprovedQueue: make([]*LabelComment, 0, QueueDefaultLength),
		},
		Id:           id,
		Name:         name,
		CommentToken: e.newToken(),
		ReviewToken:  e.newToken(),
		DisplayToken: e.newToken(),
		ReviewOn:     true,
	}

	e.ActivityMap[id] = act
	e.TokenMap[act.CommentToken] = act
	e.TokenMap[act.ReviewToken] = act
	e.TokenMap[act.DisplayToken] = act
	return act, nil
}

// get activity by token
func (e *Engine) ActivityByToken(token string) (*Activity, bool) {
	act, ok := e.TokenMap[token]
	return act, ok
}

// all activity; action permit: admin
func (e *Engine) Activities(auth_token string) ([]*Activity, error) {
	if !IsOneOf(auth_token, e.AdminToken) {
		return nil, NotAuthorizedError
	}

	r := make([]*Activity, 0, len(e.ActivityMap))
	for _, a := range e.ActivityMap {
		r = append(r, a)
	}
	return r, nil
}

// delete activity by id; action permit: admin
func (e *Engine) DelActivity(auth_token string, id int) (error) {
	if !IsOneOf(auth_token, e.AdminToken) {
		return NotAuthorizedError
	}

	act, ok := e.ActivityMap[id]
	if !ok {
		return NotExistError
	}
	delete(e.TokenMap, act.CommentToken)
	delete(e.TokenMap, act.ReviewToken)
	delete(e.TokenMap, act.DisplayToken)
	delete(e.ActivityMap, id)
	return nil
}

// rename activity by id; action permit: admin
func (e *Engine) RenameActivity(auth_token string, id int, name string) (error) {
	if !IsOneOf(auth_token, e.AdminToken) {
		return NotAuthorizedError
	}

	act, ok := e.ActivityMap[id]
	if !ok {
		return NotExistError
	}
	act.Name = name
	return nil
}

// turn review on; action permit: admin
func (e *Engine) ReviewOn(auth_token string, id int) (error) {
	if !IsOneOf(auth_token, e.AdminToken) {
		return NotAuthorizedError
	}

	act, ok := e.ActivityMap[id]
	if !ok {
		return NotExistError
	}
	act.ReviewOn = true
	return nil
}

// turn review Off; action permit: admin
func (e *Engine) ReviewOff(auth_token string, id int) (error) {
	if !IsOneOf(auth_token, e.AdminToken) {
		return NotAuthorizedError
	}

	act, ok := e.ActivityMap[id]
	if !ok {
		return NotExistError
	}
	act.ReviewOn = false
	return nil
}

// push a comment; action permit: comment, review, display
func (e *Engine) Push(auth_token string, tp string, attr map[string]string) (*LabelComment, error) {
	act, ok := e.ActivityByToken(auth_token)
	if !ok {
		return nil, NotExistError
	}

	if !IsOneOf(auth_token, act.CommentToken, act.ReviewToken, act.DisplayToken) {
		return nil, NotAuthorizedError
	}

	c, ok := NewComment(tp, attr)
	if !ok {
		return nil, IllFormatError
	}

	lc := act.Add(c)
	act.Approve(act.Review())

	return lc, nil
}

// review; action permit: review
func (e *Engine) Review(auth_token string) ([]*LabelComment, error) {
	act, ok := e.ActivityByToken(auth_token)
	if !ok {
		return nil, NotExistError
	}

	if !IsOneOf(auth_token, act.ReviewToken) {
		return nil, NotAuthorizedError
	}

	return act.Review(), nil
}

// approve; action permit: review
func (e *Engine) Approve(auth_token string, ids []int) (error) {
	act, ok := e.ActivityByToken(auth_token)
	if !ok {
		return NotExistError
	}

	if !IsOneOf(auth_token, act.ReviewToken) {
		return NotAuthorizedError
	}

	lcs := act.Fetch(ids)
	act.Approve(lcs)
	return nil
}

// deny; action permit: review
func (e *Engine) Deny(auth_token string, ids []int) (error) {
	act, ok := e.ActivityByToken(auth_token)
	if !ok {
		return NotExistError
	}

	if !IsOneOf(auth_token, act.ReviewToken) {
		return NotAuthorizedError
	}

	lcs := act.Fetch(ids)
	act.Deny(lcs)
	return nil
}

// display; action permit: display
func (e *Engine) Display(auth_token string) ([]*LabelComment, error) {
	act, ok := e.ActivityByToken(auth_token)
	if !ok {
		return nil, NotExistError
	}

	if !IsOneOf(auth_token, act.DisplayToken) {
		return nil, NotAuthorizedError
	}

	return act.Display(), nil
}

// reset; action permit: admin
func (e *Engine) Reset(auth_token string) (error) {
	act, ok := e.ActivityByToken(auth_token)
	if !ok {
		return NotExistError
	}

	if !IsOneOf(auth_token, e.AdminToken) {
		return NotAuthorizedError
	}
	act.Reset()
	return nil
}
