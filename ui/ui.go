package ui

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// Format options.
const (
	Str = iota // String
	Hex        // Hex
	Num        // Number
	Bin        // Binary
)

type ui struct {
	db   *leveldb.DB
	iter iterator.Iterator

	app    *tview.Application
	layout *tview.Grid
	left   *tview.Flex

	keys *tview.List
	max  int // max keys per page.
	page int // current page.

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
	iter := db.NewIterator(nil, nil)

	ui := &ui{
		db: db,

		iter: iter,

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

		max: opts.Max,

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
		SetTitle(" Keys - page: 1 ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(1, 1, 2, 2)
	ui.value.SetBorder(true).
		SetTitle(" Value ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(1, 1, 2, 2)
	ui.search.SetBorder(true)

	// Intercept and override relevant key bindings.
	ui.keys.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {

			switch key := event.Key(); key {

			case tcell.KeyDown:
				if ui.keys.GetCurrentItem() == ui.keys.GetItemCount()-1 {
					log.Println("next page")
					ui.keys.Clear()
					ui.nextKeyBatch()
					ui.page++
					ui.keys.SetTitle(fmt.Sprintf(" Keys - page: %d ", ui.page))
					return nil
				} else { // If not last item, pass to default handler.
					return event
				}

			case tcell.KeyUp:
				if ui.keys.GetCurrentItem() == 0 && ui.page > 1 {
					log.Println("previous page")
					ui.keys.Clear()
					ui.previousKeyBatch()
					ui.page--
					ui.keys.SetTitle(fmt.Sprintf(" Keys - page: %d ", ui.page))
					return nil
				} else { // If not first item, pass to default handler.
					return event
				}

			}
			return event
		},
	)

	// Set changed.
	ui.keys.SetChangedFunc(
		func(index int, mainText string, secondaryText string, shortcut rune) {
			value, _ := ui.db.Get([]byte(secondaryText), nil) // err not possible.
			ui.value.SetText(ui.fmtOut(value))
		},
	)

	// Load inital keys.
	iter.First() // Move to first key.
	ui.page = 1  // Set page to 1.
	ui.nextKeyBatch()

	return ui
}

func (ui *ui) Run() {
	// Run the application.
	if err := ui.app.SetRoot(ui.layout, true).Run(); err != nil {
		ui.db.Close()
		log.Fatal(err)
	}
}

func (ui *ui) nextKeyBatch() {

	var i int
	for ui.iter.Next() && (i < ui.max) {
		key := ui.iter.Key()
		ui.keys.AddItem(ui.fmtOut(key), string(key), 0, nil)
		i++
		ui.iter.Next()
	}

	ui.keys.SetCurrentItem(0)
}

func (ui *ui) previousKeyBatch() {

	var i int
	for ui.iter.Prev() && i < ui.max {
		key := ui.iter.Key()
		ui.keys.InsertItem(0, ui.fmtOut(key), string(key), 0, nil)
		i++
		ui.iter.Prev()
	}

	ui.keys.SetCurrentItem(ui.keys.GetItemCount() - 1)
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

func (ui *ui) shutdown() {
	ui.app.Stop()
	ui.iter.Release()
	ui.db.Close()
}

// Idea: implement iterator into the struct, using same for next and prev page.
// If search, we can switch out the iterator.
