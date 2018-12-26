package main

import (
	"time"
	"net"
	"net/http"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"context"
)

func TimeoutDialerC(cTimeout time.Duration, rwTimeout time.Duration) func(ctx context.Context, network, address string) (c net.Conn, err error) {

	return func(ctx context.Context,netw, addr string) (net.Conn, error) {
		dialer:=net.Dialer{
			Timeout:   cTimeout,
		}
		conn,err:=dialer.DialContext(ctx,netw,addr)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(rwTimeout))
		return conn, nil
	}

}

func main() {

	uri := "http://www.baidu.com"
	connectTimeout := 1 * time.Second
	readWriteTimeout := 3 * time.Second


	c := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			DialContext: TimeoutDialerC(connectTimeout,readWriteTimeout),
			TLSHandshakeTimeout: 1 * time.Second,
		},
	}
	req, err := http.NewRequest(http.MethodPost, uri, nil)
	if err != nil {
		fmt.Println("req error:" + err.Error())
		return
	}
	resp, err := c.Do(req)
	if err != nil {
		fmt.Println("do error,err:" + err.Error())
		return
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	fmt.Println("llllll==",string(respBody))

}
