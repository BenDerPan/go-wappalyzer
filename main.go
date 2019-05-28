package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/benderpan/go-wappalyzer/wappalyzer"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
)

var (
	runType    string
	update     bool
	outputFile string
	workers    int
	apps       string
	host       string
	hosts      string
	crawlCount int
	listen     string
	verbose    bool
)

func init() {
	flag.StringVar(&runType,"type","standalone","run type, standalone - single tool model, web - web server model")
	flag.StringVar(&outputFile, "output", "data.json", "output file")
	flag.BoolVar(&update, "update", false, "update apps file")
	flag.IntVar(&workers, "worker", 4, "number of concurrent worker")
	flag.StringVar(&apps, "apps", "apps.json", "app definition file.")
	flag.StringVar(&host, "host", "", "single host to test")
	flag.StringVar(&hosts, "hosts", "", "filename with hosts, one host per line.")
	flag.IntVar(&crawlCount, "crawl", 0, "links to follow from the root page (default 0)")
	flag.StringVar(&listen,"listen","0.0.0.0:8080","Web server listen address")
	flag.BoolVar(&verbose,"verbose",true,"show verbose message")


	if cpu := runtime.NumCPU(); cpu == 1 {
		runtime.GOMAXPROCS(2)
	} else {
		runtime.GOMAXPROCS(cpu)
	}
}

func main() {

	flag.Parse()
	if runType =="web" {
		//run as web server, which will provide restfull api.
		was, err := wappalyzer.NewWebAnalyzeServer(apps)
		if err != nil {
			fmt.Printf("Error:%v", err)
			return
		}

		http.Handle("/", http.FileServer(FS(false)))
		//http.Handle("/",http.StripPrefix("/", http.FileServer(http.Dir("./www/"))))
		http.HandleFunc("/analyze",was.HttpHandlerAnalyze)
		err=http.ListenAndServe(listen,nil)
		if err!=nil {
			log.Fatal("open web server failed, err: ", err)
		}

	}else{
		var file io.ReadCloser
		var err error

		if !update && host == "" && hosts == "" {
			flag.Usage()
			return
		}

		if update {
			err = wappalyzer.DownloadFile(wappalyzer.WappalyzerURL, "apps.json")
			if err != nil {
				log.Fatalf("error: can not update apps file: %v", err)
			}

			log.Println("app definition file updated from ", wappalyzer.WappalyzerURL)

			if host == "" && hosts == "" {
				return
			}

		}

		// check single host or hosts file
		if host != "" {
			file = ioutil.NopCloser(strings.NewReader(host))
		} else {
			file, err = os.Open(hosts)
			if err != nil {
				log.Fatalf("error: can not open host file %s: %s", hosts, err)
			}
		}
		defer file.Close()

		results, err := wappalyzer.Init(workers, file, apps, crawlCount)

		if err != nil {
			log.Fatal("error initializing: ", err)
		}

		log.Printf("Scanning with %v workers.", workers)

		outFile, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal("output file create failed. err: ", err)
		}
		defer outFile.Close()

		for result := range results {
			if result.Error != "" {
				if verbose{
					log.Printf("Error: Host=%s  Msg=%s", result.Host, result.Error)
				}

			}

			jsonValue, err := json.Marshal(result)
			if err != nil {
				if verbose{
					log.Print("{}\n")
				}
			} else {
				jsonValue = append(jsonValue, '\n')
				_,err= outFile.Write(jsonValue)
				if err!=nil{
					log.Fatal("write message to output file failed, err: ",err)
				}
				if verbose{
					log.Print(string(jsonValue))
				}


			}

		}

	}






}