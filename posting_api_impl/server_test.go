package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/url"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v3"
	"github.com/docker/distribution/uuid"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func TestMyFirstTest(t *testing.T) {
	dapi, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Fatal("Could not create docker client")
	}
	ctx := context.Background()
	images, err := dapi.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		t.Fatalf("Could not list images: %v", err)
	}
	for _, i := range images {
		fmt.Printf("Image ID: %s\n", i.ID)
	}
	uuid.Generate()
	cfg := &container.Config{
		Hostname: fmt.Sprintf("postgres-%v", uuid.Generate()),
		Image:    "postgres:11",
		Env:      []string{"POSTGRES_PASSWORD=secret", "POSTGRES_DB=database"},
	}
	cfg.ExposedPorts = make(map[nat.Port]struct{})
	cfg.ExposedPorts[nat.Port("5432/tcp")] = struct{}{}
	hostConfig := &container.HostConfig{
		AutoRemove:      false,
		PublishAllPorts: true,
	}
	resp, _ := dapi.ImageCreate(ctx, "postgres:11", types.ImageCreateOptions{})
	if resp != nil {
		_, _ = ioutil.ReadAll(resp)
	}
	netConfig := &network.NetworkingConfig{}
	container, err := dapi.ContainerCreate(ctx, cfg, hostConfig, netConfig, cfg.Hostname)
	if err != nil {
		t.Fatalf("Could not start container: %s", cfg.Hostname)
	}
	err = dapi.ContainerStart(ctx, container.ID, types.ContainerStartOptions{})
	if err != nil {
		t.Fatalf("Could not start container: %s", cfg.Hostname)
	}
	cleanup := func() {
		for i := 0; i < 10; i++ {
			err := dapi.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{Force: true})
			if err == nil {
				return
			}
			time.Sleep(1 * time.Second)
		}
	}
	defer cleanup()

	inspect, err := dapi.ContainerInspect(ctx, container.ID)
	if err != nil {
		_ = dapi.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{})
		t.Fatalf("Could not inspect container: %v", inspect.ID)
	}
	mapped, ok := inspect.NetworkSettings.Ports[nat.Port("5432/tcp")]
	if !ok || len(mapped) == 0 {
		t.Fatalf("no port mapping found for %s", "5432/tcp")
	}
	addr := fmt.Sprintf("127.0.0.1:%s", mapped[0].HostPort)

	bo := backoff.NewExponentialBackOff()
	bo.MaxInterval = time.Second * 5
	bo.MaxElapsedTime = time.Minute

	err = backoff.Retry(func() error {
		err := connect(ctx, addr)
		if err != nil {
			return err
		}
		return nil
	}, bo)

	if err != nil {
		t.Fatal("Could not test connection to db")
	}
}

func connect(ctx context.Context, addr string) error {
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword("postgres", "secret"),
		Host:     addr,
		Path:     "postgres",
		RawQuery: "sslmode=disable",
	}

	db, err := sql.Open("postgres", u.String())
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}
	fmt.Println("Connected!")
	return nil
}
