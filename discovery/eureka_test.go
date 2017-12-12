package discovery

import (
	"gopkg.in/ory-am/dockertest.v3"
	dc "github.com/fsouza/go-dockerclient"
	"github.com/jaume-pinyol/fargo"
	"log"
	"os"
	"testing"
	"strconv"
	"github.com/op/go-logging"
)

var eurekaTestPort int = 8080
var eurekaTestUrl string = "http://127.0.0.1:" + strconv.Itoa(eurekaTestPort) + "/eureka"

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository:   "netflixoss/eureka",
		Tag:          "1.3.1",
		PortBindings: map[dc.Port][]dc.PortBinding{
			dc.Port(strconv.Itoa(eurekaTestPort) + "/tcp"): {{HostIP: "", HostPort: strconv.Itoa(eurekaTestPort)}},
		},
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		_, err := NewEurekaClient(eurekaTestUrl)
		return err
	}); err != nil {
		log.Fatalf("Could not connect to the docker resource: %s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestEurekaClientNoEureka(t *testing.T) {
	_, err := NewEurekaClient("http://localhost:9999/thisshouldntwork")
	if err != errNoEurekaConnection {
		t.Fatal("We shouldnt reach eureka if Eureka hostname/port is completely wrong")
	}
}

func TestEurekaClientWrongEurekaContext(t *testing.T) {
	_, err := NewEurekaClient(eurekaTestUrl + "badsuffix")
	if err != errEurekaUnexpectedHttpResponseCode {
		t.Fatal("Eureka should be reachable but, when asking a wrong URL, it should return a non 200 response code")
	}
}

func TestEurekaClientUnknownApp(t *testing.T) {
	appName := "unknown"
	eurekaClient, err := NewEurekaClient(eurekaTestUrl)
	if err != nil {
		t.Fatal("We cannot connect to the specified eureka server:", err)
	}
	t.Log("Connection to Eureka established")
	_, err = eurekaClient.GetIPs(appName)
	if err != errNoIpsFound {
		t.Fatal("Eureka did return something different from an No-IPs-error associated to the unknown App")
	}
}

func TestEurekaClientValidApp(t *testing.T) {
	appName := "testApp"
	ipAddr := "192.0.2.1"
	port := 10080
	registerDummyAppInTestEureka(appName, ipAddr, port)
	eurekaClient, err := NewEurekaClient(eurekaTestUrl)
	if err != nil {
		t.Fatal("We cannot connect to the specified eureka server:", err)
	}
	t.Log("Connection to Eureka established")
	ipsFromEureka, err := eurekaClient.GetIPs(appName)
	if err != nil {
		t.Fatal("Eureka returned an error when requesting the IPs:", err)
	}
	if len(ipsFromEureka) != 1 || ipsFromEureka[0] != ipAddr {
		t.Fatal("Eureka returned a set of IPs we did not expect for our service:", ipsFromEureka )
	}

}
func registerDummyAppInTestEureka(appName string, ipAddr string, port int) {
	logging.SetLevel(logging.ERROR, "fargo")
	fargoclient := fargo.NewConn(eurekaTestUrl + "/v2")
	appInstance := &fargo.Instance{
		HostName:         "dummyhost",
		Port:             port,
		SecurePort:       port,
		App:              appName,
		IPAddr:           ipAddr,
		VipAddress:       ipAddr,
		SecureVipAddress: ipAddr,
		DataCenterInfo:   fargo.DataCenterInfo{Name: fargo.MyOwn},
		Status:           fargo.UP,
		Overriddenstatus: fargo.UNKNOWN,
		HealthCheckUrl:   "http://" + ipAddr + ":" + "8080" + "/healthcheck",
		StatusPageUrl:    "http://" + ipAddr + ":" + "8080" + "/healthcheck",
		HomePageUrl:      "http://" + ipAddr + ":" + "8080" + "/",
		AsgName:          "dummyAsg",
	}
	fargoclient.RegisterInstance(appInstance)
}
