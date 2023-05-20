package main

import (
	"context"
	"fmt"
	cfg "iHAapi/internal/config_readers"
	pngr "iHAapi/internal/pinger"
	rds "iHAapi/internal/redis"
	wb "iHAapi/internal/web"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	tick = 1 * time.Minute
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	// Чтение списка url
	pckPath, _ := os.Getwd()
	cfgPath := pckPath + "\\configs\\url_list.txt"
	urlList := cfg.UrlReader(cfgPath)

	// Создание клиента Redis
	rdsClient := rds.PrepareRedis(ctx)

	// Создание пустой структуры для сохранения промежуточных результатов (см. /internal/pinger/pinger.go)
	result := new(pngr.GeneralizedPingResult)

	// Запуск воркера
	go worker(urlList, ctx, tick, result, *rdsClient)

	// Запуск сервера
	go wb.Server(ctx, rdsClient)

	// Отслеживание системных сигналов (частично совместим с Windows) для остановки сервера
	sigChan := make(chan os.Signal, 1)
	exitChan := make(chan struct{})
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			s := <-sigChan
			cancel()
			switch s {
			case syscall.SIGINT:
				fmt.Println("Catch: SIGNAL INTERRUPT | Server stopping")
				close(exitChan)
				return
			case os.Interrupt:
				fmt.Println("Catch: SIGNAL INTERRUPT | Server stopping")
				close(exitChan)
				return
			case syscall.SIGTERM:
				fmt.Println("Catch: SIGNAL TERMINATE | Server stopping")
				close(exitChan)
				return
			case syscall.SIGKILL:
				fmt.Println("Catch: SIGNAL KILL | Server stopping")
				close(exitChan)
				return
			}
		}
	}()
	<-exitChan

	time.Sleep(3 * time.Second)
}

func worker(urlList []string, ctx context.Context, tick time.Duration, gpr *pngr.GeneralizedPingResult, rdsc rds.RedisClient) {
	heartbeat := time.Tick(1 * time.Millisecond)

	var wg sync.WaitGroup
	for {
		select {
		case <-heartbeat:
			heartbeat = time.Tick(tick)

			gpr.New()

			for _, url := range urlList {
				wg.Add(1)
				go func() {
					pngr.Ping(url, gpr, ctx, tick)
					wg.Done()
				}()
			}
			wg.Wait()
			rds.AddUserDataToRedis(*gpr, rdsc)
		}
	}
}
