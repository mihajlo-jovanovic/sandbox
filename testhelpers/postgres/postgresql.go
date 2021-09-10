package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/cenkalti/backoff/v3"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"io/ioutil"
	"net/url"
	"testing"
	"time"
)

const (
	TestUser   = "postgres"
	TestPasswd = "secret"
	DbDriver   = "postgres"
	DbPort     = "5432/tcp"
)

func PrepareTestContainer(t *testing.T, version string) (func(), string) {
	dapi, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Fatal("Could not create docker client")
	}
	ctx := context.Background()
	cfg := &container.Config{
		Hostname: fmt.Sprintf("postgres-%v", uuid.New()),
		Image:    fmt.Sprintf("postgres:%s", version),
		Env:      []string{fmt.Sprintf("POSTGRES_PASSWORD=%s", TestPasswd), "POSTGRES_DB=database"},
	}
	cfg.ExposedPorts = make(map[nat.Port]struct{})
	cfg.ExposedPorts[nat.Port(DbPort)] = struct{}{}
	hostConfig := &container.HostConfig{
		AutoRemove:      false,
		PublishAllPorts: true,
	}
	resp, _ := dapi.ImageCreate(ctx, cfg.Image, types.ImageCreateOptions{})
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

	inspect, err := dapi.ContainerInspect(ctx, container.ID)
	if err != nil {
		_ = dapi.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{})
		t.Fatalf("Could not inspect container: %v", inspect.ID)
	}
	mapped, ok := inspect.NetworkSettings.Ports[nat.Port(DbPort)]
	if !ok || len(mapped) == 0 {
		t.Fatalf("no port mapping found for %s", "5432/tcp")
	}
	addr := fmt.Sprintf("127.0.0.1:%s", mapped[0].HostPort)

	bo := backoff.NewExponentialBackOff()
	bo.MaxInterval = time.Second * 5
	bo.MaxElapsedTime = time.Minute

	var url string
	err = backoff.Retry(func() error {
		u, err := connect(ctx, addr)
		if err != nil {
			return err
		}
		url = u.String()
		return nil
	}, bo)

	if err != nil {
		t.Fatal("Could not test connection to db")
	}
	return cleanup, url
}

func connect(ctx context.Context, addr string) (url.URL, error) {
	u := url.URL{
		Scheme:   DbDriver,
		User:     url.UserPassword(TestUser, TestPasswd),
		Host:     addr,
		Path:     "postgres",
		RawQuery: "sslmode=disable",
	}
	db, err := sql.Open(DbDriver, u.String())
	if err != nil {
		return u, err
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		return u, err
	}
	return u, nil
}
