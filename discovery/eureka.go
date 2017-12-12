package discovery

import (
	"errors"
	"net/http"
	"io/ioutil"
	"fmt"
	"time"
	"net"
)

//TODO: This should become a discovery interface. And Eureka just the first implementation
type EurekaClient struct {
	eurekaUrl string
}

//TODO: Creating our own error type and wrapping standard net/http errors could be useful to prevent
// the original errors from being lost
var errEurekaTimesOut = errors.New("Eureka server timed out")
var errNoEurekaConnection = errors.New("Unable to reach Eureka server")
var errEurekaUnexpectedHttpResponseCode = errors.New("Eureka returned a non 200 http response code")
const eurekaClientTimeoutInSeconds  = 10

func NewEurekaClient(eurekaUrl string) (ec EurekaClient, err error) {
	ec.eurekaUrl = eurekaUrl
	httpclient := http.Client{Timeout: time.Second * eurekaClientTimeoutInSeconds}
	resp, err := httpclient.Get(eurekaUrl)
	if serr, ok := err.(net.Error); ok && serr.Timeout()  {
		return ec,errEurekaTimesOut
	} else if err != nil {
		return ec, errNoEurekaConnection
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return ec,errEurekaUnexpectedHttpResponseCode
	}
	return ec, nil
}

var errNoIpsFound = errors.New("No IPs associated to the requested App name")

// TODO: Probably we should break this function
func (ec EurekaClient) GetIPs(appName string) ([]string, error) {
	eurekaAppUrl := ec.eurekaUrl + "/v2/apps/" + appName
	// can we reuse the client but starting from 0 in terms of timeout? to review
	httpclient := http.Client{Timeout: time.Second * eurekaClientTimeoutInSeconds}
	req, err := http.NewRequest("GET", eurekaAppUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpclient.Do(req)
	if serr, ok := err.(net.Error); ok && serr.Timeout()  {
		return []string{},errEurekaTimesOut
	} else if err != nil {
		return []string{}, errNoEurekaConnection
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		return []string{},errNoIpsFound
	} else if resp.StatusCode != 200 {
		return []string{},errEurekaUnexpectedHttpResponseCode
	}
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("Response from eureka:", string(body))
	// parsing to get the right data will be done like this: https://stackoverflow.com/a/35665161
	return nil, nil
}
