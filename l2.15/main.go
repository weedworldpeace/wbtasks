package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

var (
	errBadArg = errors.New("bad argument")
)

func cdWrap(path string) error {
	err := os.Chdir(path)
	return err
}

func pwdWrap() (string, error) {
	return os.Getwd()
}

func echoWrap(s string) string {
	return os.ExpandEnv(s)
}

func killWrap(pid string) error {
	pidConv, err := strconv.Atoi(pid)
	if err != nil {
		return err
	}
	p, err := os.FindProcess(pidConv)
	if err != nil {
		return err
	}
	return p.Kill()
}

func psWrap() ([]byte, error) {
	cmd := exec.Command("bash", "-c", "ps aux")
	return cmd.Output()
}

func execWrap(ctx context.Context, cmdSplitted []string, inp string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "bash", "-c", strings.Join(cmdSplitted, " "))
	if len(inp) != 0 {
		writer, err := cmd.StdinPipe()
		if err != nil {
			return []byte{}, err
		}
		writer.Write([]byte(inp))
		writer.Close()
	}
	return cmd.CombinedOutput()
}

func runCmd(row string, ctx context.Context, callbackCh chan bool) {
	cmds := strings.Split(row, "|")
	lastOut := ""
	for idx := 0; idx < len(cmds); idx++ {
		cmdSplitted := strings.Fields(cmds[idx])
		if len(cmdSplitted) < 1 {
			close(callbackCh)
			return
		}
		switch cmdSplitted[0] {
		case "cd":
			if len(cmdSplitted) != 2 {
				fmt.Println(errBadArg)
				close(callbackCh)
				return
			}
			if len(cmds) == 1 {
				err := cdWrap(cmdSplitted[1])
				if err != nil {
					fmt.Println(err)
					close(callbackCh)
					return
				}
			}
			lastOut = ""
		case "pwd":
			res, err := pwdWrap()
			if err != nil {
				fmt.Println(err)
				close(callbackCh)
				return
			}
			lastOut = fmt.Sprintf("%s\n", res)
		case "echo":
			lastOut = fmt.Sprintf("%s\n", echoWrap(strings.Join(cmdSplitted[1:], " ")))
		case "kill":
			if len(cmdSplitted) != 2 {
				fmt.Println(errBadArg)
				close(callbackCh)
				return
			}
			err := killWrap(cmdSplitted[1])
			if err != nil {
				fmt.Println(err)
				close(callbackCh)
				return
			}
			lastOut = ""
		case "ps":
			res, err := psWrap()
			if err != nil {
				fmt.Println(err)
				close(callbackCh)
				return
			}
			lastOut = fmt.Sprintf("%s", res)
		default:
			out, err := execWrap(ctx, cmdSplitted, lastOut)
			if err != nil {
				close(callbackCh)
				return
			}
			lastOut = fmt.Sprintf("%s", string(out))
		}
	}
	fmt.Print(lastOut)
	close(callbackCh)
}

func startBash() {
	sc := bufio.NewScanner(os.Stdin)
	stopCh := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGINT, syscall.SIGTERM)
	for sc.Scan() {
		ctx, canc := context.WithCancel(context.Background())
		signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)
		callbackCh := make(chan bool)
		go runCmd(sc.Text(), ctx, callbackCh)
		select {
		case <-stopCh:
		case <-callbackCh:
		}
		canc()
		signal.Reset(syscall.SIGINT, syscall.SIGTERM)
	}
}

func main() {
	startBash()
}
