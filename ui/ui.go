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

	keys    *tview.List
	prevIdx int // previous index of key, used for pagination.
	max     int // max keys per page.
	page    int // current page.

	value  *tview.TextView
	search *tview.InputField
	format uint
}

type UiOpts struct {
	Format uint
	Max    int
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
			SetHighlightFullLine(true).
			SetWrapAround(false),

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
		SetTitle(" Keys - page: 1").
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(1, 1, 2, 2)
	ui.value.SetBorder(true).
		SetTitle(" Value ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(1, 1, 2, 2)
	ui.search.SetBorder(true)

	// Override key bindings.

	// Set changed.
	ui.keys.SetChangedFunc(
		func(index int, mainText string, secondaryText string, shortcut rune) {

			if index == 0 && index == ui.prevIdx {
			}

			value, err := ui.db.Get([]byte(secondaryText), nil)
			if err != nil {
				log.Fatal(fmt.Errorf("list change key err: %s", err))
			}
			ui.value.SetText(ui.fmtOut(value))
		},
	)

	// Load inital keys.
	ui.loadKeyBatch(nil, 100)

	return ui
}

func (ui *ui) Run() {
	// Run the application.
	if err := ui.app.SetRoot(ui.layout, true).Run(); err != nil {
		panic(err)
	}
}

func (ui *ui) nextKeyBatch(from []byte) {

	iter := ui.db.NewIterator(
		&util.Range{Start: from, Limit: nil},
		nil,
	)
	defer iter.Release()

	var i int
	for iter.Next() && i < ui.max {
		key := iter.Key()
		ui.keys.AddItem(ui.fmtOut(key), string(key), 0, nil)
		i++
	}
	ui.page++
	ui.keys.SetTitle(fmt.Sprintf(" Keys - page: %d", ui.page))
}

func (ui *ui) previousKeyBatch(to []byte) {

	iter := ui.db.NewIterator(
		&util.Range{Start: nil, Limit: to},
		nil,
	)
	defer iter.Release()

	var i int
	for iter.Prev() && i < ui.max {
		key := iter.Key()
		ui.keys.AddItem(ui.fmtOut(key), string(key), 0, nil)
		i++
	}
	ui.page--
	ui.keys.SetTitle(fmt.Sprintf(" Keys - page: %d", ui.page))
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
