package main

import (
	"fmt"
)

/*
filename.go by Allen J. Mills
    mm.d.yy

    Description
*/

type StorageKey struct {
	Domain      string
	LoginDomain string
	ID          interface{}
}

func (s StorageKey) ToString() string {
	return fmt.Sprint(s.Domain, s.LoginDomain, `-`, s.ID)
}
