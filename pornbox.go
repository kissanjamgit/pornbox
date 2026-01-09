// Package pornbox to scrape resource from pornbox.com
package pornbox

import "C"

// test clone ssh

import (
	"fmt"
	"regexp"

	"github.com/tidwall/gjson"
	"resty.dev/v3"
)

var header = map[string]string{
	"User-Agent":       "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:143.0) Gecko/20100101 Firefox/143.0",
	"Accept":           "*/*",
	"Accept-Language":  "en-US,en;q=0.5",
	"Accept-Encoding":  "gzip, deflate, br, zstd",
	"X-CSRF-Token":     "F6wyQTnF-rZeY2u7ZV8fUaY0dslunx1GQySQ",
	"X-Requested-With": "XMLHttpRequest",
	"Sec-GPC":          "1",
	"Connection":       "keep-alive",
	"Cookie":           "JDIALOG3=FIZT353D78913CYVAYLIVSOV4SGQM6MG516ZCS3SA5QEMVHWOA; http_referer=; entry_point=https%3A%2F%2Fpornbox.com%2Fapplication%2Fwatch-page%2F3842321; boxsessid=s%3APQIiI9skPgvgznr2YhxuACHu2WwahlCg.r1a4qD658kvqUWwtFlrGmTwWss7UvuwuR7%2BNUtkdfRo; version_website_id=j%3A%5B25%5D; agree18=1",
	"Sec-Fetch-Dest":   "empty",
	"Sec-Fetch-Mode":   "cors",
	"Sec-Fetch-Site":   "same-origin",
}

type ContentResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PornBox struct {
	url    string
	header map[string]string
}

func New(url string) PornBox {
	return PornBox{url, header}
}

func title(json string) (name string, err error) {
	value := gjson.Get(json, "scene_name").Value()
	var ok bool
	if name, ok = value.(string); !ok {
		err = fmt.Errorf("PornBox: (_name) `value.(string)` value: %v", value)
		return
	}
	return
}

func video(input string, client *resty.Client) (cr ContentResource, err error) {
	re := regexp.MustCompile(`\d+`)
	id := re.FindString(input)
	url := fmt.Sprintf("https://pornbox.com/contents/%v", id)
	client.SetHeaders(header)
	var res *resty.Response
	res, err = client.R().Get(url)
	if err != nil {
		return
	}
	value := gjson.Get(res.String(), "medias.@reverse.#(title==Trailer).media_id")
	mediaID, ok := value.Value().(float64)
	if !ok {
		err = fmt.Errorf("PornBox: (video) `value.(float64)` value: %v ", value)
		return
	}
	// fmt.Printf("mediaID: %d\n", int(mediaID))
	var getSrc *resty.Response
	getSrc, err = client.R().Get(fmt.Sprintf("https://pornbox.com/media/%d/stream", int(mediaID)))
	if err != nil {
		return
	}

	value = gjson.Get(getSrc.String(), `qualities.@reverse.0.src`)
	var src string
	src, ok = value.Value().(string)
	if !ok {
		err = fmt.Errorf("PornBox: (video) `value.(string)` value: %v ", value)
		return
	}
	var name string
	name, err = title(res.String())
	if err != nil {
		return
	}
	cr = ContentResource{name, src}
	return
}

func (pb *PornBox) Video() (cr ContentResource, err error) {
	client := resty.New()
	defer client.Close()
	cr, err = video(pb.url, client)
	return
}

func Queue(input []string) (result []ContentResource, errorList []string,
) {
	client := resty.New()
	defer client.Close()
	for _, url := range input {
		cr, err := video(url, client)
		if err != nil {
			errorList = append(errorList, fmt.Sprintf("ERROR: %s (%v)", url, err))
			continue
		}
		result = append(result, cr)
	}
	return
}
