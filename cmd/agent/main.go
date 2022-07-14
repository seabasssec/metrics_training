package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const pollInterval = 2
const reportInterval = 10

type bufferStorage struct {
	metricValue float64
	metricName  string
	typePseudo  string
}

func float64MetricGenerator() float64 {
	seconds := time.Now().Unix()
	rand.Seed(seconds)
	return rand.Float64() + float64(rand.Intn(10000))
}

func getMetric(c chan bufferStorage) {
	time.Sleep(pollInterval * time.Second)
	metricsNameList := [29]string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc", "PollCount", "RandomValue"}
	metricsTypeList := [29]string{"gauge", "gauge", "gauge", "gauge", "gauge", "gauge", "gauge", "gauge", "gauge", "gauge", "gauge", "gauge", "gauge", "gauge", "gauge", "gauge", "gauge", "gauge", "gauge", "gauge", "gauge", "gauge", "gauge", "gauge", "gauge", "gauge", "gauge", "counter", "gauge"}
	for i := 0; i < 29; i++ {
		var buff bufferStorage = bufferStorage{metricValue: float64MetricGenerator(), metricName: metricsNameList[i], typePseudo: metricsTypeList[i]}
		c <- buff
	}
}

func send(c chan bufferStorage) {
	for {
		time.Sleep(reportInterval * time.Second)
		result, ok := <-c
		if ok != true {
			log.Fatal("ok isn't TRUE")
		}
		client := http.Client{}
		//http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
		if result.typePseudo == "gauge" {
			url := fmt.Sprintf("http://127.0.0.1:8000/update/%s/%s/%.2f", result.typePseudo, result.metricName, result.metricValue)
			resp, err := client.Post(url, "text/plain", nil)
			if err != nil {
				log.Fatalf("Failed get: %s", err)
			}
			defer resp.Body.Close()
		} else {
			url := fmt.Sprintf("http://127.0.0.1:8000/update/%s/%s/%d", result.typePseudo, result.metricName, int64(result.metricValue))
			resp, err := client.Post(url, "text/plain", nil)
			if err != nil {
				log.Fatalf("Failed get: %s", err)
			}
			defer resp.Body.Close()
		}
	}
}

func main() {
	c := make(chan bufferStorage)
	go send(c)
	for {
		getMetric(c)
	}
}
