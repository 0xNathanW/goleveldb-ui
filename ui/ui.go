package ui

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// Format options.
const (
	Str = iota // String
	Hex        // Hex
	Num        // Number
	Bin        // Binary
)

type ui struct {
	db     *leveldb.DB
	app    *tview.Application
	layout *tview.Grid
	left   *tview.Flex
	keys   *tview.List
	value  *tview.TextView
	search *tview.InputField
	format uint
}

type UiOpts struct {
	Format uint
}

func NewUI(dbPath string, opts *UiOpts) *ui {

	// database init
	db, err := leveldb.OpenFile(dbPath, &opt.Options{ErrorIfMissing: true})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to open database: %w", err))
	}
	if err := db.SetReadOnly(); err != nil {
		log.Fatal(fmt.Errorf("failed to set database to read-only: %w", err))
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
			ShowSecondaryText(false).
			SetHighlightFullLine(true),

		value: tview.NewTextView().
			SetScrollable(true).
			SetWrap(true),

		search: tview.NewInputField().
			SetLabel("Search: ").
			SetPlaceholder("Ctrl-S").
			SetPlaceholderTextColor(tcell.ColorGreen).
			SetFieldBackgroundColor(tcell.ColorRed),

		format: opts.Format,
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

	// Set changed.
	ui.keys.SetChangedFunc(
		func(index int, mainText string, secondaryText string, shortcut rune) {
			value, err := ui.db.Get([]byte(secondaryText), nil)
			if err != nil {
				log.Fatal(fmt.Errorf("list change key err: %s", err))
			}
			ui.value.SetText(ui.fmtOut(value))
		},
	)

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
	var i int
	for iter.Next() && i < h {
		key := iter.Key()
		ui.keys.AddItem(ui.fmtOut(key), string(key), 0, nil)
	}

}

func (ui *ui) fmtOut(data []byte) string {
	switch ui.format {
	case Str:
		return string(data)
	case Hex:
		return fmt.Sprintf("%x", data)
	case Num:
		return fmt.Sprintf("%d", data)
	case Bin:
		return fmt.Sprintf("%b", data)
	}
	return ""
}
