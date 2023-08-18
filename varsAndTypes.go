package main

import (
	"fyne.io/fyne/v2/app"
	"net/http"
)

// defining types
type start struct {
	Lines           []currency  `json:"lines"`
	CurrencyDetails interface{} `json:"currencyDetails"`
}

type currency struct {
	CurrencyTypeName string  `json:"currencyTypeName"`
	ChaosEquivalent  float64 `json:"chaosEquivalent"`
	Enabled          bool
}

type currencySetting struct {
	Rounding   bool
	Multiplier float64
	Gain       float64
}

type configFile struct {
	PostLink      string `json:"postLink"`
	CurrentLeague string `json:"CurrentLeague"`
	Stash         string `json:"Stash"`
}

type sortedCurrency []currency

type currencyJSON struct {
	CurrencyName string          `json:"currencyName"`
	Settings     currencySetting `json:"settings"`
}

type currencyStart struct {
	Gianni []currencyJSON `json:"Gianni"` //sucate voglio chiamare il tag gianni :P
}

// vars for slider bindings

// maps for json debloating
var currencyName = map[string]string{
	"Gemcutter's Prism":     "gcp",
	"Glassblower's Bauble":  "bauble",
	"Orb of Regret":         "regret",
	"Vaal Orb":              "vaal",
	"Orb of Fusing":         "fuse",
	"Orb of Scouring":       "scour",
	"Chromatic Orb":         "chrom",
	"Cartographer's Chisel": "chisel",
	"Orb of Alchemy":        "alch",
	"Orb of Alteration":     "alt",
	"Jeweller's Orb":        "jew",
	"Chaos Orb":             "chaos",
}

var filter = []string{
	"Gemcutter's Prism",
	"Glassblower's Bauble",
	"Orb of Regret",
	"Vaal Orb",
	"Orb of Fusing",
	"Orb of Scouring",
	"Chromatic Orb",
	"Cartographer's Chisel",
	"Orb of Alchemy",
	"Orb of Alteration",
	"Jeweller's Orb",
	"Chaos Orb",
}

// general vars
var client http.Client
var a = app.New()
var w = a.NewWindow("pomegranate")
var cJSON = getCurrencyJSON()

var currencySettingMap = initCurrencySettingMap()

const configFileName = "./config/settings.json"
const currencyFileName = "./config/currency.json"
