package test

import (
	"encoding/json"
	"fmt"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/restApi"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/util"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"testing"
)

func TestAddGroup(t *testing.T) {

	obj := restApi.MulticastReq{
		MulticastId:   string(rune(rand.Intn(100))),
		MulticastType: util.BMULTICAST,
	}
	jsn, _ := json.Marshal(obj)
	url := "http://localhost:8080/multicast/v1/groups"
	method := "POST"
	payload := strings.NewReader(string(jsn))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func TestStartGroup(T *testing.T) {
	url := "http://localhost:" + "8080" + "/multicast/v1/groups/" + "Q"
	method := "PUT"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func TestGetInfoGroup(t *testing.T) {
	url := "http://localhost:" + "8080" + "/multicast/v1/groups/" + "Q"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func TestCloseGroup(t *testing.T) {
	url := "http://localhost:" + "8080" + "/multicast/v1/groups/" + "r"
	method := "DELETE"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
