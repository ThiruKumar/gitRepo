package db

import (
  "time"
  "github.com/influxdb/influxdb/client/v2"
  "log"
  "net/http"
  "fmt"
 // "url"
)

const (
    MyDB = "iplonDb"
    username = "admin"
    password = "admins"	
)


func GetConnection() (c client.Client, err error){

// Make client
    clnt,conErr := client.NewHTTPClient(client.HTTPConfig{
        Addr: "http://localhost:8086",
        Username: username,
        Password: password,
    })
    
    return clnt,conErr

}

func PutData(rspns http.ResponseWriter,rqst *http.Request){
	

name2 :=rqst.FormValue("onedata")
secData :=rqst.FormValue("secdata")

    //c,conErr := GetConnection()    
    
    //if conErr != nil {
		//log.Fatal(conErr)
	//}
	
	c,clntErr :=  GetConnection()

    if clntErr != nil {
		log.Fatal(clntErr)
	}
	

    // Create a new point batch
    bp,batchErr := client.NewBatchPoints(client.BatchPointsConfig{
        Database:  MyDB,
        Precision: "s",
    })
    
        if batchErr != nil {
		log.Fatal(batchErr)
	}

    // Create a point and add to batch
    tags := map[string]string{"host": "primary",
		"plant": name2,
		"second": secData,
		}
    fields := map[string]interface{}{
        "power":   27.8,
    }
    pt,ptErr := client.NewPoint("power", tags, fields, time.Now())
    bp.AddPoint(pt)
    
            if ptErr != nil {
		log.Fatal(ptErr)
	}

    // Write the batch
    c.Write(bp)
http.Redirect(rspns, rqst, "faq.html", 301)
log.Println("data written suucessfully...")

}


// queryDB convenience function to query the database
func queryDB(clnt client.Client, cmd string) (res []client.Result, err error) {
    q := client.Query{
        Command:  cmd,
        Database: MyDB,
    }
    if response, err := clnt.Query(q); err == nil {
        if response.Error() != nil {
            return res, response.Error()
        }
        res = response.Results
    }
    return res, nil
}

func GetData(MyMeasurement string) (endRes string){
			q := fmt.Sprintf("SELECT * FROM %s LIMIT %d", MyMeasurement, 20)
			var val string
			
			//clnt,conErr := GetConnection()    
    
			//if conErr != nil {
			//	log.Fatal(conErr)
			//}
			
			
	c, clntErr :=  GetConnection()
    if clntErr != nil {
		log.Fatal(clntErr)
	}
			
			res, err := queryDB(c, q)
			if err != nil {
				log.Fatal(err)
			}

			for i, row := range res[0].Series[0].Values {
				t, err := time.Parse(time.RFC3339, row[0].(string))
				if err != nil {
					log.Fatal(err)
				}
			   val = row[1].(string)
			   
				//fmt.Fprintf(w, "Hi there, I love %s!", val)
				log.Printf("[%2d] %s: %s\n", i, t.Format(time.Stamp), val)
				
			}
				
						
			return (val) 
			 			
}


