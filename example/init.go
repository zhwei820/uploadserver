package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func excuteBash(cmds ...string) error {
	args := []string{}
	if len(cmds) > 1 {
		args = strings.Split(cmds[1], " ")
	}
	cmd := exec.Command(cmds[0], args...)

	fmt.Println("cmd.Args", cmd.Args)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func doReport() {
	fileName := time.Now().In(time.FixedZone("CST", 8*3600)).Format("2006-01-02_15:04:05") + "_80%_cover_report.html"

	cmdStr := fmt.Sprintf(`cover_report.html %s`, fileName)
	excuteBash("cp", cmdStr)

	cmdStr = fmt.Sprintf(`-X POST -F files=@%s -F file_path=zero http://192.168.98.10:8000/index/upload`, fileName)
	excuteBash("curl", cmdStr)
}
