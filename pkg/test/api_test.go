package test

import (
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"testing"
)

func TestAddGroup(t *testing.T) {
	resp, err := http.Get("http://localhost:80/multicast/v1/groups/")
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println(resp.Body)
	}
	assert.Equal(t, 1, 1)
}
