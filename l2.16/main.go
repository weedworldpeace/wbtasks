package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/pflag"
)

const (
	successfulEnd   = "completed successful"
	unsuccessfulEnd = "completed with errors"
	indexFile       = "index.html"
	exitErr         = 1
	hrefAttr        = "href"
)

var (
	errBadRecursiveFlag = errors.New("bad recursive flag")
	errBadArgs          = errors.New("bad args")
)

type args struct {
	r  int
	lf *linkForms
}

type linkForms struct {
	urlAbs    string
	urlNotAbs string
	fromHtml  string
}

func urlAbs(rawUrl string) (string, error) {
	if rawUrl[:4] != "http" {
		rawUrl = "https://" + rawUrl
	}
	if rawUrl[len(rawUrl)-1] == '/' {
		rawUrl = rawUrl[:len(rawUrl)-1]
	}

	return rawUrl, nil
}

func urlNotAbs(urlAbs string) (string, error) {
	idx := 1
	for ; idx < len(urlAbs); idx++ {
		if urlAbs[idx] == '/' && urlAbs[idx-1] == '/' {
			break
		}
	}
	return urlAbs[idx+1:], nil
}

func makeLinkForms(raw string) (*linkForms, error) {
	uAbs, err := urlAbs(raw)
	if err != nil {
		return nil, err
	}

	uNotAbs, err := urlNotAbs(uAbs)
	if err != nil {
		return nil, err
	}

	res := &linkForms{uAbs, uNotAbs, ""}

	return res, nil
}

func (a *args) Valid() error {
	if a.r < 1 {
		return errBadRecursiveFlag
	}
	return nil
}

func parseArgs() (*args, error) {
	data := &args{}
	pflag.IntVarP(&data.r, "recursive", "r", 1, "level of link recursion to download")

	pflag.Parse()
	err := data.Valid()
	if err != nil {
		return data, err
	}

	if len(pflag.Args()) != 1 {
		return data, errBadArgs
	}

	lf, err := makeLinkForms(pflag.Args()[0])
	if err != nil {
		return data, err
	}

	data.lf = lf

	return data, nil
}

func wget(data *args) (string, error) {
	res, err := http.Get(data.lf.urlAbs)
	if err != nil {
		return "", err
	}
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	pathLocDir := fmt.Sprintf("./%s", data.lf.urlNotAbs)
	if idx := strings.Index(pathLocDir, "?"); idx != -1 {
		pathLocDir = pathLocDir[:idx]
	}
	err = os.MkdirAll(pathLocDir, 0766)
	if err != nil {
		return "", err
	}

	pathLocFile := fmt.Sprintf("%s/%s", pathLocDir, indexFile)

	err = os.WriteFile(pathLocFile, raw, 0666)
	if err != nil {
		return "", err
	}
	return pathLocFile, nil
}

func collectUrls(path string, parentLf *linkForms) ([]*linkForms, error) {
	res := make([]*linkForms, 0)

	f, err := os.Open(path)
	if err != nil {
		return res, err
	}

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		return res, err
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr(hrefAttr)
		if !exists || href == "" {
			return
		}
		if strings.HasPrefix(href, "#") {
			return
		}

		parsed, err := url.Parse(href)
		if err != nil {
			return
		}

		if !parsed.IsAbs() {
			curLf, err := makeLinkForms(fmt.Sprintf("%s%s", parentLf.urlAbs, href))
			if err != nil {
				return
			}
			curLf.fromHtml = href
			res = append(res, curLf)
		} else {
			uAbs, err := urlAbs(href)
			if err != nil {
				return
			}
			uNotAbs, err := urlNotAbs(uAbs)
			if err != nil {
				return
			}

			idx := strings.Index(uNotAbs, "/")
			if idx == -1 {
				idx = len(uNotAbs)
			}
			jdx := strings.Index(parentLf.urlNotAbs, "/")
			if jdx == -1 {
				jdx = len(parentLf.urlNotAbs)
			}

			if uNotAbs[:idx] == parentLf.urlNotAbs[:jdx] {
				curLf, err := makeLinkForms(href)
				if err != nil {
					return
				}
				curLf.fromHtml = href
				res = append(res, curLf)
			}
		}
	})
	return res, nil
}

func replaceUrl(path string, hm map[string]string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		return err
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr(hrefAttr)
		if !exists {
			return
		}
		if v, b := hm[href]; b {
			s.RemoveAttr(hrefAttr)
			s.SetAttr(hrefAttr, fmt.Sprintf("../%s", v))
		}
	})

	modified, err := doc.Html()
	if err != nil {
		return err
	}

	return os.WriteFile(path, []byte(modified), 0666)
}

func wgetWrap(data *args, errCh chan error, workDir string, hm map[string]string) {
	path, err := wget(data)
	if err != nil {
		errCh <- err
		return
	}
	hm[data.lf.fromHtml] = path
	if data.r > 1 {
		urlsToWget, err := collectUrls(path, data.lf)
		if err != nil {
			errCh <- err
		} else {
			for _, v := range urlsToWget {
				if _, b := hm[v.fromHtml]; !b {
					wgetWrap(&args{data.r - 1, v}, errCh, workDir, hm)
				}
			}
			err = replaceUrl(path, hm)
			if err != nil {
				errCh <- err
			}
		}
	}
}

func main() {
	data, err := parseArgs()
	if err != nil {
		fmt.Println(err)
		return
	}
	data.lf.fromHtml = data.lf.urlAbs

	exitFlag := true
	errCh := make(chan error)

	os.Mkdir(fmt.Sprintf("./%s", data.lf.urlNotAbs), 0766)

	go func() {
		fmt.Println("processing...")
		wgetWrap(data, errCh, data.lf.urlNotAbs, make(map[string]string))
		close(errCh)
	}()

	for i := range errCh {
		fmt.Println(i)
		exitFlag = false
	}
	if exitFlag {
		fmt.Println(successfulEnd)
	} else {
		fmt.Println(unsuccessfulEnd)
		os.Exit(exitErr)
	}
}
