// Copyright 2018 Yi Jin. All rights reserved.
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"crypto/rand"
)

func NewAuthToken(len int) string {
	b := make([]byte, len/2+1)
	rand.Read(b[:])
	return fmt.Sprintf("%x", b[:])[:len]
}
