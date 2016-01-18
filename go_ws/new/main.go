package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	restful "github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"github.com/influxdb/influxdb/client"
	"github.com/jinzhu/now"
	"github.com/spf13/viper"
)

var (
	hostname     string
	port         int
	topStaticDir string
)

func init() {
	// Flags
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [default_static_dir]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&hostname, "h", "0.0.0.0", "hostname")
	flag.IntVar(&port, "p", 8080, "port")
	flag.StringVar(&topStaticDir, "static_dir", "", "static directory in addition to default static directory")
}

func appendStaticRoute(sr StaticRoutes, dir string) StaticRoutes {
	if _, err := os.Stat(dir); err != nil {
		log.Fatal(err)
	}
	return append(sr, http.Dir(dir))
}

type StaticRoutes []http.FileSystem

func (sr StaticRoutes) Open(name string) (f http.File, err error) {
	for _, s := range sr {
		if f, err = s.Open(name); err == nil {
			f = disabledDirListing{f}
			return
		}
	}
	return
}

type disabledDirListing struct {
	http.File
}

func (f disabledDirListing) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

var version = "0.0.1"

var db *client.Client

// queryDB convenience function to query the database
func queryDB(cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: "testdb",
	}
	if response, err := db.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	}
	return res, nil
}

func report_inverter(r *restful.Request, w *restful.Response) {
	sd := r.QueryParameter("sd") // startdate, "day/month/year"
	startTime := now.BeginningOfDay()
	if len(sd) > 0 {
		sdA := strings.Split(sd, "/")
		if len(sdA) == 3 {
			dy, _ := strconv.Atoi(sdA[0])
			mt, _ := strconv.Atoi(sdA[1])
			yr, _ := strconv.Atoi(sdA[2])
			startTime = time.Date(yr, time.Month(mt), dy, 0, 0, 0, 0, time.Now().Location())
		}
	}
	//	block := r.QueryParameter("block")
	tpe := strings.ToLower(r.QueryParameter("type"))
	if tpe == "month" {
		//		"SELECT block,device,field,DATE(FROM_UNIXTIME(ts)) AS ts,value "
		//		"FROM day_variable where ts >= $startTime and ts < $endTime AND "
		//		"block NOT LIKE 'ALL_BLOCK' and field='EAE_DAY' AND device LIKE '%_INV%' "
		//		"$where1 order by ts,block,device"
		influyquery := `SELECT last("value") FROM "v" WHERE "b" = 'B01' AND "f" = 'EAE_DAY' AND "d" =~ /_INV.$/ `
		influyquery += `AND (time > ` + fmt.Sprintf("%vs", now.New(startTime).BeginningOfMonth().Unix()) + `) `
		influyquery += `AND (time < ` + fmt.Sprintf("%vs", now.New(startTime).EndOfMonth().Unix()) + `) `
		influyquery += `GROUP BY time(1d), "b", "d", "f"`
		log.Println(influyquery)
		res, err := queryDB(influyquery)
		if err != nil {
			fmt.Println("error:", err.Error())
			w.WriteError(http.StatusNotAcceptable, err)
			return
		}

		log.Println(res)

		//		for i, row := range res[0].Series[0].Values {
		//			t, err := time.Parse(time.RFC3339, row[0].(string))
		//			if err != nil {
		//				log.Fatal(err)
		//			}
		//			val := row[1].(string)
		//			log.Printf("[%2d] %s: %s\n", i, t.Format(time.Stamp), val)
		//		}

		//io.WriteString(w)
		w.WriteAsJson(res)

	}

}

func main() {

	viper.SetDefault("loglevel", "debug")
	loglevel, err := log.ParseLevel(viper.GetString("loglevel"))
	if err != nil {
		panic(err)
	}
	log.SetLevel(loglevel)

	log.Println(version)

	db = NewInflux()

	ws := new(restful.WebService)
	ws.Path("/irep01").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/version").
		To(func(r *restful.Request, w *restful.Response) { io.WriteString(w, version) }).
		Doc("get igo version number").
		Operation("version"))

	ws.Route(ws.GET("/report_inverter").
		To(report_inverter).
		Doc("get inverter report data").
		Operation("report_inverter").
		Param(ws.QueryParameter("type", "month, year, ...")).
		Param(ws.QueryParameter("sd", "startdate")).
		Param(ws.QueryParameter("blocl", "list of blocks")))

	restful.Add(ws)

	config := swagger.Config{
		WebServices:    restful.DefaultContainer.RegisteredWebServices(), // you control what services are visible
		WebServicesUrl: "/",
		ApiPath:        "/apidocs.json",

		// Optionally, specifiy where the UI is located
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: "./swaggerui"}
	swagger.RegisterSwaggerService(config, restful.DefaultContainer)

	// Parse flags
	flag.Parse()
	staticDir := flag.Arg(0)

	// Setup static routes
	staticRoutes := make(StaticRoutes, 0)
	if topStaticDir != "" {
		staticRoutes = appendStaticRoute(staticRoutes, topStaticDir)
	}
	if staticDir == "" {
		staticDir = "./"
	}
	staticRoutes = appendStaticRoute(staticRoutes, staticDir)

	// Handle routes
	http.Handle("/", http.FileServer(staticRoutes))

	// Listen on hostname:port
	fmt.Printf("Listening on %s:%d...\n", hostname, port)
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", hostname, port), nil)
	if err != nil {
		log.Fatal("Error: ", err)
	}

}

func NewInflux() *client.Client {
	viper.SetDefault("influxurl", "http://127.0.0.1:8086")
	u, err := url.Parse(viper.GetString("influxurl"))
	if err != nil {
		log.Fatalln("influxdb: ", err)
	}

	conf := client.Config{
		URL:      *u,
		Username: "iplon",
		Password: "iplon321",
	}

	con, err := client.NewClient(conf)
	if err != nil {
		log.Fatalln("influxdb: ", err)
	}

	return con

}
