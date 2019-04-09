func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
    m := validPath.FindStringSubmatch(r.URL.Path)
    if m == nil {
        http.NotFound(w, r)
        return "", errors.New("Invalid Page Title")
    }
    fmt.Println(m)
    return m[2], nil 
}
func renderHTML(w http.ResponseWriter, p *Page, name string){
	err := templates.ExecuteTemplate(w,"user.html",p)
	if(err!=nil){
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func loadPage(title string) (*Page, error) {
    filename := title + ".txt"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}

func handler(w http.ResponseWriter, r *http.Request) {

    renderHTML(w,nil,"user")
}
func viewHandler(w http.ResponseWriter, r *http.Request) {
	secret(w,r);
    title := r.URL.Path[len("/view/"):]
    p, _ := loadPage(title)
    renderHTML(w,p,"view")
}
func editHandler(w http.ResponseWriter, r *http.Request) {
	login(w,r);
	//fmt.Println(r)
	getTitle(w,r)
    title := r.URL.Path[len("/edit/"):]
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title, Body: []byte("Error Occured!")}
    }
    renderHTML(w,p,"edit")
}
func (p *Page) save() error {
    filename := p.Title + ".txt"
    return ioutil.WriteFile(filename, p.Body, 0600)
}
func saveHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/save/"):]
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    err := p.save()
    if(err != nil){
    	http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

