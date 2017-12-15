package discovery

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

//TODO: This should become a discovery interface. And Eureka just the first implementation
type EurekaClient struct {
	eurekaURL string
}

//TODO: Creating our own error type and wrapping standard net/http errors could be useful to prevent
// the original errors from being lost
var errEurekaTimesOut = errors.New("Eureka server timed out")
var errNoEurekaConnection = errors.New("Unable to reach Eureka server")
var errEurekaUnexpectedHTTPResponseCode = errors.New("Eureka returned a non 200 http response code")

const eurekaClientTimeoutInSeconds = 10

func NewEurekaClient(eurekaURL string) (ec EurekaClient, err error) {
	ec.eurekaURL = eurekaURL
	httpclient := http.Client{Timeout: time.Second * eurekaClientTimeoutInSeconds}
	resp, err := httpclient.Get(eurekaURL)
	if serr, ok := err.(net.Error); ok && serr.Timeout() {
		return ec, errEurekaTimesOut
	} else if err != nil {
		return ec, errNoEurekaConnection
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return ec, errEurekaUnexpectedHTTPResponseCode
	}
	return ec, nil
}

var errNoIpsFound = errors.New("No IPs associated to the requested App name")

// TODO: Probably we should break this function
func (ec EurekaClient) GetIPs(appName string) ([]string, error) {
	eurekaAppURL := ec.eurekaURL + "/v2/apps/" + appName
	// can we reuse the client but starting from 0 in terms of timeout? to review
	httpclient := http.Client{Timeout: time.Second * eurekaClientTimeoutInSeconds}
	req, err := http.NewRequest("GET", eurekaAppURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpclient.Do(req)
	if serr, ok := err.(net.Error); ok && serr.Timeout() {
		return []string{}, errEurekaTimesOut
	} else if err != nil {
		return []string{}, errNoEurekaConnection
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		return []string{}, errNoIpsFound
	} else if resp.StatusCode != 200 {
		return []string{}, errEurekaUnexpectedHTTPResponseCode
	}
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("Response from eureka:", string(body))
	// parsing to get the right data will be done like this: https://stackoverflow.com/a/35665161
	return nil, nil
}
