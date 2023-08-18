package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/pkg/browser"
	"golang.design/x/clipboard"
)

func createWindow() {

	a.Settings().SetTheme(theme.DarkTheme())

	gcpItem := fyne.NewMenuItem("Gcp", func() {
		w.SetContent(currencyContent("Gemcutter's Prism"))
	})
	chaosItem := fyne.NewMenuItem("Chaos", func() {
		w.SetContent(currencyContent("Chaos Orb"))
	})
	baubleItem := fyne.NewMenuItem("Bauble", func() {
		w.SetContent(currencyContent("Glassblower's Bauble"))
	})
	regretItem := fyne.NewMenuItem("Regret", func() {
		w.SetContent(currencyContent("Orb of Regret"))
	})
	vaalItem := fyne.NewMenuItem("Vaal", func() {
		w.SetContent(currencyContent("Vaal Orb"))
	})
	fuseItem := fyne.NewMenuItem("Fuse", func() {
		w.SetContent(currencyContent("Orb of Fusing"))
	})
	scourItem := fyne.NewMenuItem("Scour", func() {
		w.SetContent(currencyContent("Orb of Scouring"))
	})
	chromItem := fyne.NewMenuItem("Chrom", func() {
		w.SetContent(currencyContent("Chromatic Orb"))
	})
	chiselItem := fyne.NewMenuItem("Chisel", func() {
		w.SetContent(currencyContent("Cartographer's Chisel"))
	})
	alchItem := fyne.NewMenuItem("Alch", func() {
		w.SetContent(currencyContent("Orb of Alchemy"))
	})
	altItem := fyne.NewMenuItem("Alt", func() {
		w.SetContent(currencyContent("Orb of Alteration"))
	})
	jewItem := fyne.NewMenuItem("Jew", func() {
		w.SetContent(currencyContent("Jeweller's Orb"))
	})

	currencyMenu := fyne.NewMenu("Currency List", chaosItem, gcpItem, baubleItem, regretItem, vaalItem, fuseItem, scourItem, chromItem, chiselItem, altItem, alchItem, jewItem)

	sito := fyne.NewMenuItem("Click to Open Thread Post", func() {
		err := browser.OpenURL(getTradingSitePostLink())
		if err != nil {
			fmt.Println(err)
		}
		err = clipboard.Init()
		if err != nil {
			panic(err)
		}
		clipboard.Write(clipboard.FmtText, []byte(getPostData()))
	})
	siteMenu := fyne.NewMenu("Site", sito)

	settings := fyne.NewMenuItem("Settings", func() {

		cFile := getConfigFile()

		league := widget.NewEntry()
		league.Text = cFile.CurrentLeague
		leagueLabel := widget.NewLabel("Current League")
		tradeSiteLink := widget.NewEntry()
		tradeSiteLink.Text = cFile.PostLink
		tradeSiteLinkLabel := widget.NewLabel("Post Link")
		stashEntry := widget.NewEntry()
		stashEntry.Text = cFile.Stash
		stashEntryLabel := widget.NewLabel("Stash tab number")

		update := widget.NewButton("Update Configs", func() {
			updateConfigFile(tradeSiteLink.Text, league.Text, stashEntry.Text)
		})

		w.SetContent(container.NewVBox(leagueLabel, league, tradeSiteLinkLabel, tradeSiteLink, stashEntryLabel, stashEntry, update))

	})

	optionsMenu := fyne.NewMenu("Options", settings)

	mainMenu := fyne.NewMainMenu(currencyMenu, siteMenu, optionsMenu)

	w.SetOnClosed(func() {
		updateCurrencyJSON()
	})
	w.SetContent(container.NewCenter(widget.NewLabel("Select One Currency")))

	w.Resize(fyne.NewSize(1080, 720))
	w.SetMainMenu(mainMenu)
	w.ShowAndRun()

}

func currencyContent(skip string) *container.Split {

	list := getCurrencyList()
	leading := container.NewVBox()
	for _, value := range list {

		if value.CurrencyTypeName == skip {

			currencyLabel := widget.NewLabel(skip)
			currencyLabel.Alignment = fyne.TextAlignCenter
			currencyLabel.TextStyle = fyne.TextStyle{Bold: true, Italic: true}
			leading.Add(currencyLabel)

			for _, v := range list {
				if value == v {
					continue
				}
				rX, rY := rateo(value.ChaosEquivalent, v.ChaosEquivalent, value.CurrencyTypeName)
				outputLabel := widget.NewLabel(fmt.Sprintf("Give %v %s ; Receive %.0f %s", rX, value.CurrencyTypeName, rY, v.CurrencyTypeName))
				outputLabel.Alignment = fyne.TextAlignCenter
				leading.Add(outputLabel)
			}
		}
	}

	//setting the left part of the Hsplit with the currency setting sliders
	multiplierEntry := widget.NewEntry()
	multiplierEntry.Text = fmt.Sprintf("%.2f", currencySettingMap[skip].Multiplier)
	multiplierNameLabel := widget.NewLabel("Multiplier")

	gainEntry := widget.NewEntry()
	gainEntry.Text = fmt.Sprintf("%.2f", currencySettingMap[skip].Gain)
	gainNameLabel := widget.NewLabel("Chaos Gain per Trade")

	roundCheckBox := widget.NewCheck("Round?", func(b bool) {})
	roundCheckBox.Checked = currencySettingMap[skip].Rounding

	applyButton := widget.NewButton("Update", func() {
		updateCurrencySettingMap(skip, multiplierEntry.Text, gainEntry.Text, roundCheckBox.Checked)
		w.SetContent(currencyContent(skip))
	})

	trailing := container.NewVBox(multiplierNameLabel, multiplierEntry, gainNameLabel, gainEntry, roundCheckBox, applyButton)
	scatola := container.NewHSplit(leading, trailing)

	return scatola
}
