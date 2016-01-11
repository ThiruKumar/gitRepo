package main

import (
  "html/template"
  "net/http"
  "os"
  "path"
  "log"
  "time"
  "db"
  //"github.com/influxdb/influxdb/client/v2"
)


func main() {
	
  
  //http.HandleFunc("/putdata.html/",db.PutData)
 // http.HandleFunc("/PutData/", db.PutData)
 
  //db.Connection()
  
  http.HandleFunc("/PutData",func(w http.ResponseWriter, r *http.Request) {
   // r.ParseForm()
   // db.PutData(r.FormValue("onedata"))
     db.PutData(w,r)
})

  http.HandleFunc("/getdata.html",func(w http.ResponseWriter, r *http.Request) {
   
   type person struct {
        name string
}
   
    res := db.GetData("power")
    
    var val string

for i, row := range res[0].Series[0].Values {
				t, err := time.Parse(time.RFC3339, row[0].(string))
				if err != nil {
					log.Fatal(err)
				}
			   val = row[2].(string)
				//fmt.Fprintf(w, "Hi there, I love %s!", val)
				log.Printf("[%2d] %s: %s\n", i, t.Format(time.Stamp), val)
				
			
}  


type Person struct {
    UserName string
//    Age  int
}

p := Person{UserName: val}
//var P Person  

//P.Name = "Thiru"  
//P.Age = 25     
    
//var message string = "Hi ths is my"    
lp := path.Join("templates", "layout.html")
  fp := path.Join("templates", r.URL.Path)
log.Println(r.URL.Path)

  // Return a 404 if the template doesn't exist
  info, err := os.Stat(fp)
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

  tmpl, err := template.ParseFiles(lp, fp)
  if err != nil {
    // Log the detailed error
    log.Println(err.Error())
    // Return a generic "Internal Server Error" message
    http.Error(w, http.StatusText(500), 500)
    return
  }

  if err := tmpl.ExecuteTemplate(w, "layout", p); err != nil {
    log.Println(err.Error())
    http.Error(w, http.StatusText(500), 500)
  }    
    
  

    
})
  
  

  fs := http.FileServer(http.Dir("static"))
  http.Handle("/static/", http.StripPrefix("/static/", fs))
  http.HandleFunc("/", serveTemplate)

  log.Println("Listening...")
  http.ListenAndServe(":3010", nil)

}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
  lp := path.Join("templates", "layout.html")
  fp := path.Join("templates", r.URL.Path)


  // Return a 404 if the template doesn't exist
  info, err := os.Stat(fp)
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

  tmpl, err := template.ParseFiles(lp, fp)
  if err != nil {
    // Log the detailed error
    log.Println(err.Error())
    // Return a generic "Internal Server Error" message
    http.Error(w, http.StatusText(500), 500)
    return
  }

  if err := tmpl.ExecuteTemplate(w, "layout", nil); err != nil {
    log.Println(err.Error())
    http.Error(w, http.StatusText(500), 500)
  }
}
