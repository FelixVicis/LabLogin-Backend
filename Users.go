package main

import (
	"fmt"
	"github.com/Esseh/retrievable"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine"
	"net/http"
	"time"
)

func UserState(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)                     // Make Context
	res.Header().Set("Access-Control-Allow-Origin", "*") // Allow for outside access.
	domain := req.FormValue("Domain")                    // Get Domain
	uuid := req.FormValue("UUID")                        // Get UUID
	// Retrieve Student
	s := Student{}
	retErr := retrievable.GetFromDatastore(ctx, StorageInfo{
		Domain: domain,
		ID:     uuid,
	}, &s)
	// If student doesn't exist...
	if retErr != nil {
		fmt.Fprint(res, `{"result":"no user"}`)
		return
	}
	// Check if student is logged in.
	l := LoginRecord{}
	retErr2 := retrievable.GetFromDatastore(ctx, StorageInfo{
		Domain: domain,
		ID:     s.MostRecent,
	}, &l)
	if retErr2 == nil && l.Out.IsZero() {
		fmt.Fprint(res, `{"result":"logged in"}`)
		return
	}
	fmt.Fprint(res, `{"result":"logged off"}`)
	// if retErr2 != nil {
	// 	fmt.Fprint(res, `{"result":"logged off"}`)
	// 	return
	// }
}

func RegisterStudent(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)                     // Make Context
	res.Header().Set("Access-Control-Allow-Origin", "*") // Allow for outside access.

	// Get parameter First
	first := req.FormValue("First")
	if first == "" {
		fmt.Fprint(res, `{"result":"failure","reason":"Missing Parameter First"}`)
		return
	}
	// Get parameter Last
	last := req.FormValue("Last")
	if last == "" {
		fmt.Fprint(res, `{"result":"failure","reason":"Missing Parameter Last"}`)
		return
	}
	// Get parameter UUID
	uuid := req.FormValue("UUID")
	if uuid == "" {
		fmt.Fprint(res, `{"result":"failure","reason":"Missing Parameter UUID"}`)
		return
	}
	// Get parameter Domain
	domain := req.FormValue("Domain")
	if domain == "" {
		fmt.Fprint(res, `{"result":"failure","reason":"Missing Parameter Domain"}`)
		return
	}
	// Check if student already exists.
	s := Student{}
	retErr := retrievable.GetFromDatastore(ctx, StorageInfo{
		Domain: domain,
		ID:     uuid,
	}, &s)
	if retErr == nil {
		fmt.Fprint(res, `{"result":"failure","reason":"Student Already Exists"}`)
		return
	}

	// Make student
	s.First = first
	s.Last = last
	s.UUID = uuid

	_, plaErr := retrievable.PlaceInDatastore(ctx, StorageInfo{
		Domain: domain,
		ID:     uuid,
	}, &s)
	if plaErr != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"Internal Server Error"}`)
		return
	}

	// Success
	fmt.Fprint(res, `{"result":"success"}`)

}

func LoginUser(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Access-Control-Allow-Origin", "*") // Allow for outside access.
	ctx := appengine.NewContext(req)                     // Make Context
	// Get parameter UUID
	uuid := req.FormValue("UUID")
	if uuid == "" {
		fmt.Fprint(res, `{"result":"failure","reason":"Missing Parameter UUID"}`)
		return
	}
	// Get parameter Domain
	domain := req.FormValue("Domain")
	if domain == "" {
		fmt.Fprint(res, `{"result":"failure","reason":"Missing Parameter Domain"}`)
		return
	}
	// Get Student
	s := Student{}
	retErr := retrievable.GetFromDatastore(ctx, StorageInfo{
		Domain: domain,
		ID:     uuid,
	}, &s)
	if retErr != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"user does not exist"}`)
		return
	}
	// Make sure the user is not logged in.
	l := LoginRecord{}
	retErr2 := retrievable.GetFromDatastore(ctx, StorageInfo{
		Domain: domain,
		ID:     s.MostRecent,
	}, &l)
	if retErr2 == nil {
		fmt.Fprint(res, `{"result":"failure","reason":"user already logged in"}`)
		return
	}

	// Log the user in.
	// Initialize LoginRecord
	l.In = time.Now()
	l.UUID = uuid
	// Place LoginRecord
	key, plaErr := retrievable.PlaceInDatastore(ctx, StorageInfo{
		Domain: domain,
		ID:     int64(0),
	}, &l)
	if plaErr != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"internal server error (1)"}`)
		return
	}

	// Make Record ID the user's most recent login record
	s.MostRecent = key.IntID()
	_, plaErr2 := retrievable.PlaceInDatastore(ctx, StorageInfo{
		Domain: domain,
		ID:     uuid,
	}, &s)
	if plaErr2 != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"internal server error (2)"}`)
		return
	}
	fmt.Fprint(res, `{"result":"success"}`)
}
func LogoutUser(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Access-Control-Allow-Origin", "*") // Allow for outside access.
	ctx := appengine.NewContext(req)                     // Make Context
	// Get parameter UUID
	uuid := req.FormValue("UUID")
	if uuid == "" {
		fmt.Fprint(res, `{"result":"failure","reason":"Missing Parameter UUID"}`)
	}
	// Get parameter Domain
	domain := req.FormValue("Domain")
	if domain == "" {
		fmt.Fprint(res, `{"result":"failure","reason":"Missing Parameter Domain"}`)
	}
	// Get Student
	s := Student{}
	retErr := retrievable.GetFromDatastore(ctx, StorageInfo{
		Domain: domain,
		ID:     uuid,
	}, &s)
	if retErr != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"user does not exist"}`)
	}
	// Make sure the user is actually logged in.
	l := LoginRecord{}
	retErr2 := retrievable.GetFromDatastore(ctx, StorageInfo{
		Domain: domain,
		ID:     s.MostRecent,
	}, &l)
	if retErr2 != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"user is not logged in"}`)
		return
	}

	// Log the user out.
	l.Out = time.Now()
	_, plaErr := retrievable.PlaceInDatastore(ctx, StorageInfo{
		Domain: domain,
		ID:     s.MostRecent,
	}, &l)
	if plaErr != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"internal server error (1)"}`)
		return
	}

	s.MostRecent = 0
	_, plaErr2 := retrievable.PlaceInDatastore(ctx, StorageInfo{
		Domain: domain,
		ID:     uuid,
	}, &s)
	if plaErr2 != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"internal server error (2)"}`)
		return
	}
	fmt.Fprint(res, `{"result":"success"}`)
}
