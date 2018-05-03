package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const (
	spider_base_url string = "http://xxx/zabbix.php?action=dashboard.view"
	login_url       string = "http://xxx/index.php"
	api_url         string = "http://xxx/api_jsonrpc.php"
	graph_url       string = "http://xxx/chart2.php"
	username        string = "xxx"
	password        string = "xxx"
)

var cookies []*http.Cookie

type AragsHost struct {
	Ip        string              `json:"ip"`
	Url       string              `json:"url"`
	Headermap map[string][]string `json:"headermap"`
	Authid    string              `json:"authid"`
}

func NewArgsHost(url, authid, ip string, headermap map[string][]string) *AragsHost {
	return &AragsHost{ip, url, headermap, authid}
}

type AragsGraph struct {
	AragsHost
	Hostid string `json:"hostid"`
}

func NewArgsGraph(hostid string, aragshost *AragsHost) *AragsGraph {
	return &AragsGraph{*aragshost, hostid}
}

func loginrpc(url string) (string, map[string][]string) {

	//json序列化
	data := fmt.Sprintf(`{
            "jsonrpc": "2.0",
            "method": "user.login",
            "params": {
                "user": "%s",
                "password": "%s"
                },
            "id": 0
    }`, username, password)

	var jsonStr = []byte(data)

	//提交请求
	reqest, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))

	if err != nil {
		panic(err)
	}
	reqest.Header.Set("Content-Type", "application/json")
	//生成client 参数为默认
	client := &http.Client{}
	//处理返回结果
	response, _ := client.Do(reqest)
	defer response.Body.Close()

	//fmt.Println("response Status:", response.Status)
	//fmt.Println("response Headers:", response.Header)
	headermap := response.Header
	body, _ := ioutil.ReadAll(response.Body)
	//fmt.Println("response Body:", string(body))
	//反序列化
	type Resp struct {
		Jsonrpc string `json:"jsonrpc"`
		Result  string `json:"result"`
		Id      int    `json:id`
	}
	var resp Resp
	var authID string
	if err := json.Unmarshal(body, &resp); err == nil {
		//fmt.Println(resp.Result)
		authID = resp.Result
	}
	return authID, headermap
}

func get_data(url string, headermap map[string][]string, data string) *http.Response {
	reqest, _ := http.NewRequest("GET", url, bytes.NewBuffer([]byte(data)))
	for k, v := range headermap {
		//fmt.Println(k, v[0])
		reqest.Header.Add(k, v[0])
	}
	client := &http.Client{}
	response, _ := client.Do(reqest)
	//defer response.Body.Close()
	//body, _ := ioutil.ReadAll(response.Body)
	//fmt.Println("response Body:", string(body))
	return response
}
func get_hostid(araghost *AragsHost) string {

	data := fmt.Sprintf(`{
        "jsonrpc": "2.0",
        "method": "host.get",
        "params": {
            "output":["name","status","host","groups"],
            "selectGroups":"extend",
            "filter": {"ip": ["%s"]}},
            "auth": "%s",
            "id": 1}`, araghost.Ip, araghost.Authid)
	//fmt.Println(data)
	res := get_data(araghost.Url, araghost.Headermap, data)
	body, _ := ioutil.ReadAll(res.Body)
	//fmt.Println("response Body:", string(body))
	type Group struct {
		Groupid  string `json:"groupid"`
		Name     string `json:"name"`
		Internal string `json:"internal"`
		Flags    string `json:"flags"`
	}
	type Result struct {
		Hostid string  `json:"hostid"`
		Name   string  `json:"name"`
		Status string  `json:"status"`
		Host   string  `json:"host"`
		Groups []Group `json:"groups"`
	}
	type Resp struct {
		Jsonrpc string   `json:"jsonrpc"`
		Result  []Result `json:"result"`
		Id      int      `json:id`
	}

	var resp Resp
	var hostid string
	if err := json.Unmarshal(body, &resp); err == nil {
		if len(resp.Result) != 0 {
			hostid = resp.Result[0].Hostid
		} else {
			hostid = "0"
		}
	}
	return hostid
}

