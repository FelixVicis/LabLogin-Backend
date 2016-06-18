package main

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

/*
filename.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"time"
)

var (
	StudentTable     = "Student"
	LoginRecordTable = "LoginRecord"
)

type Student struct { // UUID
	First, Last, UUID string
}

func (s *Student) Key(ctx context.Context, k interface{}) *datastore.Key {
	return datastore.NewKey(ctx, StudentTable, "", k.(int64), nil)
}

type LoginRecord struct { // in.string - uuid
	UUID     string
	In, Out  time.Time
	LoggedIn bool
}

func (l *LoginRecord) Key(ctx context.Context, k interface{}) *datastore.Key {
	return datastore.NewKey(ctx, LoginRecordTable, k.(string), 0, nil)
}

////////////////////////////////////
//Function: NewRecord
///~
func NewRecord(uuid string) *LoginRecord {
	r := &LoginRecord{}
	r.UUID = uuid
	r.In = time.Now()
	r.LoggedIn = true
	return r
}
