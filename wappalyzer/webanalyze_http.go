package wappalyzer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type WebAnalyzeHttpServer struct {

}

type ResponseMsg struct {
	Msg string `json:"msg"`
	Code int `json:"code"`
	Data string `json:"data"`
}

var webAnalyzeServer WebAnalyzeHttpServer

func NewWebAnalyzeServer(appFile string) (*WebAnalyzeHttpServer,error) {
	was:=new(WebAnalyzeHttpServer)
	if err:=was.init(appFile);err!=nil{
		return nil,err
	}
	return was,nil
}

func (was *WebAnalyzeHttpServer) init(appsFile string) error {
	//加载 应用识别定义字典
	if err:=loadApps(appsFile);err!=nil{
		return err
	}
	return nil
}

func (was *WebAnalyzeHttpServer) Analyze(url string) (Result) {
	job:=NewOnlineJob(url, "", nil, 0)
	result:=dowork(job)
	return result
}

func (was *WebAnalyzeHttpServer) HttpHandlerAnalyze(w http.ResponseWriter,r *http.Request)  {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers","*")
	}

	if r.Method == "OPTIONS" {
		return
	}
	respMsg:=ResponseMsg{}
	if r.Method!="POST"{
		respMsg.Code=100
		respMsg.Msg=fmt.Sprintf("Error:Ivalid Request!")
		respMsg.Data=""
	}else{
		err:=r.ParseForm()
		if err!=nil{
			respMsg.Code=100
			respMsg.Msg=fmt.Sprintf("Error:Ivalid Form Data!")
			respMsg.Data=""
		}else{
			url:=r.PostForm.Get("url")
			if url!=""{
				result:=was.Analyze(string(url))
				jsonValue, err := json.Marshal(result)

				if err!=nil{
					respMsg.Code=1
					respMsg.Msg=fmt.Sprintf("Error:%v",err)
					respMsg.Data=""
				}else {
					respMsg.Code=0
					respMsg.Msg="ok"
					respMsg.Data=string(jsonValue)
				}
			}else{
				respMsg.Code=100
				respMsg.Msg=fmt.Sprintf("Error:Ivalid Form Data!")
				respMsg.Data=""
			}
		}

	}
	w.Header().Set("Content-Type", "application/json")
	respMsgJson,_:=json.Marshal(respMsg)
	w.Write(respMsgJson)
}

func dowork(job *Job) Result {
	u, err := url.Parse(job.URL)
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	job.URL = u.String()

	t0 := time.Now()
	result, err := process(job)
	t1 := time.Now()
	result.Duration = t1.Sub(t0)
	if err != nil {
		result.Error = fmt.Sprintf("%s", err)
	} else {
		result.Error = ""
	}

	return result
}