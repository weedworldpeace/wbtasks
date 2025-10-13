package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/pflag"
)

const (
	exitErrOccured = 1
)

var (
	errHostClosed     = errors.New("host closed connection")
	errBadHostFlag    = errors.New("bad host flag")
	errBadPortFlag    = errors.New("bad port flag")
	errBadTimeoutFlag = errors.New("bad timeout flag")
)

type args struct {
	h string
	p int
	t int
}

func (a *args) valid() error {
	if a.h == "" {
		return errBadHostFlag
	}
	if a.p < 0 || a.p > 65535 {
		return errBadPortFlag
	}
	if a.t < 1 {
		return errBadTimeoutFlag
	}
	return nil
}

func getArgs() (*args, error) {
	data := &args{}

	pflag.StringVarP(&data.h, "host", "h", "", "host bratan")
	pflag.IntVarP(&data.p, "port", "p", 80, "port bratan")
	pflag.IntVarP(&data.t, "timeout", "t", 10, "timeout bratan")

	pflag.Parse()

	err := data.valid()
	return data, err
}

func connect(data *args, scanCh, writeCh chan string, graceCh chan os.Signal) error {
	conn, err := net.DialTimeout("tcp4", fmt.Sprintf("%s:%d", data.h, data.p), time.Second*time.Duration(data.t))
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var resultErr error

	ctxWrite, cancWrite := context.WithCancel(context.Background())
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()

		for {
			select {
			case v, b := <-scanCh:
				if !b {
					return
				}
				_, err := conn.Write([]byte(v))
				if err != nil {
					if err != io.EOF {
						mu.Lock()
						resultErr = multierror.Append(resultErr, errHostClosed)
						mu.Unlock()
					}
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}(ctxWrite)

	wg.Add(1)
	go func() {
		defer wg.Done()

		rd := bufio.NewReader(conn)
		for {
			msg, err := rd.ReadString('\n')
			if err != nil {
				// if err != io.EOF {
				// 	resultErr = multierror.Append(resultErr, err)
				// }
				cancWrite()
				close(writeCh)
				return
			}
			writeCh <- msg
		}
	}()

	go func() {
		<-graceCh
		err := conn.Close()
		if err != nil {
			mu.Lock()
			resultErr = multierror.Append(resultErr, err)
			mu.Unlock()
		}
	}()

	wg.Wait()

	return resultErr
}

func scan(scanCh chan string) {
	sc := bufio.NewScanner(os.Stdin)

	for sc.Scan() {
		scanCh <- fmt.Sprintf("%s\n", sc.Text())
	}

	close(scanCh)
}

func write(writeCh chan string) {
	for v := range writeCh {
		fmt.Printf("host: %s", v)
	}
}

func main() {
	data, err := getArgs()
	if err != nil {
		fmt.Println(err)
		os.Exit(exitErrOccured)
	}

	scanCh := make(chan string)
	go scan(scanCh)

	writeCh := make(chan string)
	go write(writeCh)

	graceCh := make(chan os.Signal, 1)
	signal.Notify(graceCh, syscall.SIGINT, syscall.SIGTERM)

	err = connect(data, scanCh, writeCh, graceCh)
	if err != nil {
		fmt.Println(err)
		os.Exit(exitErrOccured)
	}
}
