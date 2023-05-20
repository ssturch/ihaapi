package redis

import (
	"encoding/json"
	pngr "iHAapi/internal/pinger"
	"strings"
)

// Функция для добавления данных в БД Redis
func AddUserDataToRedis(res pngr.GeneralizedPingResult, rdsc RedisClient) {
	var jsonByte []byte
	var redisKey string
	for _, pingres := range res.Pr {
		jsonByte, _ = json.Marshal(pingres)
		redisKey = "urllist:" + strings.TrimLeft(pingres.Url, "https://")
		rdsc.AddData(redisKey, jsonByte)
	}

	jsonByte, _ = json.Marshal(res.Minlr)
	redisKey = "minLatUrl"
	rdsc.AddData(redisKey, jsonByte)

	jsonByte, _ = json.Marshal(res.Maxlr)
	redisKey = "maxLatUrl"
	rdsc.AddData(redisKey, jsonByte)
}
