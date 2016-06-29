package main

import (
	// "fmt"
	"github.com/Esseh/retrievable"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"net/http"
	// "time"
)

type User struct {
	// Key: UUID
	First, Last, UUID string
}

// Implements: Retrievable
func (u *User) Key(ctx context.Context, k interface{}) *datastore.Key {
	si := k.(StorageKey)
	return datastore.NewKey(ctx, UsersTable, si.ToString(), 0, nil)
}

/////=========================
// Functions
/////

func GetUserInfo(req *http.Request) (User, bool) {
	u := User{}
	ok := true

	if n := req.FormValue("First"); n != "" {
		u.First = n
	} else {
		ok = false
	}

	if n := req.FormValue("Last"); n != "" {
		u.Last = n
	} else {
		ok = false
	}

	if i := req.FormValue("UUID"); i != "" {
		u.UUID = i
	} else {
		ok = false
	}

	return u, ok
}

//////----------------
// Handlers
/////

func CreateTableUser(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Access-Control-Allow-Origin", "*") // Allow for outside access.
	ServeJsonOfStruct(res, JsonOptions{
		Status: "Success",
	}, NewUUID())
}

func CreateUser(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Access-Control-Allow-Origin", "*") // Allow for outside access.

	// Get parameter Domain
	udomain := req.FormValue("User-Domain")
	if udomain == "" {
		ServeJsonOfStruct(res, JsonOptions{
			Status: "Failure",
			Reason: "Missing User-Domain parameter",
			Code:   http.StatusNotAcceptable,
		}, nil)
		return
	}

	uin, ok := GetUserInfo(req)
	if !ok {
		ServeJsonOfStruct(res, JsonOptions{
			Status: "Failure",
			Reason: "Missing User parameters, check documentation.",
			Code:   http.StatusNotAcceptable,
		}, uin)
		return
	}

	ctx := appengine.NewContext(req) // Make Context

	// Check if user already exists.
	retErr := retrievable.GetEntity(ctx, StorageKey{
		LoginDomain: udomain,
		ID:          uin.UUID,
	}, &User{})
	if retErr == nil {
		ServeJsonOfStruct(res, JsonOptions{
			Status: "Failure",
			Reason: "User already registered.",
			Code:   http.StatusForbidden,
		}, nil)
		return
	}

	_, putErr := retrievable.PlaceEntity(ctx, StorageKey{
		LoginDomain: udomain,
		ID:          uin.UUID,
	}, &uin)
	if putErr != nil {
		ServeJsonOfStruct(res, JsonOptions{
			Status: "Failure",
			Reason: "Internal Services Error.",
			Code:   http.StatusInternalServerError,
		}, nil)
		return
	}

	ServeJsonOfStruct(res, JsonOptions{
		Status: "Success",
	}, nil)
}

func DropUser(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Access-Control-Allow-Origin", "*") // Allow for outside access.

	// Get parameter Domain
	udomain := req.FormValue("User-Domain")
	if udomain == "" {
		ServeJsonOfStruct(res, JsonOptions{
			Status: "Failure",
			Reason: "Missing User-Domain parameter",
			Code:   http.StatusNotAcceptable,
		}, nil)
		return
	}

	uin, _ := GetUserInfo(req)
	if uin.UUID == "" {
		ServeJsonOfStruct(res, JsonOptions{
			Status: "Failure",
			Reason: "Missing UUID parameter.",
			Code:   http.StatusNotAcceptable,
		}, nil)
		return
	}

	ctx := appengine.NewContext(req) // Make Context

	delErr := retrievable.DeleteEntity(ctx, uin.Key(ctx, StorageKey{
		LoginDomain: udomain,
		ID:          uin.UUID,
	}))
	// delErr := retrievable.DeleteEntity(ctx, StorageInfo{
	// 	LoginDomain: udomain,
	// 	ID:          uin.UUID,
	// }, &uin)
	if delErr != nil {
		ServeJsonOfStruct(res, JsonOptions{
			Status: "Failure",
			Reason: delErr.Error(),
			Code:   http.StatusInternalServerError,
		}, nil)
		return
	}

	ServeJsonOfStruct(res, JsonOptions{
		Status: "Success",
	}, nil)
}
