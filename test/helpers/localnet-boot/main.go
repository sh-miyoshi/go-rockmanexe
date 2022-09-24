package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

func build(wg *sync.WaitGroup, dir string) {
	defer wg.Done()

	var stderr bytes.Buffer
	cmd := exec.Command("go", "build")
	if dir != "." {
		cmd.Dir = dir
	}
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("go build in %s error: %v\n", dir, err)
		fmt.Println(stderr.String())
		os.Exit(1)
	}
}

func main() {
	os.Chdir("../../../")

	fmt.Println("building binaries ...")
	var wg sync.WaitGroup
	wg.Add(4)
	go build(&wg, ".")
	go build(&wg, "cmd/fakeserver")
	go build(&wg, "cmd/router")
	go build(&wg, "cmd/botclient")
	wg.Wait()
	fmt.Println("done")

	fmt.Println("Run fakeserver")
	var fakeserverStderr bytes.Buffer
	fakeserverCmd := exec.Command("./fakeserver.exe", "--config", "config.yaml")
	fakeserverCmd.Dir = "cmd/fakeserver"
	fakeserverCmd.Stderr = &fakeserverStderr
	if err := fakeserverCmd.Start(); err != nil {
		fmt.Printf("Failed to run fakeserver: %v\n", err)
		return
	}
	defer fakeserverCmd.Process.Kill()

	fmt.Println("Run router")
	var routerStderr bytes.Buffer
	routerCmd := exec.Command("./router.exe", "--config", "config.yaml")
	routerCmd.Dir = "cmd/router"
	routerCmd.Stderr = &routerStderr
	if err := routerCmd.Start(); err != nil {
		fmt.Printf("Failed to run router: %v\n", err)
		return
	}
	defer routerCmd.Process.Kill()

	// Waiting router wakeup
	time.Sleep(1 * time.Second)

	fmt.Println("Run botclient")
	var clientStderr bytes.Buffer
	clientCmd := exec.Command("./botclient.exe", "-c", "tester2", "-log", "botclient.log")
	clientCmd.Dir = "cmd/botclient"
	clientCmd.Stderr = &clientStderr
	if err := clientCmd.Start(); err != nil {
		fmt.Printf("Failed to run botclient: %v\n", err)
		return
	}
	defer clientCmd.Process.Kill()

	fmt.Println("Run main app")
	var appStderr bytes.Buffer
	appCmd := exec.Command("go-rockmanexe.exe", "--config", "data/config_debug.yaml")
	appCmd.Stderr = &appStderr
	if err := appCmd.Run(); err != nil {
		fmt.Printf("Failed to run main app: %v\n", err)
	}

	fmt.Printf("Router Stderr: %s\n", routerStderr.String())
	fmt.Printf("Test Client Stderr: %s\n", clientStderr.String())
	fmt.Printf("App Stderr: %s\n", appStderr.String())
}
