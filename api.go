package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
)

func deblobJson(body string) sortedCurrency {

	var list sortedCurrency
	var read start
	err := json.Unmarshal([]byte(body), &read)
	if err != nil {
		log.Fatal("unmarshal error:  ", err)
	}

	for _, value := range read.Lines {

		for _, v := range filter {
			if value.CurrencyTypeName == v {
				value.Enabled = true
				list = append(list, value)
			}
		}

	}

	list = append(list, currency{"Chaos Orb", 1.0, false})
	return list
}

func request() (string, sortedCurrency) {
	req, err := http.NewRequest("GET", getLeagueLink(), nil)
	if err != nil {
		log.Fatal("req creation err: ", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("resp err: ", err)
	}
	defer resp.Body.Close()
	//fmt.Printf("StatusCode: %d\n", resp.StatusCode)

	bytes, _ := io.ReadAll(resp.Body)
	body := string(bytes)

	listCurrency := deblobJson(body)
	sort.Sort(listCurrency)

	return postData(listCurrency), listCurrency

}

func getPostData() string {
	data, _ := request()
	return data
}

func getCurrencyList() sortedCurrency {
	_, data := request()
	return data
}

func rateo(x, y float64, currencyType string) (int64, float64) {

	ratio := x / y
	//Moltiplicatore di scam
	newRatio := ratio * currencySettingMap[currencyType].Multiplier

	increment := (newRatio - ratio) * y
	//rappresenta il ricavo in chaos
	noX := currencySettingMap[currencyType].Gain / increment
	if currencySettingMap[currencyType].Rounding {
		noX = float64(round(int(noX)))
	} else {
		app := (int64(noX*10) % 10)
		if app >= 5 {
			noX = ((noX * 10) + (10 - float64(int64(noX*10)%10))) / 10
		} else {
			noX = ((noX * 10) - float64(int64(noX*10)%10)) / 10
		}
	}

	noY := noX * newRatio

	return int64(noX), noY
}

func postData(list sortedCurrency) string {
	y := 1
	var output = "[spoiler]"

	for _, value := range list {

		//fmt.Println("\n")
		//fmt.Print(value.CurrencyTypeName, " : ", value.ChaosEquivalent)
		//fmt.Print("\t Ratios:")
		x := 0

		for _, v := range list {

			if value == v {
				continue
			}
			rX, rY := rateo(value.ChaosEquivalent, v.ChaosEquivalent, value.CurrencyTypeName)
			//fmt.Println("")
			//fmt.Printf("Give %v : %v\t", value.CurrencyTypeName, rX)
			//fmt.Printf("Receive %v : %.0f\t", v.CurrencyTypeName, rY)
			//fmt.Printf("\n")
			output += fmt.Sprintf("[spoiler=\" ~price %.0f/%v %v\"][linkItem location=\"Stash%s\" league=\"%s\" x=\"%v\" y=\"%v\"]", rY, rX, currencyName[v.CurrencyTypeName], getConfigFile().Stash, getConfigFile().CurrentLeague, x, y)
			x++
			output += "[/spoiler]"
		}

		y++
	}

	return output

}

func round(x int) int {
	if (x % 20) >= 10 {
		return x - (x % 20) + 20
	} else if (x % 20) == 0 {
		return x
	}
	return x - (x % 20)
}

func getTradingSitePostLink() string {
	return getConfigFile().PostLink
}

func (v sortedCurrency) Len() int      { return len(v) }
func (v sortedCurrency) Swap(i, j int) { v[i], v[j] = v[j], v[i] }
func (v sortedCurrency) Less(i, j int) bool {
	return currencyName[v[i].CurrencyTypeName] < currencyName[v[j].CurrencyTypeName]
}

func updateConfigFile(p string, c string, s string) {
	cFile := getConfigFile()
	cFile.Stash = s
	cFile.PostLink = p
	cFile.CurrentLeague = c
	b, err := json.Marshal(cFile)
	if err != nil {
		fmt.Println("Error marshaling json")
	}
	err = os.WriteFile(configFileName, b, 0666)
	if err != nil {
		fmt.Println("Error writing file")
	}
}

func getLeagueLink() string {
	return "https://poe.ninja/api/data/currencyoverview?league=" + getConfigFile().CurrentLeague + "&type=Currency"
}

func getConfigFile() configFile {
	b, err := os.ReadFile(configFileName)
	if err != nil {
		fmt.Println("Error Reading File")
	}
	//fmt.Println(string(b))
	var cFile configFile
	err = json.Unmarshal(b, &cFile)
	if err != nil {
		fmt.Println("Error unmarshaling json")
	}
	//fmt.Println(tradeSiteLink)
	return cFile
}

func getCurrencyJSON() currencyStart {
	b, err := os.ReadFile(currencyFileName)
	if err != nil {
		fmt.Println("Error Reading File")
	}
	var sFile currencyStart
	err = json.Unmarshal(b, &sFile)
	if err != nil {
		fmt.Println("Error unmarshaling json")
	}
	return sFile

}

func updateCurrencyJSON() {
	var g currencyStart
	for k, v := range currencySettingMap {
		var app currencyJSON
		app.CurrencyName = k
		app.Settings = currencySetting{Rounding: v.Rounding, Multiplier: v.Multiplier, Gain: v.Gain}
		g.Gianni = append(g.Gianni, app)
	}
	b, err := json.Marshal(g)
	if err != nil {
		fmt.Println("Error marshaling json")
	}
	err = os.WriteFile(currencyFileName, b, 0666)
	if err != nil {
		fmt.Println("Error writing file")
	}
}

func initCurrencySettingMap() map[string]currencySetting {
	var app = make(map[string]currencySetting)
	for _, v := range cJSON.Gianni {
		app[v.CurrencyName] = currencySetting{Rounding: v.Settings.Rounding, Multiplier: v.Settings.Multiplier, Gain: v.Settings.Gain}
	}
	return app
}

func updateCurrencySettingMap(skip string, mult string, gain string, rounding bool) {
	delete(currencySettingMap, skip)
	m, _ := strconv.ParseFloat(mult, 64)
	g, _ := strconv.ParseFloat(gain, 64)
	currencySettingMap[skip] = currencySetting{Rounding: rounding, Multiplier: m, Gain: g}
}
