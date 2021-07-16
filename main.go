package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

func requestString(URL string) (string, error) {
	buff := bytes.Buffer{}
	// json.NewEncoder(&buff).Encode()
	req, err := http.NewRequest(http.MethodGet, URL, &buff)
	if err != nil {
		log.Println(err)
		return "", err
	}
	// req.Header.Set("Content-Type", "application/json")
	// req.AddCookie(&http.Cookie{Name: "authsession", Value: ""})
	req.Close = true
	// выполняем запрос
	response, err := (&http.Client{
		Timeout: time.Minute,
		// Transport: &http.Transport{
		// 	TLSClientConfig: &tls.Config{
		// 		InsecureSkipVerify: true,
		// 	},
		// },
	}).Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode > 299 {
		err = fmt.Errorf("ошибка %s", response.Status)
		log.Println(err)
		return "", err
	}
	buff.Reset()
	_, err = buff.ReadFrom(response.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return buff.String(), nil
}

func requestBytes(URL string) ([]byte, error) {
	buff := bytes.Buffer{}
	req, err := http.NewRequest(http.MethodGet, URL, &buff)
	if err != nil {
		log.Println(URL, err)
		return []byte{}, err
	}
	req.Close = true
	// выполняем запрос
	response, err := (&http.Client{
		Timeout: time.Minute,
		// Transport: &http.Transport{
		// 	TLSClientConfig: &tls.Config{
		// 		InsecureSkipVerify: true,
		// 	},
		// },
	}).Do(req)
	if err != nil {
		log.Println(URL, err)
		return []byte{}, err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode > 299 {
		err = fmt.Errorf("ошибка %s", response.Status)
		log.Println(URL, err)
		return []byte{}, err
	}
	buff.Reset()
	_, err = buff.ReadFrom(response.Body)
	if err != nil {
		log.Println(URL, err)
		return []byte{}, err
	}
	return buff.Bytes(), nil
}

func stringToUrlAndNumPages(s string) (string, string, int) {
	rp := `(?sm)` + regexp.QuoteMeta(`<div class="row gallery_first">
	    <div class="col-md-4 col left_cover">
	      <a href="`) + `.+?` + regexp.QuoteMeta(`"><img src="`) + `(.+?)\".+?` +
		regexp.QuoteMeta(`<h1>`) + `(.+?)` + regexp.QuoteMeta(`</h1>`) + `.*?` +
		regexp.QuoteMeta(`<li class="pages">Pages: `) + `(\d+).*?` + regexp.QuoteMeta(`</li>`)
	res := regexp.MustCompile(rp).FindAllStringSubmatch(s, -1)
	if len(res) < 1 || len(res[0]) != 4 {
		return "", "", 0
	}
	u, n, p := res[0][1:][0], res[0][1:][1], res[0][1:][2]
	uTmp := strings.Split(u, "/")
	if len(uTmp) < 2 {
		return "", "", 0
	}
	u = strings.Join(uTmp[:len(uTmp)-1], "/")
	pTmp, err := strconv.Atoi(p)
	if err != nil {
		return "", "", 0
	}
	return u, n, pTmp
}

func escapeFileName(n string) string {
	const r = " "
	if len([]rune(n)) > 200 {
		n = string([]rune(n)[:200])
	}
	return strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ReplaceAll(
					strings.ReplaceAll(
						strings.ReplaceAll(
							strings.ReplaceAll(
								strings.ReplaceAll(
									strings.ReplaceAll(
										n, `\`, r,
									), `/`, r,
								), `|`, r,
							), `:`, r,
						), `"`, r,
					), `*`, r,
				), `?`, r,
			), `<`, r,
		), `>`, r,
	)
}

func load(URL string) bool {
	log.Printf("начато %s\n", URL)
	rawData, err := requestString(URL)
	if err != nil {
		return false
	}
	u, n, p := stringToUrlAndNumPages(rawData)
	if p < 1 {
		return false
	}
	hasErr := false
	buff := &bytes.Buffer{}
	zw := zip.NewWriter(buff)
	for i := 1; i <= p; i++ {
		tp := "jpg"
		data, err := requestBytes(fmt.Sprintf("%s/%d.%s", u, i, tp))
		if err != nil {
			tp = "png"
			data, err = requestBytes(fmt.Sprintf("%s/%d.%s", u, i, tp))
			if err != nil {
				tp = "gif"
				data, err = requestBytes(fmt.Sprintf("%s/%d.%s", u, i, tp))
				if err != nil {
					log.Printf("ошибка %s %s\n", URL, err.Error())
					hasErr = true
					break
				}
			}
		}
		w, err := zw.Create(fmt.Sprintf("%d.%s", i, tp))
		if err != nil {
			log.Printf("ошибка %s %s\n", URL, err.Error())
			hasErr = true
			break
		}
		w.Write(data)
	}
	w, err := zw.Create("info.txt")
	if err != nil {
		log.Printf("ошибка %s %s\n", URL, err.Error())
		hasErr = true
	}
	fmt.Fprintln(w, "URL:", URL)
	fmt.Fprintln(w, "name:", n)
	fmt.Fprintln(w, "pages:", p)
	zw.Close()
	if hasErr {
		log.Printf("завершено %s %s\n", URL, n)
	} else {
		f, err := os.Create(fmt.Sprintf("loads/%s.zip", escapeFileName(n)))
		if err != nil {
			log.Printf("ошибка %s %s\n", URL, err.Error())
			return false
		}
		_, err = buff.WriteTo(f)
		if err != nil {
			log.Printf("ошибка %s %s\n", URL, err.Error())
			return false
		}
		f.Close()
		log.Printf("успешно %s %s\n", URL, n)
	}
	return !hasErr
}

func main() {
	lf, err := os.Create("log.txt")
	if err != nil {
		log.Println(err)
		return
	}
	log.SetOutput(io.MultiWriter(os.Stderr, lf))
	URL := flag.String("url", "", "Адрес для закачки")
	flag.Parse()
	_, err = os.Stat("loads")
	if os.IsNotExist(err) {
		os.MkdirAll("loads", 0777)
	}
	if *URL != "" {
		load(*URL)
		return
	}
	f, err := os.Open("task.txt")
	if err != nil {
		log.Println(err)
		return
	}
	sc := bufio.NewScanner(f)
	wg := &sync.WaitGroup{}
	notComplete := make([]string, 0)
	for sc.Scan() {
		if sc.Text() == "" {
			continue
		}
		wg.Add(1)
		go func(u string) {
			if !load(u) {
				notComplete = append(notComplete, u)
			}
			wg.Done()
		}(sc.Text())
	}
	wg.Wait()
	f.Close()
	f, err = os.Create("task.txt")
	if err != nil {
		log.Println(err)
		return
	}
	for _, u := range notComplete {
		fmt.Fprintln(f, u)
	}
	f.Close()
	time.Sleep(time.Second)
	// load("https://imhentai.xxx/gallery/686547/")
}
