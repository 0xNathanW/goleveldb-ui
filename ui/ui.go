package ui

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type ui struct {
	db     *leveldb.DB
	app    *tview.Application
	layout *tview.Grid
	left   *tview.Flex
	keys   *tview.List
	value  *tview.TextView
	search *tview.InputField
}

func NewUI(dbPath string) *ui {

	// database init
	db, err := leveldb.OpenFile(dbPath, &opt.Options{ErrorIfMissing: true})
	if err != nil {
		log.Fatal(err)
	}
	if err := db.SetReadOnly(); err != nil {
		log.Fatal(err)
	}

	ui := &ui{
		db: db,

		app: tview.NewApplication().
			EnableMouse(false),

		layout: tview.NewGrid().
			// Two evenly sized columns.
			SetColumns(0, 0),

		left: tview.NewFlex().
			SetDirection(tview.FlexRow),

		keys: tview.NewList().
			ShowSecondaryText(false),

		value: tview.NewTextView(),

		search: tview.NewInputField().
			SetLabel("Search: ").
			SetPlaceholder("Ctrl-S").
			SetPlaceholderTextColor(tcell.ColorGreen).
			SetFieldBackgroundColor(tcell.ColorRed),
	}

	// Left side flex layout.
	ui.left.AddItem(ui.search, 1, 0, false)
	ui.left.AddItem(ui.keys, 0, 1, true)

	// Grid layout.
	ui.layout.AddItem(ui.left, 0, 0, 1, 1, 0, 0, true)
	ui.layout.AddItem(ui.value, 0, 1, 1, 1, 0, 0, false)

	// Box setup.
	ui.keys.SetBorder(true).
		SetTitle(" Keys ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(1, 1, 2, 2)
	ui.value.SetBorder(true).
		SetTitle(" Value ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(1, 1, 2, 2)
	ui.search.SetBorder(true)

	// Load inital keys.
	ui.loadKeyBatch(nil)

	return ui
}

func (ui *ui) Run() {

	// Run the application.
	if err := ui.app.SetRoot(ui.layout, true).Run(); err != nil {
		panic(err)
	}
}

func (ui *ui) loadKeyBatch(from []byte) {

	iter := ui.db.NewIterator(
		&util.Range{Start: from, Limit: nil},
		nil,
	)
	defer iter.Release()

	_, _, _, h := ui.keys.GetInnerRect()
	for i := 0; i < h; i++ {
		if !iter.Next() {
			break
		}
		key := iter.Key()
		ui.keys.AddItem(string(key), "", 0, nil)
	}

}
