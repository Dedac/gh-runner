package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
)

func runStart(name string) {
	var runner string
	if runtime.GOOS != "windows" {
		runner = fmt.Sprintf("%s/run.sh", name)
	} else {
		runner = fmt.Sprintf(".\\%s\\run.cmd", name)
	}
	runcmd := exec.Command(runner)
	runcmd.Stdout = os.Stdout
	runcmd.Stderr = os.Stderr
	err := runcmd.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func runStop(name string) {
	if runtime.GOOS != "windows" {

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
	} else {
		//TODO: Implement stop for windows
	}
}

func createService(name string) {
	err := os.Chdir(name)
	if err != nil {
		log.Fatal(err)
	}

	//Create a service to run the runner
	var runcmd *exec.Cmd
	if runtime.GOOS == "darwin" {
		runcmd = exec.Command("./svc.sh", "install")
	} else if runtime.GOOS == "linux" {
		runcmd = exec.Command("sudo", "./svc.sh", "install")
	} else if runtime.GOOS == "windows" {
		log.Fatal("On windows, you must configure the service when creating the runner.")
	} else {
		log.Fatal("Unsupported OS")
	}

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
	var runcmd *exec.Cmd
	if runtime.GOOS == "darwin" {
		runcmd = exec.Command("./svc.sh", "start")
	} else if runtime.GOOS == "linux" {
		runcmd = exec.Command("sudo", "./svc.sh", "start")
	} else if runtime.GOOS == "windows" {
		runcmd = exec.Command("powershell", "Start-Service", "actions.runner.*")
	} else {
		log.Fatal("Unsupported OS")
	}
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
	var runcmd *exec.Cmd
	if runtime.GOOS == "darwin" {
		runcmd = exec.Command("./svc.sh", "stop")
	} else if runtime.GOOS == "linux" {
		runcmd = exec.Command("sudo", "./svc.sh", "stop")
	} else if runtime.GOOS == "windows" {
		runcmd = exec.Command("powershell", "Stop-Service", "actions.runner.*")
	} else {
		log.Fatal("Unsupported OS")
	}
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
	var runcmd *exec.Cmd
	if runtime.GOOS == "darwin" {
		runcmd = exec.Command("./svc.sh", "uninstall")
	} else if runtime.GOOS == "linux" {
		runcmd = exec.Command("sudo", "./svc.sh", "uninstall")
	} else if runtime.GOOS == "windows" {
		runcmd = exec.Command("powershell", "Remove-Service", "actions.runner.*")
	} else {
		log.Fatal("Unsupported OS")
	}
	runcmd.Stdout = os.Stdout
	runcmd.Stderr = os.Stderr
	err = runcmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
