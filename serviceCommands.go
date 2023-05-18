package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func runStart(name string) {
	runner := fmt.Sprintf("%s/run.sh", name)
	runcmd := exec.Command(runner)
	runcmd.Stdout = os.Stdout
	runcmd.Stderr = os.Stderr
	err := runcmd.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func runStop(name string) {
	//kill the 3 created processes
	runnerprocs := fmt.Sprintf("%[1]s/run.sh|%[1]s/bin/Runner.Listener|%[1]s/run-helper.sh", name)
	//Find the pid of the runner
	c1 := exec.Command("pgrep", "-f", runnerprocs)
	//kill the processes
	c2 := exec.Command("xargs", "kill")
	c2.Stdin, _ = c1.StdoutPipe()
	c2.Stdout = os.Stdout
	_ = c2.Start()
	err := c1.Run()
	err2 := c2.Wait()

	if err != nil {
		log.Fatal(err)
	}
	if err2 != nil {
		log.Fatal(err)
	}
}

func createService(name string) {
	err := os.Chdir(name)
	if err != nil {
		log.Fatal(err)
	}

	//Create a service to run the runner
	runcmd := exec.Command("./svc.sh", "install")
	runcmd.Stdout = os.Stdout
	runcmd.Stderr = os.Stderr
	err = runcmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	os.Chdir("..")
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
	err = runcmd.Run()
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
	err = runcmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func removeService(name string) {
	err := os.Chdir(name)
	if err != nil {
		log.Fatal(err)
	}
	//Remove the service
	runcmd := exec.Command("./svc.sh", "uninstall")
	runcmd.Stdout = os.Stdout
	runcmd.Stderr = os.Stderr
	err = runcmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
