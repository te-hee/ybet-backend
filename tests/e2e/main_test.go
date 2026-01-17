package e2e

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go/modules/compose"
)

var (
	GatewayURL = "http://localhost:8080"
	WSURL      = "ws://localhost:8081/ws"
)

func TestMain(m *testing.M) {
	composeFilePaths := []string{"../../docker-compose.yml"}

	stack, err := compose.NewDockerCompose(composeFilePaths...)
	if err != nil {
		fmt.Printf("failed init docker compose QwQ: %v\n", err)
		os.Exit(1)
	}

	defer func() {
		fmt.Println("cleaning docker compose...")
		if err := stack.Down(context.Background(), compose.RemoveOrphans(true), compose.RemoveImagesLocal); err != nil {
			fmt.Printf("error cleaning :c : %v\n", err)
		}
	}()

	fmt.Println("starting docker compose...")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if err := stack.Up(ctx, compose.Wait(true)); err != nil {
		fmt.Printf("error starting docker compose >~<: %v\n", err)
		os.Exit(1)
	}

	if err := waitForGateway(GatewayURL); err != nil {
		fmt.Printf("gateway did not start in time TwT: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()

	os.Exit(code)
}

func waitForGateway(url string) error {
	for range 30 {
		if resp, _ := http.Get(url + "/health"); resp.StatusCode == http.StatusOK {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("gateway did not start :c")
}
