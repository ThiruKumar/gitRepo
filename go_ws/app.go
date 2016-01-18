package main

import (
  "html/template"
  "net/http"
  "os"
  "path"
  "log"
  "encoding/json"
  //"time"
  "db"
  //"github.com/influxdb/influxdb/client/v2"
)


func main() {
	
  http.HandleFunc("/PutData",func(w http.ResponseWriter, r *http.Request) {
     db.PutData(w,r)
})




  http.HandleFunc("/getData",func(w http.ResponseWriter, r *http.Request) {
  
   // res := db.GetData("power")
    
   type Todo struct {
    Name      string    `json:"name"`
    Completed bool      `json:"completed"`
}

type Todos []Todo
    
    
        todos := Todos{
        Todo{Name: "Write presentation"},
        Todo{Name: "Host meetup"},
    }

    json.NewEncoder(w).Encode(todos)
    
 
})





  fs := http.FileServer(http.Dir("static"))
  http.Handle("/static/", http.StripPrefix("/static/", fs))
  fonts := http.FileServer(http.Dir("fonts"))
  http.Handle("/fonts/", http.StripPrefix("/fonts/", fonts))
  images := http.FileServer(http.Dir("images"))
  http.Handle("/images/", http.StripPrefix("/images/", images))
  scripts := http.FileServer(http.Dir("scripts"))
  http.Handle("/scripts/", http.StripPrefix("/scripts/", scripts))
  styles := http.FileServer(http.Dir("styles"))
  http.Handle("/styles/", http.StripPrefix("/styles/", styles))
  templates := http.FileServer(http.Dir("templates"))
  http.Handle("/templates/", http.StripPrefix("/templates/", templates))
  
  
  //log.Println(fs)
  
  http.HandleFunc("/", serveTemplate)

  log.Println("Listening...")
  http.ListenAndServe(":3010", nil)
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
  lp := path.Join("templates", "layout.html")
 // fp := path.Join("templates", r.URL.Path)


  // Return a 404 if the template doesn't exist
  info, err := os.Stat(lp)
  if err != nil {
    if os.IsNotExist(err) {
      http.NotFound(w, r)
      return
    }
  }

  // 404 error for the requests that asks for directory
  if info.IsDir() {
    http.NotFound(w, r)
    return
  }

  //tmpl, err := template.ParseFiles(lp, fp)
   tmpl, err := template.ParseFiles(lp)
  if err != nil {
    // Log the detailed error
    log.Println(err.Error())
    // Return a generic "Internal Server Error" message
    http.Error(w, http.StatusText(500), 500)
    return
  }
type Person struct {
    UserName string
//    Age  int
}

p := Person{UserName: "{{themeActive}}"}
  if err := tmpl.ExecuteTemplate(w, "layout", p); err != nil {
    log.Println(err.Error())
    http.Error(w, http.StatusText(500), 500)
  }
}
