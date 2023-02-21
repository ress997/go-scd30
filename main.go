package main

import (
	"log"
	"time"
	"flag"

	"github.com/pvainio/scd30"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

var (
	interval int
	i2c string
)

func init() {
	flag.StringVar(&i2c, "i2c", "", "I²C bus to use")
	flag.IntVar(&interval, "interval", 5, "The time in seconds between CO₂ readings")
	flag.Parse()
}

func main() {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
	bus, err := i2creg.Open(i2c)
	if err != nil {
		log.Fatal(err)
	}
	defer bus.Close()
	dev, err := scd30.Open(bus)
	if err != nil {
		log.Fatal(err)
	}
	dev.StartMeasurements(uint16(interval))

	for {
		time.Sleep(time.Duration(interval) * time.Second)

		if hasMeasurement, err := dev.HasMeasurement(); err != nil {
			log.Fatalf("error %v", err)
		} else if !hasMeasurement {
			return
		}

		m, err := dev.GetMeasurement()
		if err != nil {
			log.Fatalf("error %v", err)
		}

		log.Printf("Temp: %.4g°C, Hum: %.3g%%, CO₂: %.4g ppm", m.Temperature, m.Humidity, m.CO2)
	}
}
