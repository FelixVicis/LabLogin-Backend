package main

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"time"
)

/*
filename.go by Allen J. Mills
    mm.d.yy

    Description
*/

type StorageInfo struct {
	Domain string
	ID     interface{}
}

// Students
//------------------------

type Student struct { // Key: UUID
	First, Last, UUID string
	MostRecent        int64
}

func (s *Student) Key(ctx context.Context, k interface{}) *datastore.Key {
	si := k.(StorageInfo)
	return datastore.NewKey(ctx, si.Domain+"-"+StudentTable, si.ID.(string), 0, nil)
}

type LoginRecord struct { // Key: in.string - uuid
	UUID    string
	In, Out time.Time
}

func (l *LoginRecord) Key(ctx context.Context, k interface{}) *datastore.Key {
	si := k.(StorageInfo)
	return datastore.NewKey(ctx, si.Domain+"-"+LoginRecordTable, "", si.ID.(int64), nil)
}

////////////////////////////////////
//Function: NewRecord
///~
func NewRecord(uuid string) *LoginRecord {
	r := &LoginRecord{}
	r.UUID = uuid
	r.In = time.Now()
	return r
}
