package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"strconv"
	"strings"
	"time"
)

type WeatherData struct {
	Min, Max, Sum float32
	Count         int
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
var executionprofile = flag.String("execprofile", "", "write tarce execution to `file`")
var input = flag.String("input", "", "path to the input file to evaluate")

func main() {
	flag.Parse()

	if *executionprofile != "" {
		f, err := os.Create("./profiles/" + *executionprofile)
		if err != nil {
			log.Fatal("could not create trace execution profile: ", err)
		}
		defer f.Close()
		trace.Start(f)
		defer trace.Stop()
	}

	if *cpuprofile != "" {
		f, err := os.Create("./profiles/" + *cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	executeWrapper()

	if *memprofile != "" {
		f, err := os.Create("./profiles/" + *memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}

func executeWrapper() {
	start := time.Now()

	executeLogic()

	elapsed := time.Since(start)
	st, _ := os.Stat(*input)
	size := st.Size()

	fmt.Fprintf(os.Stderr, "Processed %.1fGB in %s\n",
		float64(size)/(1024*1024*1024), elapsed)
}

func executeLogic() {
	weatherMap := readAndProcessData()
	stationsSorted := sortWeatherStations(weatherMap)
	fmt.Print("{")
	for _, station := range stationsSorted {
		if weatherStationData := weatherMap[station]; weatherStationData != nil {
			mean := weatherStationData.Sum / float32(weatherStationData.Count)
			fmt.Printf("%s=%.1f/%.1f/%.1f", station, weatherStationData.Min, mean, weatherStationData.Max)
		}
	}
	fmt.Print("}\n")
}

func readAndProcessData() map[string]*WeatherData {
	weatherMap := make(map[string]*WeatherData)

	file, err := os.Open(*input)
	if err != nil {
		log.Fatal("Could not find file: ", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		weatherData := scanner.Text()

		// Ignore comments
		if strings.HasPrefix(weatherData, "#") {
			continue
		}

		weatherDataArr := strings.Split(weatherData, ";")
		stationName := weatherDataArr[0]
		weatherValue := weatherDataArr[1]
		weatherValueFloat, _ := strconv.ParseFloat(weatherValue, 32)
		weatherStationData := weatherMap[stationName]

		if weatherStationData == nil {
			weatherMap[stationName] = &WeatherData{Min: float32(weatherValueFloat), Max: float32(weatherValueFloat), Count: 1, Sum: float32(weatherValueFloat)}
		} else {
			weatherStationData.Count++
			if weatherStationData.Max < float32(weatherValueFloat) {
				weatherStationData.Max = float32(weatherValueFloat)
			}
			if weatherStationData.Min > float32(weatherValueFloat) {
				weatherStationData.Min = float32(weatherValueFloat)
			}
			weatherStationData.Sum += float32(weatherValueFloat)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error reading file: ", err)
	}

	return weatherMap
}

func sortWeatherStations(weatherMap map[string]*WeatherData) (keys []string) {
	keys = make([]string, 0, len(weatherMap))
	for k := range weatherMap {
		keys = append(keys, k)
	}
	return keys
}
