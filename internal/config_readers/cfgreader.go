package config_readers

import (
	"bufio"
	"os"
)

// Функция для чтения url из файла
func UrlReader(path string) []string {
	var res []string
	cfgData, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(cfgData)
	for fileScanner.Scan() {
		url := fileScanner.Text()
		res = append(res, url)
	}

	defer cfgData.Close()
	return res
}
