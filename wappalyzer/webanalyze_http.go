package wappalyzer

import (
	"fmt"
	"net/url"
	"time"
)

type WebAnalyzeHttpServer struct {

}

var webAnalyzeServer WebAnalyzeHttpServer

func NewWebAnalyzeServer() (*WebAnalyzeHttpServer,error) {
	was:=new(WebAnalyzeHttpServer)
	if err:=was.Init("apps.json");err!=nil{
		return nil,err
	}

	return was,nil
}

func (was *WebAnalyzeHttpServer) Init(appsFile string) error {
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