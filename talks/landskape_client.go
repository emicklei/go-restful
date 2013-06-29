package main

import (
	"io/ioutil"
	"net/http"
)

func main() {
	resp, _ := http.Get("http://localhost:9090/test/systems/")
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	println(string(body))
}
