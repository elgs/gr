package main

import (
	"encoding/json"
	"fmt"
	"github.com/elgs/gorest"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	input := args()
	config := parseConfig(input[0])
	if config == nil {
		return
	}
	ds := config["data_source"].(string)
	dbo := &gorest.MySqlDataOperator{
		Ds:              ds,
		DbNameExtractor: extractDbNameFromDs,
	}
	r := &gorest.Gorest{
		EnableHttp: config["enable_http"].(bool),
		HostHttp:   config["host_http"].(string),
		PortHttp:   uint16(config["port_http"].(float64)),

		EnableHttps:   config["enable_https"].(bool),
		HostHttps:     config["host_https"].(string),
		PortHttps:     uint16(config["port_https"].(float64)),
		CertFileHttps: config["cert_file_https"].(string),
		KeyFileHttps:  config["key_file_https"].(string),

		UrlPrefix: config["url_prefix"].(string),
		Dbo:       dbo}
	r.Serve()
}

func extractDbNameFromDs(ds string) string {
	a := strings.LastIndex(ds, "/")
	b := ds[a+1:]
	c := strings.Index(b, "?")
	if c < 0 {
		return b
	}
	return b[:c]
}

func parseConfig(configFile string) map[string]interface{} {
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Println(configFile, "not found")
		return nil
	}
	var config map[string]interface{}
	if err := json.Unmarshal(b, &config); err != nil {
		fmt.Println("Error parsing", configFile)
		return nil
	}
	return config
}

func args() []string {
	ret := []string{}
	if len(os.Args) <= 1 {
		ret = append(ret, "gorest.json")
	} else {
		ret = os.Args[1:]
	}
	return ret
}
