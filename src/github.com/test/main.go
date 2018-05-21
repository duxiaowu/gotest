package main

import (
	//	"github.com/sendpush/funs"
	"fmt"
	"encoding/json"
	"strconv"
	"time"
	"net/http"
	"io/ioutil"
	"net"
	"strings"
	"log"
)
var (
	httpClient *http.Client
)
const (
	MaxIdleConns        int = 1000
	MaxIdleConnsPerHost int = 1000
	IdleConnTimeout     int = 90
)

const URL string = "http://10.10.24.190/duweibin/a.php?id=%d&name=%s"

// init HTTPClient
func init() {
	httpClient = createHTTPClient()
}

func main()  {
	fmt.Println(time.Now())
	//file := "./a.txt"
	var mapC = map[int]map[int]string{}
	for i := 1; i <= 1000; i++ {
		name := "jack" + strconv.Itoa(i)
		mapC[i] = map[int]string{i: name}
	}
	/*mapB, _ := json.Marshal(mapC)
	funcs.WriteFile(file, string(mapB))
	con := funcs.ReadFile(file)
	byt := []byte(con)
	var dat map[string]interface{}
	if err := json.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}*/
	jobs := make(chan string, 10000)
	results := make(chan int, 10000)
	//for w := 1; w <= 10; w++ {
	go worker(1, jobs, results)
	//}
	// 这里我们发送9个任务，然后关闭通道，告知任务发送完成
	for _,value := range mapC{
		for key,value1 := range value{
			url :=  fmt.Sprintf(URL, key, value1)
			jobs <- url
		}
	}
	close(jobs)

	// 然后我们从results里面获得结果
	for a := 1; a <= 6; a++ {
		fmt.Println(<-results)
	}
	fmt.Println(time.Now())
}

func worker(id int, jobs <-chan string, results chan<- int) {
	//var parm  map[string]string
	//parm["sss"] = "sss"
	for j := range jobs {
		//httpPost(j,parm)
		httpGet(j)
		results <- id
	}
}

func httpGet(url string) string{
	resp, err := http.Get(url)
	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	return string(body)
}

func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        MaxIdleConns,
			MaxIdleConnsPerHost: MaxIdleConnsPerHost,
			IdleConnTimeout:     time.Duration(IdleConnTimeout) * time.Second,
		},

		Timeout: 30 * time.Second,
	}
	return client
}

func httpPost(url string,arr map[string]string) map[string]interface{} {
	var dataArr map[string]interface{}
	postStr:=http_build_query(arr)
	req, err := http.NewRequest("POST", url, strings.NewReader(postStr))
	if err != nil {
		return dataArr
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// use httpClient to send request
	response, err := httpClient.Do(req)
	if err != nil && response == nil {
	}

	// Close the connection to reuse it
	defer response.Body.Close()
	// Let's check if the work actually is done
	// We have seen inconsistencies even when we get 200 OK response
	body, err := ioutil.ReadAll(response.Body)
	dataArr=json_decode(string(body))
	return dataArr
}

func http_build_query(arr map[string]string) string{
	var str=""
	for key,value:= range(arr) {
		str=str+key+"="+value+"&"
	}
	if str!="" {
		str = str[0 : len(str)-1]
	}
	return str
}

func json_decode(jsonStr string) map[string]interface{} {
	var dat map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &dat)
	if err!=nil {
		msg:="json decode errror,err:"+err.Error()+",jsonStr:"+jsonStr
		log.Println(msg)
	}
	return dat
}