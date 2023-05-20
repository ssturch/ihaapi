package pinger

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type GeneralizedPingResult struct {
	Pr    []PingResult
	Minlr MinLatencyResult
	Maxlr MaxLatencyResult
}

func (r *GeneralizedPingResult) New() {
	r.Pr = []PingResult{}
	r.Minlr = MinLatencyResult{Url: "Undefined", Latency: 0}
	r.Maxlr = MaxLatencyResult{Url: "Undefined", Latency: 0}
}

type PingResult struct {
	Url        string
	Err        error
	Latency    time.Duration
	StatusCode int
}

type MinLatencyResult struct {
	Url     string
	Latency time.Duration
}

type MaxLatencyResult struct {
	Url     string
	Latency time.Duration
}

//Функция для пинга url и сохранения данных в структурах выше

func Ping(url string, gpr *GeneralizedPingResult, ctx context.Context, timeout time.Duration) {
	//var temp PingResult

	temp := make(chan PingResult, 1)
	start := time.Now()

	go func() {
		if !strings.Contains(url, "https://") {
			url = "https://" + url
		}
		if resp, err := http.Get(url); err != nil {
			temp <- PingResult{Url: url, Err: err, Latency: 0, StatusCode: 0}
			close(temp)
		} else {
			t := time.Since(start).Round(time.Millisecond)
			temp <- PingResult{Url: url, Err: nil, Latency: t, StatusCode: resp.StatusCode}
			resp.Body.Close()
			close(temp)
		}
	}()

	select {
	case t, _ := <-temp:
		gpr.Pr = append(gpr.Pr, t)
		if gpr.Minlr.Latency == 0 || gpr.Minlr.Latency > t.Latency && t.StatusCode == 200 {
			gpr.Minlr.Latency = t.Latency
			gpr.Minlr.Url = t.Url
		}
		if gpr.Maxlr.Latency == 0 || gpr.Maxlr.Latency < t.Latency && t.StatusCode == 200 {
			gpr.Maxlr.Latency = t.Latency
			gpr.Maxlr.Url = t.Url
		}
		return

	case <-time.After(timeout):
		tempStruct := PingResult{Url: url, Err: errors.New("Ping timeout"), Latency: 0, StatusCode: 0}
		gpr.Pr = append(gpr.Pr, tempStruct)
		return

	case <-ctx.Done():
		fmt.Println("Pinger stopped!")
		return
	}
}
