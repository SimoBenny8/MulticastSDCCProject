package test

import (
	"encoding/json"
	"fmt"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/restApi"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/util"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestAddGroups(t *testing.T) {

	AddGroupTesting("8080")
	AddGroupTesting("8081")
	AddGroupTesting("8082")
}

func AddGroupTesting(host string) {

	obj := restApi.MulticastReq{
		MulticastId:   "BM",
		MulticastType: util.BMULTICAST,
	}
	jsn, _ := json.Marshal(obj)
	url := "http://localhost:" + host + "/multicast/v1/groups"
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
	url := "http://localhost:" + "8080" + "/multicast/v1/groups/" + "Q"
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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func TestSendMessage(t *testing.T) {
	url := "http://localhost:" + "8080" + "/multicast/v1/messaging/" + "Q"
	method := "POST"
	m := []byte("message")
	//obj := Message{m}
	jsonMessage, _ := json.Marshal(m)
	reader := strings.NewReader(string(jsonMessage))
	client := &http.Client{}
	req, err := http.NewRequest(method, url, reader)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/jsonMessage")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func TestGetMessages(t *testing.T) {
	url := "http://localhost:" + "8080" + "/multicast/v1/messaging/" + "Q"
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
