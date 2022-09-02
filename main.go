package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

type Data struct {
	Id     int32   `json:"id"`
	Market int32   `json:"market"`
	Price  float64 `json:"price"`
	Volume float64 `json:"volume"`
	IsBuy  bool    `json:"is_buy"`
}

type Stats struct {
	Market         int32   `json:"market"`
	TotalVolume    float64 `json:"total_volume"`
	MeanPrice      float64 `json:"mean_price"`
	MeanVolume     float64 `json:"mean_volume"`
	VolumeWeighted float64 `json:"volume_weighted_average_price"`
	PercentBuy     int32   `json:"percent_buy"`
	TotalPoints    int32   `json:"-"`
}

func main() {
	var datapoint Data

	markets := map[int32]Stats{}

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		line := sc.Text()
		if len(line) == 0 || line == "BEGIN" {
			continue
		}

		if line == "END" {
			break
		}

		err := json.Unmarshal([]byte(line), &datapoint)
		if err != nil {
			log.Printf("Failed to unmarshal: %s\n", line)
		}

		stat, ok := markets[datapoint.Market]
		if !ok {
			stat = Stats{}
		}

		isBuy := 0
		if datapoint.IsBuy {
			isBuy = 1
		}

		stat.TotalVolume = stat.TotalVolume + datapoint.Volume
		stat.MeanPrice = (stat.MeanPrice*float64(stat.TotalPoints) + datapoint.Price) / (float64(stat.TotalPoints) + 1)
		stat.MeanVolume = (stat.MeanVolume*float64(stat.TotalPoints) + datapoint.Volume) / (float64(stat.TotalPoints) + 1)
		stat.VolumeWeighted = (stat.VolumeWeighted*stat.TotalVolume + datapoint.Price*datapoint.Volume) / (stat.TotalVolume + datapoint.Volume)
		stat.PercentBuy = (stat.PercentBuy*stat.TotalPoints + int32(isBuy)) / (stat.TotalPoints + 1)
		stat.TotalPoints += 1

		markets[datapoint.Market] = stat

	}
	if err := sc.Err(); err != nil {
		log.Fatalf("scan file error: %v", err)
	}

	for market, stat := range markets {
		stat.Market = market
		output, err := json.Marshal(stat)
		if err != nil {
			log.Fatal("failed to unmarshal")
		}
		log.Println(string(output))
	}
}
