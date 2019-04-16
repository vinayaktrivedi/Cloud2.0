package main

import (
    "go/build"
    "fmt"
    "log"
    "net/http"
    "html/template"
    "regexp"
    "encoding/json"
    "storeit"
    "userlib"
    //"fmt"
    "io/ioutil"
    "strings"
    "mime"
    //"mysessions"

)
var templates *template.Template 
var validPath = regexp.MustCompile("^/(upload|view|download)/([a-zA-Z0-9]+)$")
type MyUser struct {
    Name string
    Username string 
    Image string
    Files []string   //only exported (uppercase) variables can be used in template
}

var global_files map[string][]string

func renderHTML(w http.ResponseWriter, p *MyUser, name string){
    err := templates.ExecuteTemplate(w,name,p)
    if(err!=nil){
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
func registerHandler (w http.ResponseWriter, r *http.Request) {

    if r.Method == http.MethodGet {
        val,_ := check_login(w,r)
        if val == true{
            http.Redirect(w, r, "/view/" , http.StatusFound)
            return
        }else{
            renderHTML(w,nil,"register.html")   
            return
        }
        

    }else if r.Method == http.MethodPost {
        username := r.FormValue("username")
        password := r.FormValue("password")
        _,err := storeit.InitUser(username,password)
        if err!=nil{
            http.Redirect(w, r, "/register/" , http.StatusFound)
            return
        }
        login(w,r)
        var myslice []string
        global_files[username] = myslice
        var html_user MyUser 
        html_user.Name = username
        html_user.Username = username
        html_user.Image = "default"
        html_user.Files = global_files[username]
        renderHTML(w,&html_user,"view.html")
        return
    }else{
        http.Error(w, "Invalid request", http.StatusInternalServerError)
        return
    }
}
func loginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        val,_ := check_login(w,r)
        if val == true{
            http.Redirect(w, r, "/view/" , http.StatusFound)
            return
        }else{
            renderHTML(w,nil,"login.html")    
            return
        }
        
    }else if r.Method == http.MethodPost {

        temp := login(w,r)
        //fmt.Println("temp is ",temp)
        if(temp == nil){
            http.Redirect(w, r, "/" , http.StatusFound)
            return
        }
        var html_user MyUser 
        html_user.Name = temp.Username
        html_user.Username = temp.Username
        html_user.Image = "default"
        html_user.Files = global_files[temp.Username]
        renderHTML(w,&html_user,"view.html")
        return
    }else{
        http.Error(w, "Invalid request", http.StatusInternalServerError)
        return
    }

}
func uploadHandler(w http.ResponseWriter, r *http.Request) {
    val,username := check_login(w,r)
    fmt.Println(val)
    if val == false{
        http.Redirect(w, r, "/" , http.StatusFound)
        return
    }
    if r.Method == http.MethodGet {
        renderHTML(w,nil,"upload.html")
        return
    }else if r.Method == http.MethodPost {

        value := getUser(r)
        var User storeit.User 
        unmarshal_err := json.Unmarshal(value,&User)
        if(unmarshal_err!=nil){
            http.Error(w, "Invalid request", http.StatusInternalServerError)
            return
        }
        r.ParseMultipartForm(10 << 20)
    
        file, handler, err := r.FormFile("files")
        if err != nil {
            http.Redirect(w, r, "/" , http.StatusFound)
            //fmt.Println("Error Retrieving the File")
            //fmt.Println(err)
            return
        }
        defer file.Close()

        // tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
        // if err != nil {
        //     fmt.Println(err)
        // }
        // defer tempFile.Close()

        fileBytes, err := ioutil.ReadAll(file)
        //fmt.Println(fileBytes)
        if err != nil {
            http.Redirect(w, r, "/" , http.StatusFound)
            return
        }

        //iterations := len(fileBytes)/256 
        // var flag int
        // if len(fileBytes)%256 == 0 {
        //     flag = 0
        // }else{
        //     flag = 1
        // }
        //var i int
        //This loop is when you use RSA encryption. Else below. For RSA, break into 256 bytes and encrypt.
        //fmt.Println("file name is ",handler.Filename)
        // for i=0;i<iterations;i++ {
        //     if i==0 {
        //         User.StoreFile(handler.Filename,fileBytes[i*256:(i+1)*256])
        //     }else{
        //         err = User.AppendFile(handler.Filename,fileBytes[i*256:(i+1)*256])
        //         if err != nil{
        //             http.Redirect(w, r, "/" , http.StatusFound)
        //             return
        //         }
        //     }
        // }
        // if (flag== 1){
        //     if i==0 {
        //         User.StoreFile(handler.Filename,fileBytes[i*256:])
        //     }else{
        //         err = User.AppendFile(handler.Filename,fileBytes[i*256:])
        //         if err != nil{
        //             http.Redirect(w, r, "/" , http.StatusFound)
        //             return
        //         }
        //     }
        // }
        //fmt.Printf("MIME Header: %+v\n", handler.Header)
        User.StoreFile(handler.Filename,fileBytes)
        global_files[username] = append(global_files[username],handler.Filename)
        http.Redirect(w, r, "/" , http.StatusFound)
        return


    }else{
        http.Error(w, "Invalid request", http.StatusInternalServerError)
        return
    }

}
func viewHandler(w http.ResponseWriter, r *http.Request) {
    val,username := check_login(w,r)
    if val == false{
        http.Redirect(w, r, "/" , http.StatusFound)
        return
    }
    if r.Method == http.MethodGet {
        var html_user MyUser 
        html_user.Name = username
        html_user.Username = username
        html_user.Image = "default"
        html_user.Files = global_files[username]
        renderHTML(w,&html_user,"view.html")
        return
    }else{
        http.Error(w, "Invalid request", http.StatusInternalServerError)
        return
    }

}
func downloadHandler(w http.ResponseWriter, r *http.Request) {
    val,_ := check_login(w,r)
    if val == false{
        http.Redirect(w, r, "/" , http.StatusFound)
        return
    }
    if r.Method == http.MethodGet {

        value := getUser(r)
        var User storeit.User 
        unmarshal_err := json.Unmarshal(value,&User)
        if(unmarshal_err!=nil){
            fmt.Println("dd")
            fmt.Println(unmarshal_err)
            http.Redirect(w, r, "/" , http.StatusFound)
            return
        }
        filename := r.URL.Query()["file"]
        //fmt.Println(filename)
        data,err := User.LoadFile(filename[0])
        if err!=nil{
            fmt.Println(err)
            http.Redirect(w, r, "/" , http.StatusFound)
            return
        }
        header := strings.Split(filename[0],".")
        //fmt.Println(mime.TypeByExtension("."+header[1]))
        w.Header().Set("Content-type", mime.TypeByExtension("."+header[1]))
        w.Write(data)
        return
    }else{
        http.Error(w, "Invalid request", http.StatusInternalServerError)
        return
    }
}
var path string

func logoutHandler(w http.ResponseWriter, r *http.Request) {
    val,_ := check_login(w,r)
    if val == false{
        http.Redirect(w, r, "/" , http.StatusFound)
        return
    }
    logout(w,r)
    http.Redirect(w, r, "/" , http.StatusFound)
    return
}
func main() {
    var list []string
    data,_ := json.Marshal(list) 
    userlib.DatastoreSet("users",data)
    path = build.Default.GOPATH
    template_folder := path+"/templates"
    global_files = make(map[string][]string)
    templates = template.Must(template.ParseFiles(template_folder+"/login.html", template_folder+"/upload.html", template_folder+"/register.html",template_folder+"/view.html"))
    http.HandleFunc("/", loginHandler)
    http.HandleFunc("/view/", viewHandler)
    http.HandleFunc("/upload/", uploadHandler)
    http.HandleFunc("/register/", registerHandler)
    http.HandleFunc("/download/", downloadHandler)
    http.HandleFunc("/logout/", logoutHandler)
    http.Handle("/static/assets/", http.StripPrefix("/static/assets/", http.FileServer(http.Dir(path+"/static/assets/"))))
    log.Fatal(http.ListenAndServe(":8090", nil))
}


