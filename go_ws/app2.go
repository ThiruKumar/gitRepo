package main


import (
   // "net/url"
   // "fmt"
    "log"
   // "os"
   "time"

    "github.com/influxdb/influxdb/client/v2"
)

const (
    MyDB = "iplonDb"
    username = "admin"
    password = "admins"
)

func main() {
    // Make client
    c,conErr := client.NewHTTPClient(client.HTTPConfig{
        Addr: "http://localhost:8086",
        Username: username,
        Password: password,
    })
    
    if conErr != nil {
		log.Fatal(conErr)
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
		"plant": "chennai2",
		}
    fields := map[string]interface{}{
        "power":   14.9,
    }
    pt,ptErr := client.NewPoint("power", tags, fields, time.Now())
    bp.AddPoint(pt)
    
            if ptErr != nil {
		log.Fatal(ptErr)
	}

    // Write the batch
    c.Write(bp)
    
    
    
    
    
    
}



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