func get_graphid(aragsgraph *AragsGraph) map[string]string {
	//graphid := "0"
	graph_map := make(map[string]string)
	if aragsgraph.Hostid != "0" {
		data := fmt.Sprintf(`{
                    "jsonrpc": "2.0",
                    "method": "graph.get",
                    "params": {
                        "output": "extend",
                        "hostids": "%s",
                        "sortfield": "name"
                    },
                    "auth": "%s",
                    "id": 1}`, aragsgraph.Hostid, aragsgraph.Authid)
		//fmt.Println(data)
		res := get_data(aragsgraph.Url, aragsgraph.Headermap, data)
		body, _ := ioutil.ReadAll(res.Body)
		//fmt.Println("response Body:", string(body))
		type Result struct {
			Graphid          string `json:graphid`
			Name             string `json:name`
			Width            string `json:width`
			Height           string `json:height`
			Yaxismin         string `json:yaxismin`
			Yaxismax         string `json:yaxismax`
			Templateid       string `json:templateid`
			Show_work_period string `json:show_work_period`
			Show_triggers    string `json:show_triggers`
			Graphtype        string `json:graphtype`
			Show_legend      string `json:show_legend`
			Show_3d          string `json:show_3d`
			Percent_left     string `json:percent_left`
			Percent_right    string `json:percent_right`
			Ymin_type        string `json:ymin_type`
			Ymax_type        string `json:ymax_type`
			Ymin_itemid      string `json:ymin_itemid`
			Ymax_itemid      string `json:ymax_itemid`
			Flags            string `json:flags`
		}
		type Resp struct {
			Jsonrpc string   `json:"jsonrpc"`
			Result  []Result `json:"result"`
			Id      int      `json:id`
		}
		var resp Resp

		regex_arr := []string{"Mysql_RW", "Network", "MySQL_Thread", "Seconds_Behind_Master", "Network_MySQL", "Cpu_Load", "Tcp_conect", "CPU_Used"}
		if err := json.Unmarshal(body, &resp); err == nil {
			//fmt.Println("graphid:", resp.Result)
			sc_arr := resp.Result
			for _, graph := range sc_arr {
				//fmt.Println(graph.Graphid, graph.Name)
				for _, ra := range regex_arr {
					if ra == graph.Name {
						graph_map[graph.Name] = graph.Graphid
					}
				}
			}
		}
	}
	return graph_map
}

func graphids(ip string) map[string]string {
	//生成要访问的url
	url := api_url
	//ip := "10.181.142.187"
	//fmt.Println(user_login(url))
	authid, headermap := loginrpc(url)
	argshost := NewArgsHost(url, authid, ip, headermap)
	//hostid := get_hostid(url, headermap, authid, "172.20.82.68")
	hostid := get_hostid(argshost)
	//arags := &Arags{Ip: "172.20.82.68", Url: url, Headermap: headermap, Authid: authid}
	//fmt.Println(*argshost)
	argsgraph := NewArgsGraph(hostid, argshost)
	graphids := get_graphid(argsgraph)
	//graphids := get_graphid(url, hostid, headermap, authid)
	return graphids
}

func getResultHtml(get_url string) *http.Response {
	c := &http.Client{}
	Jar, _ := cookiejar.New(nil)
	getURL, _ := url.Parse(get_url)
	Jar.SetCookies(getURL, cookies)
	c.Jar = Jar
	res, _ := c.Get(get_url)
	return res
}
func login() {
	//获取登陆界面的cookie
	c := &http.Client{}
	req, err := http.NewRequest("POST", login_url, nil)
	res, err := c.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	var temp_cookies = res.Cookies()
	for _, v := range res.Cookies() {
		req.AddCookie(v)
	}
	//post数据
	postValues := url.Values{}
	postValues.Add("name", username)
	postValues.Add("password", password)
	postValues.Add("autologin", "1")
	postValues.Add("enter", "Sign in")
	postURL, _ := url.Parse(login_url)
	Jar, _ := cookiejar.New(nil)
	Jar.SetCookies(postURL, temp_cookies)
	c.Jar = Jar
	res, _ = c.PostForm(login_url, postValues)
	cookies = res.Cookies()
	//data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	//fmt.Println(string(data))
}

func getgraph(ip string) map[string]string {
	login()
	graphids := graphids(ip)
	results := make(map[string]string)
	for k, v := range graphids {
		//fmt.Println("Screen", k, fmt.Sprintf("%s?graphid=%s&width=500&height=100&legend=1&updateProfile=1&profileIdx=web.screens&profileIdx2=812&period=86400&stime=20180815102116&sid=4bb26f2166536de4&curtime=1502850087598", graph_url, v))
		url := fmt.Sprintf("%s?graphid=%s&width=500&height=100&legend=1&updateProfile=1&profileIdx=web.screens&profileIdx2=812&period=86400&stime=20180815102116&sid=4bb26f2166536de4&curtime=1502850087598", graph_url, v)
		results[k] = url
	}
	return results
}

var (
	GraphResults map[string]*GraphResult
)

type Gata struct {
	Data map[string]string `json:data`
}

type GraphResult struct {
	Code    int    `json:code`
	Message string `json:message`
	Gata
}

func NewGata(data map[string]string) *Gata {
	return &Gata{data}
}

func NewGraphsResult(code int, message string, gata *Gata) *GraphResult {
	return &GraphResult{code, message, *gata}
}

func main() {
	//ip := "10.181.142.187"
	ip := "10.191.65.126"
	url_map := getgraph(ip)
	gatamap := NewGata(url_map)
	var graphresult *GraphResult
	if true {
		graphresult = NewGraphsResult(0, "message", gatamap)
	} else {
		graphresult = NewGraphsResult(1, "message", gatamap)
	}

	//fmt.Println(gatamap)
	fmt.Println(*graphresult)
}
