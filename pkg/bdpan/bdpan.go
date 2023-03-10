package bdpan

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gizmo-ds/toybox/internal/utils"
	"github.com/go-resty/resty/v2"
)

const (
	appId     = "250528"
	panUA     = "netdisk;P2SP;3.0.0.8;netdisk;11.12.3;ANG-AN00;android-android;10.0;JSbridge4.4.0;jointBridge;1.1.0;"
	browserUA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/110.0"
)

var client = resty.New()

func GetFileHash(file *os.File) (filename, md5 string, size int, err error) {
	filename = filepath.Base(file.Name())
	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		return
	}
	size = buf.Len()
	md5 = utils.MD5(buf.Bytes())
	return
}

func RapidUpload(bduss, savename, md5 string, size int, retry ...int) (err error) {
	_url := "https://pan.baidu.com/rest/2.0/xpan/file?method=create&app_id=" + appId

	body := url.Values{}
	body.Set("path", savename)
	body.Set("size", strconv.Itoa(size))
	body.Set("isdir", "0")
	body.Set("rtype", "0")
	blockList, _ := json.Marshal([]string{md5})
	body.Set("block_list", string(blockList))
	body.Set("mode", "1")

	var result struct {
		Errno int `json:"errno"`
	}
	_, err = client.R().
		SetHeader("user-agent", panUA).
		SetCookie(&http.Cookie{Name: "BDUSS", Value: bduss}).
		SetBody(body.Encode()).
		SetResult(&result).
		Post(_url)
	if err != nil {
		return
	}
	if result.Errno != 0 {
		if result.Errno == 2 && len(retry) > 0 && retry[0] > 0 {
			return RapidUpload(bduss, savename, md5, size, retry[0]-1)
		}
		return errors.New("errno: " + strconv.Itoa(result.Errno))
	}
	return
}

func CheckBduss(bduss string) (name string, err error) {
	var result struct {
		Data any `json:"data"`
	}
	_, err = client.R().
		SetCookie(&http.Cookie{Name: "BDUSS", Value: bduss}).
		SetHeader("user-agent", browserUA).
		SetResult(&result).
		Get("https://news.baidu.com/passport")
	if err != nil {
		return
	}
	if _, ok := result.Data.([]any); ok {
		err = errors.New("BDUSS invalid")
		return
	}
	name = result.Data.(map[string]any)["displayname"].(string)
	return
}
