// Copyright 2018 Yi Jin. All rights reserved.
// license that can be found in the LICENSE file.

package main

import "fmt"

// flat format for outputting
type FlatComment struct {
	Id         int               `json:"id"`
	Type       string            `json:"type"`
	Content    string            `json:"content"`
	Attributes map[string]string `json:"attributes"`
}

func FlattenComment(c *LabelComment) *FlatComment {
	return &FlatComment{
		Id:         c.Id,
		Type:       c.Type,
		Content:    c.Content,
		Attributes: c.Attributes,
	}
}

type FlatActivity struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	CommentToken   string `json:"comment_token"`
	ReviewToken    string `json:"review_token"`
	DisplayToken   string `json:"display_token"`
	ReviewOn       bool   `json:"review_on"`
	TotalCount     int    `json:"total_count"`
	ApprovedCount  int    `json:"approved_count"`
	DeniedCount    int    `json:"denied_count"`
	DisplayedCount int    `json:"displayed_count"`
}

func FlattenActivity(act *Activity) *FlatActivity {
	return &FlatActivity{
		Id:             act.Id,
		Name:           act.Name,
		CommentToken:   act.CommentToken,
		ReviewToken:    act.ReviewToken,
		DisplayToken:   act.DisplayToken,
		ReviewOn:       act.ReviewOn,
		TotalCount:     act.TotalCount,
		ApprovedCount:  act.ApprovedCount,
		DeniedCount:    act.DeniedCount,
		DisplayedCount: act.DeniedCount,
	}
}

type FlatActivityDigest struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	TotalCount     int    `json:"total_count"`
	ApprovedCount  int    `json:"approved_count"`
	DeniedCount    int    `json:"denied_count"`
	DisplayedCount int    `json:"displayed_count"`
}

func FlattenActivityDigest(act *Activity) *FlatActivityDigest {
	return &FlatActivityDigest{
		Id:             act.Id,
		Name:           act.Name,
		TotalCount:     act.TotalCount,
		ApprovedCount:  act.ApprovedCount,
		DeniedCount:    act.DeniedCount,
		DisplayedCount: act.DeniedCount,
	}
}

/*
Danmaku Service
*/

type Context struct{}

type DanmakuService struct {
	E *Engine
}

func (s *DanmakuService) Login(ctx *Context,
	args *struct {
		Token string
	}, reply *struct {
		Type string `json:"type"`
	}) error {
	tp, err := s.E.Login(args.Token)
	if err != nil {
		return err
	}
	reply.Type = tp
	return nil
}

// new activity
func (s *DanmakuService) NewActivity(ctx *Context,
	args *struct {
		Token string
		Name  string
	}, reply *struct {
		Activity *FlatActivity `json:"activity"`
	}) error {
	act, err := s.E.NewActivity(args.Token, args.Name)
	if err != nil {
		return err
	}
	reply.Activity = FlattenActivity(act)
	return nil
}

// get all activities
func (s *DanmakuService) Activities(ctx *Context, args *struct {
	Token string
}, reply *struct {
	Activities []*FlatActivity `json:"activities"`
}) error {
	acts, err := s.E.Activities(args.Token)
	if err != nil {
		return err
	}
	reply.Activities = make([]*FlatActivity, 0, len(acts))
	for _, act := range acts {
		reply.Activities = append(reply.Activities, FlattenActivity(act))
	}
	return nil
}

// delete activity
func (s *DanmakuService) DelActivity(ctx *Context, args *struct {
	Token string
	Id    int
}, reply *struct{}) error {
	err := s.E.DelActivity(args.Token, args.Id)
	if err != nil {
		return err
	}
	return nil
}

// rename activity
func (s *DanmakuService) RenameActivity(ctx *Context, args *struct {
	Token string
	Id    int
	Name  string
}, reply *struct{}) error {
	err := s.E.RenameActivity(args.Token, args.Id, args.Name)
	if err != nil {
		return err
	}
	return nil
}

// turn review on
func (s *DanmakuService) ReviewOn(ctx *Context, args *struct {
	Token string
	Id    int
}, reply *struct{}) error {
	err := s.E.ReviewOn(args.Token, args.Id)
	if err != nil {
		return err
	}
	return nil
}

// turn review off
func (s *DanmakuService) ReviewOff(ctx *Context, args *struct {
	Token string
	Id    int
}, reply *struct{}) error {
	err := s.E.ReviewOff(args.Token, args.Id)
	if err != nil {
		return err
	}
	return nil
}

// get activity
func (s *DanmakuService) GetActivityDigest(ctx *Context,
	args *struct {
		Token string
	}, reply *struct {
		Activity *FlatActivityDigest `json:"activity"`
	}) error {
	act, ok := s.E.ActivityByToken(args.Token)
	if !ok {
		return fmt.Errorf("not exist")
	}
	reply.Activity = FlattenActivityDigest(act)
	return nil
}

// push a comment
func (s *DanmakuService) Push(ctx *Context,
	args *struct {
		Token string
		Type  string
		Attr  map[string]string
	}, reply *struct {
		Comment *FlatComment `json:"comment"`
	}) error {
	c, err := s.E.Push(args.Token, args.Type, args.Attr)
	if err != nil {
		return err
	}
	reply.Comment = FlattenComment(c)
	return nil
}

// review
func (s *DanmakuService) Review(ctx *Context,
	args *struct {
		Token string
	}, reply *struct {
		Comments []*FlatComment `json:"comments"`
	}) error {
	cs, err := s.E.Review(args.Token)
	if err != nil {
		return err
	}
	reply.Comments = make([]*FlatComment, 0, len(cs))
	for _, c := range cs {
		reply.Comments = append(reply.Comments, FlattenComment(c))
	}
	return nil
}

// approve
func (s *DanmakuService) Approve(ctx *Context,
	args *struct {
		Token string
		Ids   []int
	}, reply *struct{}) error {
	err := s.E.Approve(args.Token, args.Ids)
	if err != nil {
		return err
	}
	return nil
}

// approve
func (s *DanmakuService) Deny(ctx *Context,
	args *struct {
		Token string
		Ids   []int
	}, reply *struct{}) error {
	err := s.E.Deny(args.Token, args.Ids)
	if err != nil {
		return err
	}
	return nil
}

// display
func (s *DanmakuService) Display(ctx *Context,
	args *struct {
		Token string
	}, reply *struct {
		Comments []*FlatComment `json:"comments"`
	}) error {
	cs, err := s.E.Display(args.Token)
	if err != nil {
		return err
	}
	reply.Comments = make([]*FlatComment, 0, len(cs))
	for _, c := range cs {
		reply.Comments = append(reply.Comments, FlattenComment(c))
	}
	return nil
}

// reset
func (s *DanmakuService) Reset(ctx *Context,
	args *struct {
		Token string
	}, reply *struct{}) error {
	err := s.E.Reset(args.Token)
	if err != nil {
		return err
	}
	return nil
}
