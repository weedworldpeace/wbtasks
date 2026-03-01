package service

import (
	"app/internal/models"
	"app/pkg/logger"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
)

func getArgs() (models.Arguments, error) {
	args := models.Arguments{}
	pflag.BoolVarP(&args.S, "separated", "s", false, "only strings with delimeter")
	rawD := pflag.StringP("delimeter", "d", "\t", "another delimeter")
	rawF := pflag.StringP("fields", "f", "", "fields to print")

	pflag.Parse()

	if len(*rawD) > 1 {
		return args, models.ErrBadFlagD
	}
	args.D = *rawD

	if len(*rawF) == 0 {
		return args, models.ErrNoFlagF
	}
	fields := strings.Split(*rawF, ",")
	for i := range fields {
		splitted := strings.Split(fields[i], "-")
		if len(splitted) == 1 {
			conv, err := strconv.Atoi(splitted[0])
			if err != nil {
				return args, fmt.Errorf("%v %s", models.ErrInvalidFieldValue, splitted[0])
			}
			if conv < 1 {
				return args, models.ErrFieldsNumberedFrom
			}
			args.F = append(args.F, conv)
		} else if len(splitted) == 2 {
			first, err := strconv.Atoi(splitted[0])
			if err != nil {
				return args, models.ErrInvalidFieldValue
			}
			second, err := strconv.Atoi(splitted[1])
			if err != nil {
				return args, models.ErrInvalidFieldValue
			}
			if second < first {
				return args, models.ErrInvalidDecreasingRange
			}
			if first < 1 {
				return args, models.ErrFieldsNumberedFrom
			}
			for i := first; i <= second; i++ {
				args.F = append(args.F, i)
			}
		} else {
			return args, models.ErrInvalidFieldRange
		}
	}
	slices.Sort(args.F)
	args.F = slices.Compact(args.F)

	return args, nil
}

type Service struct {
	cfg    ServiceConfig
	ctx    context.Context
	doneCh chan struct{}
}

type Wrk struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type ServiceConfig struct {
	Wrks     []Wrk `yaml:"wrks"`
	WrkCount int   `yaml:"wrkCount"`
}

func New(cfg ServiceConfig, ctx context.Context, doneCh chan struct{}) *Service {
	return &Service{cfg, ctx, doneCh}
}

func (s Service) SendData(addr string, ent models.Entity, ch chan models.ChStruct) {
	jsonData, err := json.Marshal(ent)
	if err != nil {
		ch <- models.ChStruct{Err: err}
		return
	}

	resp, err := http.Post(addr, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		ch <- models.ChStruct{Err: err}
		return
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		var errResp models.ErrorResponse

		raw, err := io.ReadAll(resp.Body)
		if err != nil {
			ch <- models.ChStruct{Err: err}
			return
		}
		defer resp.Body.Close()

		err = json.Unmarshal(raw, &errResp)
		if err != nil {
			ch <- models.ChStruct{Err: err}
			return
		}

		ch <- models.ChStruct{Err: fmt.Errorf("%v", errResp.Error)}
	} else {
		var respData models.Response

		raw, err := io.ReadAll(resp.Body)
		if err != nil {
			ch <- models.ChStruct{Err: err}
			return
		}
		defer resp.Body.Close()

		err = json.Unmarshal(raw, &respData)
		if err != nil {
			ch <- models.ChStruct{Err: err}
			return
		}

		ch <- models.ChStruct{Data: respData.Data}
	}
}

func (s *Service) Scale(data []string, args models.Arguments) ([]string, error) {
	lg := logger.LoggerFromCtx(s.ctx).Lg

	count := s.cfg.WrkCount
	if count > len(data) {
		count = len(data)
	}
	dataCh := make(chan models.ChStruct, count)

	dataSize := len(data) - len(data)%count
	for i := range count {
		var ent models.Entity
		ent.Args = args

		if i == count-1 {
			ent.Data = data[i*(dataSize/count):]
		} else {
			ent.Data = data[i*(dataSize/count) : (i+1)*(dataSize/count)]
		}

		go s.SendData(fmt.Sprintf("http://%s:%d/cut", s.cfg.Wrks[i].Host, s.cfg.Wrks[i].Port), ent, dataCh)
	}

	res := make([]string, 0, len(data))
	cter := 0
	requiredCount := count/2 + 1

	for {
		select {
		case <-s.ctx.Done():
			return nil, s.ctx.Err()
		case d := <-dataCh:
			if d.Err != nil {
				lg.Error().Err(d.Err).Send()
				count--
				if count < requiredCount {
					return nil, models.ErrNoQuorum
				}
			} else {
				res = append(res, d.Data...)
				cter++
				if cter >= requiredCount {
					return res, nil
				}
			}
		}
	}
}

func (s *Service) Cut() {
	lg := logger.LoggerFromCtx(s.ctx).Lg

	args, err := getArgs()
	if err != nil {
		lg.Error().Err(err).Send()
		os.Exit(models.ExitErrOccurred)
	}

	data := make([]string, 0)
	sc := bufio.NewScanner(os.Stdin)

	for sc.Scan() {
		data = append(data, sc.Text())
	}

	if sc.Err() != nil {
		lg.Error().Err(err).Send()
		os.Exit(models.ExitErrOccurred)
	}

	res, err := s.Scale(data, args)
	if err != nil {
		lg.Error().Err(err).Send()
		os.Exit(models.ExitErrOccurred)
	}

	for i := range res {
		fmt.Println(res[i])
	}

	close(s.doneCh)
}
