package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func runStart(name string) {
	err := os.Chdir(name)
	if err != nil {
		log.Fatal(err)
	}

	runcmd := exec.Command("./run.sh")
	runcmd.Stdout = os.Stdout
	runcmd.Stderr = os.Stderr
	err = runcmd.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func runStop(name string) {
	//Find the pid of the runner
	pocessname := fmt.Sprintf("%s/run.sh", name)
	pidcmd := exec.Command("pgrep", "-f", pocessname)
	err := pidcmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	//Kill the pid
	killcmd := exec.Command("xargs", "kill", "-SIGINT")
	err = killcmd.Run()

	if err != nil {
		log.Fatal(err)
	}
}

func createService(name string) {
	err := os.Chdir(name)
	if err != nil {
		log.Fatal(err)
	}

	//Create a service to run the runner
	runcmd := exec.Command("./run.sh")
	runcmd.Stdout = os.Stdout
	runcmd.Stderr = os.Stderr
	err = runcmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	runService(name)
}

func runService(name string) {
	err := os.Chdir(name)
	if err != nil {
		log.Fatal(err)
	}

	//run the configured service
	runcmd := exec.Command("./svc.sh", "start")
	runcmd.Stdout = os.Stdout
	runcmd.Stderr = os.Stderr
	err = runcmd.Start()
	if err != nil {
		log.Fatal(err)
	}

}
func stopService(name string) {
	err := os.Chdir(name)
	if err != nil {
		log.Fatal(err)
	}

	//Stop the running runner service
	runcmd := exec.Command("./svc.sh", "stop")
	runcmd.Stdout = os.Stdout
	runcmd.Stderr = os.Stderr
	err = runcmd.Start()
	if err != nil {
		log.Fatal(err)
	}
}
