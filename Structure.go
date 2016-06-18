package main

/*
filename.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"time"
)

type Client struct { // ID
	ID string
}

type Student struct { // UUID
	First, Last, UUID string
}

type LoginRecord struct { // in.string - uuid
	UUID     string
	In, Out  time.Time
	LoggedIn bool
}

func NewRecord(uuid string) *LoginRecord {
	r := &LoginRecord{}
	r.UUID = uuid
	r.In = time.Now()
	r.LoggedIn = true
	return r
}
