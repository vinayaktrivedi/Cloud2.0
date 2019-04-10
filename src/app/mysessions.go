// sessions.go
package main

import (
	"github.com/fenilfadadu/CS628-assn1/userlib"
	"storeit"
	"net/http"
	"github.com/gorilla/sessions"
	"encoding/json"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key = []byte("MTU1NDg5NDg0OHxEdi1CQkFFQ180SUFB")
	store = sessions.NewCookieStore(key)
)

func check_login(w http.ResponseWriter, r *http.Request) (bool,string) {
	session, _ := store.Get(r, "sessionid")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		return false,""
	}
	return true,session.Values["username"].(string)
}

func login(w http.ResponseWriter, r *http.Request) *storeit.User {
	session, _ := store.Get(r, "sessionid")
	username := r.FormValue("username")
    password := r.FormValue("password")
    User,err := storeit.GetUser(username,password)
    if err!=nil{
        
        return nil
    }

    marshalled_user_struct, err := json.Marshal(&User)
    userlib.DatastoreSet(username,marshalled_user_struct)
	session.Values["authenticated"] = true
	session.Values["username"] = username
	session.Save(r, w)
	return User
	
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "sessionid")
	session.Values["authenticated"] = false
	session.Save(r, w)
}
