package main

/*
filename.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"fmt"
	"github.com/Esseh/retrievable"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"net/http"
	"time"
)

type Record struct {
	// Key(Completed): in.string - uuid
	// Key(Live): uuid
	UUID    string
	In, Out time.Time
}

func MakeRecordTable(si StorageInfo) string {
	return si.Domain + "-" + RecordsTable
}

// Implements: Retrivable
func (l *Record) Key(ctx context.Context, k interface{}) *datastore.Key {
	si := k.(StorageInfo)
	return datastore.NewKey(ctx, MakeRecordTable(si), si.ID.(string), 0, nil)
}

func NewRecord(uuid string) Record {
	r := Record{}
	r.UUID = uuid
	r.In = time.Now()
	return r
}

/////=========================
// Functions
/////

func GetStorageInfo(req *http.Request) (struct {
	Domain      string
	LoginDomain string
	UUID        string
}, bool) {
	si := struct {
		Domain      string
		LoginDomain string
		UUID        string
	}{}
	ok := true

	if i := req.FormValue("Domain"); i != "" {
		si.Domain = i
	} else {
		ok = false
	}

	if i := req.FormValue("User-Domain"); i != "" {
		si.LoginDomain = i
	} else {
		ok = false
	}

	if i := req.FormValue("UUID"); i != "" {
		si.UUID = i
	} else {
		ok = false
	}

	return si, ok
}

////////------------------
// Handlers
////

func SelectStateFromRecord(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Access-Control-Allow-Origin", "*") // Allow for outside access.

	// Get Parameters
	sinfo, ok := GetStorageInfo(req)
	if !ok {
		ServeJsonOfStruct(res, JsonOptions{
			Status: "Failure",
			Reason: "Missing Paramters, check documentation.",
			Code:   http.StatusNotAcceptable,
		}, sinfo)
		return
	}

	ctx := appengine.NewContext(req) // Make Context

	// Ensure user exists
	getErr1 := retrievable.GetFromDatastore(ctx, StorageInfo{
		LoginDomain: sinfo.LoginDomain,
		ID:          sinfo.UUID,
	}, &User{})
	if getErr1 != nil {
		ServeJsonOfStruct(res, JsonOptions{
			Status: "Failure",
			Reason: "No such user exists.",
			Code:   http.StatusNotFound,
		}, nil)
		return
	}

	// Get Record
	r := Record{}
	getErr2 := retrievable.GetFromDatastore(ctx, StorageInfo{
		Domain: sinfo.Domain,
		ID:     sinfo.UUID,
	}, &r)
	if getErr2 != nil {
		ServeJsonOfStruct(res, JsonOptions{
			Status: "Success",
		}, "User Not Logged In")
		return
	}

	ServeJsonOfStruct(res, JsonOptions{
		Status: "Success",
	}, "User Logged In")
}

func ToggleStateFromRecord(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Access-Control-Allow-Origin", "*") // Allow for outside access.

	// Get Parameters
	sinfo, ok := GetStorageInfo(req)
	if !ok {
		ServeJsonOfStruct(res, JsonOptions{
			Status: "Failure",
			Reason: "Missing Paramters, check documentation.",
			Code:   http.StatusNotAcceptable,
		}, sinfo)
		return
	}

	ctx := appengine.NewContext(req) // Make Context

	// Ensure user exists
	getErr1 := retrievable.GetFromDatastore(ctx, StorageInfo{
		LoginDomain: sinfo.LoginDomain,
		ID:          sinfo.UUID,
	}, &User{})
	if getErr1 != nil {
		ServeJsonOfStruct(res, JsonOptions{
			Status: "Failure",
			Reason: "No such user exists.",
			Code:   http.StatusNotFound,
		}, nil)
		return
	}

	// Get Record
	r := Record{}
	getErr2 := retrievable.GetFromDatastore(ctx, StorageInfo{
		Domain: sinfo.Domain,
		ID:     sinfo.UUID,
	}, &r)
	if getErr2 != nil { // user is not logged in, lets do that.
		r = NewRecord(sinfo.UUID)
		_, putErr1 := retrievable.PlaceInDatastore(ctx, StorageInfo{
			Domain: sinfo.Domain,
			ID:     sinfo.UUID,
		}, &r)
		if putErr1 != nil {
			ServeJsonOfStruct(res, JsonOptions{
				Status: "Failure",
				Reason: "Internal Services Error (ts1)",
				Code:   http.StatusInternalServerError,
			}, nil)
			return
		}

		ServeJsonOfStruct(res, JsonOptions{
			Status: "Success",
		}, "User is Logged In")
		return
	}
	// User is logged in, lets move them to logout.
	r.Out = time.Now()

	delErr1 := retrievable.DeleteFromDatastore(ctx, StorageInfo{
		Domain: sinfo.Domain,
		ID:     sinfo.UUID,
	}, &r)
	if delErr1 != nil {
		ServeJsonOfStruct(res, JsonOptions{
			Status: "Failure",
			Reason: "Internal Services Error (ts2)",
			Code:   http.StatusInternalServerError,
		}, nil)
		return
	}

	_, putErr2 := retrievable.PlaceInDatastore(ctx, StorageInfo{
		Domain: sinfo.Domain,
		ID:     fmt.Sprint(r.UUID, "-", r.In),
	}, &r)
	if putErr2 != nil {
		ServeJsonOfStruct(res, JsonOptions{
			Status: "Failure",
			Reason: "Internal Services Error (ts3)",
			Code:   http.StatusInternalServerError,
		}, nil)
		return
	}

	ServeJsonOfStruct(res, JsonOptions{
		Status: "Success",
	}, "User is Logged out")
}

/////////===================================
// Queries
///////

func SelectAllFromRecord(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Access-Control-Allow-Origin", "*") // Allow for outside access.

	// Get Parameters
	sinfo, _ := GetStorageInfo(req)
	if sinfo.Domain == "" {
		ServeJsonOfStruct(res, JsonOptions{
			Status: "Failure",
			Reason: "Missing Paramter Domain",
			Code:   http.StatusNotAcceptable,
		}, nil)
		return
	}

	q := datastore.NewQuery(MakeRecordTable(StorageInfo{
		Domain: sinfo.Domain,
	}))
	q = q.Order("In").Order("UUID")

	ctx := appengine.NewContext(req) // Make Context

	recordList := make([]Record, 0)
	for t := q.Run(ctx); ; {
		var x Record
		_, qErr := t.Next(&x)
		if qErr == datastore.Done {
			break
		} else if qErr != nil {
			ServeJsonOfStruct(res, JsonOptions{
				Status: "Failure",
				Reason: qErr.Error(),
				Code:   500,
			}, nil)
			return
		}
		recordList = append(recordList, x)
	}

	ServeJsonOfStruct(res, JsonOptions{
		Status: "Success",
	}, recordList)
}

func SelectCurrentFromRecord(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Access-Control-Allow-Origin", "*") // Allow for outside access.

	// Get Parameters
	sinfo, _ := GetStorageInfo(req)
	if sinfo.Domain == "" {
		ServeJsonOfStruct(res, JsonOptions{
			Status: "Failure",
			Reason: "Missing Paramter Domain",
			Code:   http.StatusNotAcceptable,
		}, nil)
		return
	}

	q := datastore.NewQuery(MakeRecordTable(StorageInfo{
		Domain: sinfo.Domain,
	}))
	q = q.Filter("Out =", time.Time{})
	q = q.Order("In").Order("UUID")

	ctx := appengine.NewContext(req) // Make Context

	recordList := make([]Record, 0)
	for t := q.Run(ctx); ; {
		var x Record
		_, qErr := t.Next(&x)
		if qErr == datastore.Done {
			break
		} else if qErr != nil {
			ServeJsonOfStruct(res, JsonOptions{
				Status: "Failure",
				Reason: qErr.Error(),
				Code:   500,
			}, nil)
			return
		}
		recordList = append(recordList, x)
	}

	ServeJsonOfStruct(res, JsonOptions{
		Status: "Success",
	}, recordList)
}

func DropAllFromRecord(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Access-Control-Allow-Origin", "*") // Allow for outside access.

	// Get Parameters
	sinfo, _ := GetStorageInfo(req)
	if sinfo.Domain == "" {
		ServeJsonOfStruct(res, JsonOptions{
			Status: "Failure",
			Reason: "Missing Paramter Domain",
			Code:   http.StatusNotAcceptable,
		}, nil)
		return
	}

	q := datastore.NewQuery(MakeRecordTable(StorageInfo{
		Domain: sinfo.Domain,
	})).KeysOnly()

	ctx := appengine.NewContext(req) // Make Context

	recordKeys, qErr := q.GetAll(ctx, nil)
	if qErr != nil {
		ServeJsonOfStruct(res, JsonOptions{
			Status: "Failure",
			Reason: qErr.Error(),
			Code:   500,
		}, nil)
		return
	}

	delErr := datastore.DeleteMulti(ctx, recordKeys)
	if delErr != nil {
		ServeJsonOfStruct(res, JsonOptions{
			Status: "Failure",
			Reason: delErr.Error(),
			Code:   500,
		}, nil)
		return
	}

	ServeJsonOfStruct(res, JsonOptions{
		Status: "Success",
	}, nil)
}
