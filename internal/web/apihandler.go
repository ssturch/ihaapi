package web

import (
	db "iHAapi/internal/redis"
	"net/http"
	"strconv"
	//"strings"
)

// Хэндлер для /user
func getLatencyHandler(rdsc db.RedisClient, gep *GeneralEndPointStat) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		urlReceived := r.URL.Query().Get("url")
		if r.Method == "POST" {
			if urlReceived == "" {
				ReturnResponseCode(&w, 404, "error: 'Url is null'", nil)
				return
			}

			redisKey := "urllist:" + urlReceived
			res, err := rdsc.GetData(redisKey)

			if err != nil {
				ReturnResponseCode(&w, 404, "error: 'This url not available'", nil)
				return
			}
			w.Write(res)
			(*gep)[r.URL.String()]++
			rdsc.AddData("admin:endpoint_stat", gep.ToJSON())

		}
	}
}

// Хэндлер для /user/min
func getMinLatencyHandler(rdsc db.RedisClient, gep *GeneralEndPointStat) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			redisKey := "minLatUrl"
			res, err := rdsc.GetData(redisKey)
			if err != nil || res == nil {
				ReturnResponseCode(&w, 404, "error: 'This data not available'", nil)
				return
			}
			w.Write(res)
			(*gep)[r.URL.String()]++
			rdsc.AddData("admin:endpoint_stat", gep.ToJSON())
		}
	}
}

// Хэндлер для /user/max
func getMaxLatencyHandler(rdsc db.RedisClient, gep *GeneralEndPointStat) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			redisKey := "maxLatUrl"
			res, err := rdsc.GetData(redisKey)
			if err != nil || res == nil {
				ReturnResponseCode(&w, 404, "error: 'This data not available'", nil)
				return
			}
			w.Write(res)
			(*gep)[r.URL.String()]++
			rdsc.AddData("admin:endpoint_stat", gep.ToJSON())
		}
	}
}

// Хэндлер для /admin/statistics
func getAdminStatHandler(rdsc db.RedisClient, gep *GeneralEndPointStat) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			redisKey := "admin:endpoint_stat"
			jsonByte, err := rdsc.GetData(redisKey)
			if err != nil {
				ReturnResponseCode(&w, 404, "error: 'Url is null'", nil)
				return
			}
			w.Write(jsonByte)
		}
	}
}

// Функция для возврата желаемого кода
func ReturnResponseCode(w *http.ResponseWriter, code int, msg string, details []byte) {
	(*w).Header().Add("code", strconv.Itoa(code))
	(*w).Header().Add("message", msg)
	(*w).Header().Add("details", string(details))
	(*w).WriteHeader(http.StatusNotFound)
}
