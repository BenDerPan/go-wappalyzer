package main

import (
	"encoding/json"
	"fmt"
	"github.com/benderpan/go-wappalyzer/wappalyzer"
	"os"
)

func main() {
	was,err:=wappalyzer.NewWebAnalyzeServer()
	if err!=nil{
		fmt.Printf("Error:%v",err)
		return
	}

	result:=was.Analyze("https://github.com")
	jsonValue, err := json.Marshal(result)
	if err != nil {
		os.Stdout.Write([]byte("{}\n"))
	} else {
		jsonValue = append(jsonValue, '\n')
		os.Stdout.Write(jsonValue)
	}

}