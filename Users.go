package main
import (
	"github.com/julienschmidt/httprouter"
	"github.com/Esseh/retrievable"
	"google.golang.org/appengine"	
	"net/http"
	"fmt"
	"time"
)

func UserState(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)	// Make Context
	res.Header().Set("Access-Control-Allow-Origin", "*") // Allow for outside access.
	domain := req.FormValue("Domain")	// Get Domain
	uuid := req.FormValue("UUID")		// Get UUID
	// Retrieve Student
	s := Student{}	
	retErr := retrievable.GetEntity(ctx, &s, StorageInfo{
		Domain:domain,
		ID:uuid,
	})
	
	// Output result
	if retErr != nil {
		fmt.Fprint(res,`{"result":"failure"}`)
		return
	}
	fmt.Fprint(res,`{"result":"success"}`)
}

func RegisterStudent(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)	// Make Context	
	res.Header().Set("Access-Control-Allow-Origin", "*") // Allow for outside access.
	
	// Get parameter First
	first := req.FormValue("First")
	if first == "" {
		fmt.Fprint(res,`{"result":"failure","reason":"Missing Parameter First"}`)
		return
	}
	// Get parameter Last
	last  := req.FormValue("Last")
	if last == "" {
		fmt.Fprint(res,`{"result":"failure","reason":"Missing Parameter Last"}`)
		return
	}
	// Get parameter UUID
	uuid  := req.FormValue("UUID")
	if uuid == "" {
		fmt.Fprint(res,`{"result":"failure","reason":"Missing Parameter UUID"}`)
		return
	}
	// Get parameter Domain
	domain:= req.FormValue("Domain")
	if domain == "" {
		fmt.Fprint(res,`{"result":"failure","reason":"Missing Parameter Domain"}`)
		return
	}	
	// Check if student already exists.
	s := Student{}	
	retErr := retrievable.GetEntity(ctx, &s, StorageInfo{
		Domain:domain,
		ID:uuid,
	})
	if retErr == nil {
		fmt.Fprint(res,`{"result":"failure","reason":"Student Already Exists"}`)
		return
	}
	
	// Make student
	s.First = first
	s.Last = last
	s.UUID = uuid
	
	_ , plaErr := retrievable.PlaceEntity(ctx,StorageInfo{
		Domain:domain,
		ID:uuid,
	},&s)
	if plaErr != nil {
		fmt.Fprint(res,`{"result":"failure","reason":"Internal Server Error"}`)	
		return
	}
	
	// Success
	fmt.Fprint(res,`{"result":"success"}`)		
	
}


func LoginUser(res http.ResponseWriter, req *http.Request, params httprouter.Params){
	ctx := appengine.NewContext(req)	// Make Context
	// Get parameter UUID
	uuid  := req.FormValue("UUID")
	if uuid == "" {
		fmt.Fprint(res,`{"result":"failure","reason":"Missing Parameter UUID"}`)
		return
	}
	// Get parameter Domain
	domain:= req.FormValue("Domain")
	if domain == "" {
		fmt.Fprint(res,`{"result":"failure","reason":"Missing Parameter Domain"}`)
		return
	}
	// Get Student
	s := Student{}	
	retErr := retrievable.GetEntity(ctx, &s, StorageInfo{
		Domain:domain,
		ID:uuid,
	})
	if retErr != nil {
		fmt.Fprint(res,`{"result":"failure","reason":"user does not exist"}`)	
		return
	}	
	// Make sure the user is not logged in.
	l := LoginRecord{}
	retErr2 := retrievable.GetEntity(ctx, &l, StorageInfo{
		Domain:domain,
		ID:s.MostRecent,
	})
	if retErr2 == nil {
		fmt.Fprint(res,`{"result":"failure","reason":"user already logged in"}`)
		return
	}
	
	// Log the user in.
	// Initialize LoginRecord
	l.In = time.Now()
	l.UUID = uuid
	// Place LoginRecord
	key , plaErr := retrievable.PlaceEntity(ctx,StorageInfo{
		Domain:domain,
		ID:0,
	},&l)
	if plaErr != nil {
		fmt.Fprint(res,`{"result":"failure","reason":"internal server error (1)"}`)
		return
	}
	
	// Make Record ID the user's most recent login record
	s.MostRecent = key.IntID()
	_ , plaErr2 := retrievable.PlaceEntity(ctx,StorageInfo{
		Domain:domain,
		ID:uuid,	
	},&s)
	if plaErr2 != nil {
		fmt.Fprint(res,`{"result":"failure","reason":"internal server error (2)"}`)
		return
	}
	fmt.Fprint(res,`{"result":"success"}`)
}
func LogoutUser(res http.ResponseWriter, req *http.Request, params httprouter.Params){
	ctx := appengine.NewContext(req)	// Make Context
	// Get parameter UUID
	uuid  := req.FormValue("UUID")
	if uuid == "" {
		fmt.Fprint(res,`{"result":"failure","reason":"Missing Parameter UUID"}`)
	}
	// Get parameter Domain
	domain:= req.FormValue("Domain")
	if domain == "" {
		fmt.Fprint(res,`{"result":"failure","reason":"Missing Parameter Domain"}`)
	}
	// Get Student
	s := Student{}	
	retErr := retrievable.GetEntity(ctx, &s, StorageInfo{
		Domain:domain,
		ID:uuid,
	})
	if retErr != nil {
		fmt.Fprint(res,`{"result":"failure","reason":"user does not exist"}`)	
	}	
	// Make sure the user is actually logged in.
	l := LoginRecord{}
	retErr2 := retrievable.GetEntity(ctx, &l, StorageInfo{
		Domain:domain,
		ID:s.MostRecent,
	})
	if retErr2 != nil {
		fmt.Fprint(res,`{"result":"failure","reason":"user is not logged in"}`)
		return
	}
	
	// Log the user out.
	l.Out = time.Now()
	_ , plaErr := retrievable.PlaceEntity(ctx,StorageInfo{
		Domain:domain,
		ID:s.MostRecent,
	},&l)
	if plaErr != nil {
		fmt.Fprint(res,`{"result":"failure","reason":"internal server error (1)"}`)
		return
	}
	
	s.MostRecent = 0
	_ , plaErr2 := retrievable.PlaceEntity(ctx,StorageInfo{
		Domain:domain,
		ID:uuid,	
	},&s)
	if plaErr2 != nil {
		fmt.Fprint(res,`{"result":"failure","reason":"internal server error (2)"}`)
		return
	}
	fmt.Fprint(res,`{"result":"success"}`)
}