package main

import (
    "go/build"
    "log"
    "net/http"
    "html/template"
    "regexp"
    "encoding/json"
    "storeit"
    "github.com/fenilfadadu/CS628-assn1/userlib"
    "fmt"
    "io/ioutil"
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


func renderHTML(w http.ResponseWriter, p *MyUser, name string){
    err := templates.ExecuteTemplate(w,name,p)
    if(err!=nil){
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
func registerHandler (w http.ResponseWriter, r *http.Request) {
    fmt.Println("called")
    if r.Method == http.MethodGet {
        renderHTML(w,nil,"register.html")

    }else if r.Method == http.MethodPost {
        username := r.FormValue("username")
        password := r.FormValue("password")
        User,err := storeit.InitUser(username,password)
        fmt.Println(err)
        marshalled_user_struct, err := json.Marshal(&User)
        userlib.DatastoreSet(username,marshalled_user_struct)
        if err!=nil{
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        var html_user MyUser 
        html_user.Name = username
        html_user.Username = username
        html_user.Image = "default"
        renderHTML(w,&html_user,"view.html")
    }else{
        http.Error(w, "Invalid request", http.StatusInternalServerError)
    }
}
func loginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        renderHTML(w,nil,"login.html")
    }else if r.Method == http.MethodPost {
        username := r.FormValue("username")
        password := r.FormValue("password")
        User,err := storeit.GetUser(username,password)
        fmt.Println(err)
        marshalled_user_struct, err := json.Marshal(&User)
        userlib.DatastoreSet(username,marshalled_user_struct)
        if err!=nil{
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        var html_user MyUser 
        html_user.Name = username
        html_user.Username = username
        html_user.Image = "default"
        html_user.Files = []string{"hello.pdf"}
        renderHTML(w,&html_user,"view.html")
    }else{
        http.Error(w, "Invalid request", http.StatusInternalServerError)
    }

}
func uploadHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("method:", r.Method)
    fmt.Println("haha: ",http.MethodPost)
    if r.Method == http.MethodGet {
        renderHTML(w,nil,"upload.html")
    }else if r.Method == http.MethodPost {

        value, _ := userlib.DatastoreGet("vinayakt")
        var User storeit.User 
        unmarshal_err := json.Unmarshal(value,&User)
        if(unmarshal_err!=nil){
            http.Error(w, "Invalid request", http.StatusInternalServerError)
        }
        r.ParseMultipartForm(10 << 20)
    
        file, handler, err := r.FormFile("myFile")
        if err != nil {
            fmt.Println("Error Retrieving the File")
            fmt.Println(err)
            return
        }
        defer file.Close()

        // tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
        // if err != nil {
        //     fmt.Println(err)
        // }
        // defer tempFile.Close()

        fileBytes, err := ioutil.ReadAll(file)
        fmt.Println(fileBytes)
        if err != nil {
            fmt.Println(err)
        }
    // write this byte array to our temporary file
    //tempFile.Write(fileBytes)
    // return that we have successfully uploaded our file!
        iterations := len(fileBytes)/256 
        var flag int
        if len(fileBytes)%256 == 0 {
            flag = 0
        }else{
            flag = 1
        }
        var i int
        fmt.Println("file name is ",handler.Filename)
        for i=0;i<iterations;i++ {
            if i==0 {
                User.StoreFile(handler.Filename,fileBytes[i*256:(i+1)*256])
            }else{
                err = User.AppendFile(handler.Filename,fileBytes[i*256:(i+1)*256])
                if err != nil{
                    fmt.Println("wrong :",err)
                }
            }
        }
        if (flag== 1){
            if i==0 {
                User.StoreFile(handler.Filename,fileBytes[i*256:])
            }else{
                err = User.AppendFile(handler.Filename,fileBytes[i*256:])
                if err != nil{
                    fmt.Println("wrong :",err)
                }
            }
        }
        //fmt.Printf("MIME Header: %+v\n", handler.Header)
        http.Error(w, "Good job", http.StatusInternalServerError)

    }else{
        http.Error(w, "Invalid request", http.StatusInternalServerError)
    }

}
func viewHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        var html_user MyUser 
        html_user.Name = "vinayakt"
        html_user.Username = "vinayakt"
        html_user.Image = "default"
        html_user.Files = []string{"hello.pdf"}
        renderHTML(w,&html_user,"view.html")
        renderHTML(w,nil,"view.html")
    }else{
        http.Error(w, "Invalid request", http.StatusInternalServerError)
    }

}
func downloadHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        value, _ := userlib.DatastoreGet("vinayakt")
        var User storeit.User 
        unmarshal_err := json.Unmarshal(value,&User)
        if(unmarshal_err!=nil){
            http.Error(w, "Invalid request", http.StatusInternalServerError)
        }
        filename := "hello.pdf"
        data,err := User.LoadFile(filename)
        if err!=nil{
            fmt.Println("laudap ",err)
        }
        w.Header().Set("Content-type", "application/pdf")
        w.Write(data)

    }else{
        http.Error(w, "Invalid request", http.StatusInternalServerError)
    }
}
var path string
func main() {
    path = build.Default.GOPATH
    template_folder := path+"/templates"
    templates = template.Must(template.ParseFiles(template_folder+"/login.html", template_folder+"/upload.html", template_folder+"/register.html",template_folder+"/view.html"))
    http.HandleFunc("/", loginHandler)
    http.HandleFunc("/view/", viewHandler)
    http.HandleFunc("/upload/", uploadHandler)
    http.HandleFunc("/register/", registerHandler)
    http.HandleFunc("/download/", downloadHandler)
    http.Handle("/static/assets/", http.StripPrefix("/static/assets/", http.FileServer(http.Dir(path+"/static/assets/"))))
    log.Fatal(http.ListenAndServe(":8080", nil))
}


