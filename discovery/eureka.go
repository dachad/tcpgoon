package discovery

import (
	"errors"
	"net/http"
)

//TODO: This should become a discovery interface. And Eureka just the first implementation
type EurekaClient struct {
	eurekaUrl string
}

var errNoEurekaConnection = errors.New("Unable to reach Eureka server")
var errEurekaUnexpectedHttpResponseCode = errors.New("Eureka returned a non 200 http response code")

func NewEurekaClient(eurekaUrl string) (ec EurekaClient, err error) {
	ec.eurekaUrl = eurekaUrl
	resp, err := http.Get(eurekaUrl)
	if err != nil {
		return ec, errNoEurekaConnection
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return ec,errEurekaUnexpectedHttpResponseCode
	}
	return ec, nil
}

var errNoIpsFound = errors.New("No IPs associated to the requested App name")

func (ec EurekaClient) GetIPs(appName string) ([]string, error) {
	eurekaAppUrl := ec.eurekaUrl + "/v2/apps/" + appName
	resp, err := http.Get(eurekaAppUrl)
	if err != nil {
		return []string{}, errNoEurekaConnection
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		return []string{},errNoIpsFound
	} else if resp.StatusCode != 200 {
		return []string{},errEurekaUnexpectedHttpResponseCode
	}
	return nil, nil
}
