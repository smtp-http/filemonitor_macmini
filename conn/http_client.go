package conn

import (
	"sync"
	"bytes"
	"fmt"
	"net/url"
    "net/http"
    "io/ioutil"
)

type HttpClient struct {
	url string
}

func (c *HttpClient) HttpSetUrl (url string) bool {

	if url != "" {
		c.url = url
		return true
	}

	return false
}



var instance *HttpClient
var once sync.Once
 
func GetHttpClientInstance() *HttpClient {
    once.Do(func() {
        instance = &HttpClient{}
    })
    return instance
}


func (c *HttpClient) HttpGet() {
	resp, err := http.Get("http://www.01happy.com/demo/accept.php?id=1")
	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}

func (c *HttpClient) HttpPost(b []byte) {
   

	body := bytes.NewBuffer(b)

	fmt.Printf("b: %v\n",b)
	resp, err := http.Post(c.url,"application/json;charset=utf-8",body)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	_, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		// handle error
	}
}

func (c *HttpClient) HttpPostForm() {
	resp, err := http.PostForm(c.url,
		url.Values{"key": {"Value"}, "id": {"123"}})

	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))

}
