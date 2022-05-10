package memcache

import (
	"bufio"
	"fmt"
	"google.golang.org/grpc/status"
	"strings"
	"time"
)

func Get(key string) (string, error) {
	_, err := fmt.Fprint(connection, "get "+key+"\n")
	if err != nil {
		return "", status.Error(13, err.Error())
	}

	resp := make([]byte, 256)
	err = connection.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return "", status.Error(13, err.Error())
	}

	_, err = bufio.NewReader(connection).Read(resp)
	if err != nil {
		return "", status.Error(13, err.Error())
	}

	return parseResponse(resp)
}

//parseResponse парсит ответ от memcached
func parseResponse(resp []byte) (string, error) {
	stringsResp := strings.Split(string(resp), "\n")
	resultResp := strings.TrimSuffix(stringsResp[len(stringsResp)-1], "\r")

	if resultResp == "ERROR" || resultResp == "CLIENT_ERROR" {
		return "", status.Error(2, "memcache error")
	}

	if resultResp == "NOT_STORED" || resultResp == "NOT_FOUND" {
		return "", status.Error(3, "wrong value")
	}

	if len(stringsResp) > 2 {
		stringsResp[1] = strings.TrimSuffix(stringsResp[1], "\r")
		return stringsResp[1], nil
	}

	return "", nil
}

func Set(key string, value string, expiration int) error {
	request := fmt.Sprintf("set %s 0 %d %d", key, expiration, len(value))

	_, err := fmt.Fprint(connection, request+"\r\n")
	if err != nil {
		return status.Error(13, err.Error())
	}

	_, err = fmt.Fprint(connection, value+"\r\n")
	if err != nil {
		return status.Error(13, err.Error())
	}

	resp := make([]byte, 256)
	err = connection.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return status.Error(13, err.Error())
	}

	_, err = bufio.NewReader(connection).Read(resp)
	if err != nil {
		return status.Error(13, err.Error())
	}

	_, err = parseResponse(resp)

	return err
}

func Delete(key string) error {
	_, err := fmt.Fprint(connection, "delete "+key+"\n")
	if err != nil {
		return status.Error(13, err.Error())
	}

	resp := make([]byte, 256)
	err = connection.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return status.Error(13, err.Error())
	}

	_, err = bufio.NewReader(connection).Read(resp)
	if err != nil {
		return status.Error(13, err.Error())
	}

	_, err = parseResponse(resp)

	return err
}

func FlushAll() error {
	_, err := fmt.Fprint(connection, "flush_all\n")
	if err != nil {
		return status.Error(13, err.Error())
	}

	resp := make([]byte, 256)
	err = connection.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return status.Error(13, err.Error())
	}

	_, err = bufio.NewReader(connection).Read(resp)
	if err != nil {
		return status.Error(13, err.Error())
	}

	_, err = parseResponse(resp)

	return err
}
