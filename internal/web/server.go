package web

import (
	"context"
	"encoding/json"
	"fmt"
	rds "iHAapi/internal/redis"
	"net/http"
)

const (
	Host = "localhost"
	Port = "8080"
)

// Структура для сохранения статистики по endpoint'ам
type GeneralEndPointStat map[string]int

func (g *GeneralEndPointStat) ToJSON() []byte {
	jsonByte, _ := json.Marshal(g)
	return jsonByte
}

// Отслеживание endpoint'ов

func Server(ctx context.Context, rdsc *rds.RedisClient) {
	gep := make(GeneralEndPointStat)

	GetOldStatData(rdsc, &gep)

	srv := &http.Server{Addr: ":" + Port}

	http.HandleFunc("/user", getLatencyHandler(*rdsc, &gep))
	http.HandleFunc("/user/min", getMinLatencyHandler(*rdsc, &gep))
	http.HandleFunc("/user/max", getMaxLatencyHandler(*rdsc, &gep))
	http.HandleFunc("/admin/statistics", getAdminStatHandler(*rdsc, &gep))

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			fmt.Println("Listen and Serve process stopped!")
			return
		}
	}()

	<-ctx.Done()
	srv.Shutdown(ctx)
	fmt.Println("Server stopped!")
}

//Функция для получения старых данных по статистике из Redis

func GetOldStatData(rdsc *rds.RedisClient, gep *GeneralEndPointStat) {
	jsonData, err := rdsc.GetData("admin:endpoint_stat")
	if err != nil && jsonData == nil {
		return
	}
	var oldData GeneralEndPointStat
	err = json.Unmarshal(jsonData, &oldData)
	if err != nil {
		fmt.Println(err)
		return
	}
	*gep = oldData
}
