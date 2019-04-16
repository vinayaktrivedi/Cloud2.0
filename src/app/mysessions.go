// sessions.go
package main

import (
	"userlib"
	"storeit"
	"net/http"
	"github.com/gorilla/sessions"
	"encoding/json"
	"fmt"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key = []byte("MTU1NDg5NDg0OHxEdi1CQkFFQ180SUFB")
	store = sessions.NewCookieStore(key)
)

func check_login(w http.ResponseWriter, r *http.Request) (bool,string) {
	session, err := store.Get(r, "sessionid")
	if err!=nil{
		fmt.Println(err)
	}
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		return true,"vinayakt"
	}
	return true,session.Values["username"].(string)
}

func login(w http.ResponseWriter, r *http.Request) *storeit.User {
	
	username := r.FormValue("username")
    password := r.FormValue("password")
    if _,err:=userlib.DatastoreGet(username); err==false {	//check if registered!
    	return nil
    }

    session, _ := store.Get(r, "sessionid")
    User,err := storeit.GetUser(username,password)
    if err!=nil{
        
        return nil
    }

    vara, err := json.Marshal(&User)
    //userlib.DatastoreSet(username,marshalled_user_struct)
	session.Values["authenticated"] = true
	session.Values["username"] = username
	session.Values["data"] = vara
	session.Save(r, w)
	return User
	
}

func getUser(r *http.Request) []byte {
	session, _ := store.Get(r, "sessionid")
	return session.Values["data"].([]byte)
}

func logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("yes")
	session, _ := store.Get(r, "sessionid")
	session.Values["authenticated"] = false
	session.Save(r, w)
}
