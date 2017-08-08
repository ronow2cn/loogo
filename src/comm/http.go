package comm

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

// ============================================================================
const (
	DefaultTimeOutSecond = 5
)

// ============================================================================

func HttpGet(addr string) (ret string) {
	res, err := http.Get(addr)
	if err != nil {
		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	return string(body)
}

// ============================================================================

func HttpPost(addr string, data url.Values) (ret string) {
	res, err := http.PostForm(addr, data)
	if err != nil {
		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	return string(body)
}

func HttpGetT(addr string, timeout int) (ret string, err error) {
	if timeout < 0 {
		timeout = DefaultTimeOutSecond
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*time.Duration(timeout)) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeout))) //设置发送接受数据超时
				return conn, nil
			},
		},
	}

	res, err := client.Get(addr)
	if err != nil {
		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	ret = string(body)

	return
}

func HttpPostT(addr string, data url.Values, timeout int) (ret string, err error) {
	if timeout < 0 {
		timeout = DefaultTimeOutSecond
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*time.Duration(timeout)) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeout))) //设置发送接受数据超时
				return conn, nil
			},
		},
	}

	res, err := client.PostForm(addr, data)
	if err != nil {
		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	ret = string(body)

	return
}
