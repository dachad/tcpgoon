package docker

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net"

	//"os"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	"github.com/dachad/tcpgoon/debugging"
)

type IPandPort struct {
	IP   string
	Port int
}

func lookForAvailableLocalPort() (availablePort int){
    minPort := 1025
	maxPort := 65535

	for {
		availablePort = rand.Intn(maxPort - minPort) + minPort
		ln, err := net.Listen("tcp", ":" + strconv.Itoa(availablePort))
		ln.Close()
		if err == nil {
			break
		}
	}
	return availablePort
}

func DownloadAndRun(image string, port int) (targetinfo IPandPort, containerID string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(debugging.DebugOut,"Pulling Docker image", image)

	out, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(debugging.DebugOut, out)

	defer out.Close()

	fmt.Fprintln(debugging.DebugOut, "Docker image", image, "downloaded")

	mappedPort := ""
	mapperProto := ""
	if port == 0 {
		// if not specific mapping defined we look for exposed ports
		// getting the exposed port before decide the mapping

		// TODO support mapping volumes, some containers required them
		respCreate, err := cli.ContainerCreate(ctx, &container.Config{
			Image: image,
		}, nil, nil, "")
		if err != nil {
			panic(err)
		}

		if err := cli.ContainerStart(ctx, respCreate.ID, types.ContainerStartOptions{}); err != nil {
			panic(err)
		}

		respInspect := inspect(respCreate.ID)
		natPorts := respInspect.NetworkSettings.NetworkSettingsBase.Ports

		for k, _ := range natPorts {
			mappedPort = k.Port()
			mapperProto = k.Proto()
			break
		}
		// Stopping container used to get the port mapping
		Stop(respCreate.ID)
	} else {
		mappedPort = strconv.Itoa(port)
		mapperProto = "tcp"
	}

	if mappedPort == "" || mapperProto != "tcp" {
		panic("Not found any TCP port mapping")
	} else {
		fmt.Println("Internal Docker port binding:", mapperProto, mappedPort)
	}

	// hack https://stackoverflow.com/questions/47395973/issue-while-using-docker-api-for-go-cannot-import-nat

	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: strconv.Itoa(lookForAvailableLocalPort()),
	}
	containerPort, err := nat.NewPort(mapperProto, mappedPort)
	if err != nil {
		panic("Unable to get the port")
	}

	portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}
	resp, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: image,
		},
		&container.HostConfig{
			PortBindings: portBinding,
		}, nil, "")
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(debugging.DebugOut, "Docker container built")

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	fmt.Fprintln(debugging.DebugOut, "Started docker container", resp.ID)

	targetinfo.IP = "127.0.0.1"
	targetinfo.Port, _ = strconv.Atoi(hostBinding.HostPort)

	return targetinfo, resp.ID
}

func Stop(containerID string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(debugging.DebugOut, "Stopping container", containerID, "...")
	if err := cli.ContainerStop(ctx, containerID, nil); err != nil {
		panic(err)
	}
	fmt.Fprintln(debugging.DebugOut, "Container stopped successfully")
}

func inspect(containerID string) types.ContainerJSON {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(debugging.DebugOut, "Getting container info", containerID, "...")
	resp, err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		panic(err)
	}

	return resp
}
