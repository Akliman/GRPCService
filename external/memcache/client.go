package memcache

import (
	"GRPCService/logger"
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

var (
	connection net.Conn
	reconTime  int
)

func NewConnection(ctx context.Context, adr string) error {
	err := connect(adr)
	if err != nil {
		return err
	}

	// Калждые 2 секунды проверяем коннект по серверу
	go func(ctx context.Context, adr string) {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				ping(adr)
			case <-ctx.Done():
				return
			}
		}
	}(ctx, adr)

	return nil
}

func Close() error {
	logger.LogrusLogger.Info("Stopping memcache connection")
	return connection.Close()
}

func connect(adr string) error {
	var err error
	connection, err = net.DialTimeout("tcp", adr, 1*time.Second)
	if err != nil {
		return err
	}

	return nil
}

//ping проверяет подлкючение к серверу memcached
func ping(adr string) {
	resp := make([]byte, 256)
	err := connection.SetReadDeadline(time.Now().Add(1 * time.Second))
	if err != nil {
		logger.LogrusLogger.Error("error setting deadline ", err)
	}
	_, err = fmt.Fprint(connection, "version"+"\n")
	if err != nil {
		logger.LogrusLogger.Error("error sending message ", err)
		reconTime++
		reconnect(adr)
		return
	}

	_, err = bufio.NewReader(connection).Read(resp)
	if err != nil {
		logger.LogrusLogger.Error("error getting message ", err)
		reconTime++
		reconnect(adr)
		return
	}

	if !strings.Contains(string(resp), "VERSION") {
		logger.LogrusLogger.Error("wrong response ", err)
		reconTime++
		reconnect(adr)
		return
	}

	reconTime = 0

}

func reconnect(adr string) {
	//Максимально количетсво попыток реконекта 3
	if reconTime > 3 {
		logger.LogrusLogger.Fatal("Maximum reconncetion try")
		log.Fatal("Maximum reconncetion try")
	}

	err := connect(adr)
	if err != nil {
		logger.LogrusLogger.Fatal("Error conncet to Memcache ", err)
		log.Fatal(err)
	}
}
