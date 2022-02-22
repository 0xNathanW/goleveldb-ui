package ui

import (
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
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
	db   *leveldb.DB
	iter iterator.Iterator

	app    *tview.Application
	layout *tview.Grid
	left   *tview.Flex
	ratio  int // Proportion left to right. Between -4 and 4.

	keys *tview.List
	max  int // max keys per page.
	page int // current page.
	prev int // 0 if last page move was forward, 1 if backward.

	value  *tview.TextView
	search *tview.InputField
	keyFmt uint
	valFmt uint
}

type UiOpts struct {
	KeyFmt uint // Format of keys.
	ValFmt uint // Format of values.
	Max    int  // Max keys per page.
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
			SetPlaceholderTextColor(tcell.ColorGreen).
			SetPlaceholderStyle(tcell.StyleDefault).
			SetFieldBackgroundColor(tcell.ColorDefault).
			SetLabel("Search: ").
			SetPlaceholder("Ctrl+S").
			SetPlaceholderTextColor(tcell.ColorGreen),

		keyFmt: opts.KeyFmt,
		valFmt: opts.ValFmt,
	}

	// Left side flex layout.
	ui.left.AddItem(ui.search, 3, 0, false)
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
					ui.nextKeyBatch()
					return nil
				} else { // If not last item, pass to default handler.
					return event
				}

			case tcell.KeyUp:
				if ui.keys.GetCurrentItem() == 0 && ui.page > 1 {
					ui.previousKeyBatch()
					return nil
				} else { // If not first item, pass to default handler.
					return event
				}

			case tcell.KeyCtrlS:
				ui.app.SetFocus(ui.search)
				return nil

			case tcell.KeyPgUp:
				if ui.page > 1 {
					ui.previousKeyBatch()
				}
				return nil
			case tcell.KeyPgDn:
				ui.nextKeyBatch()
				return nil

			case tcell.KeyCtrlLeftSq:
				ui.shiftRatio(false)
			case tcell.KeyCtrlRightSq:
				ui.shiftRatio(true)
			}

			return event
		},
	)
	ui.value.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {

			switch key := event.Key(); key {

			case tcell.KeyBackspace:
				ui.app.SetFocus(ui.keys)
				return nil

			case tcell.KeyCtrlLeftSq:
			}
			return event
		},
	)

	// Set changed.
	ui.keys.SetChangedFunc(
		func(index int, mainText string, secondaryText string, shortcut rune) {
			value, _ := ui.db.Get([]byte(secondaryText), nil) // err not possible.
			ui.value.SetText(ui.valOut(value))
		},
	)

	// Set selected.
	ui.keys.SetSelectedFunc(
		func(index int, mainText string, secondaryText string, shortcut rune) {
			ui.app.SetFocus(ui.value)
		},
	)

	// Set search.
	ui.search.SetDoneFunc(
		func(key tcell.Key) {

			if key == tcell.KeyEnter {
				ui.handleInput(ui.search.GetText())
				ui.app.SetFocus(ui.keys)
			}
			if key == tcell.KeyEscape {
				ui.app.SetFocus(ui.keys)
			}
		},
	)

	// Load inital keys.
	iter.First() // Move to first key.
	ui.nextKeyBatch()

	return ui
}

func (ui *ui) Run() {
	// Run the application.
	if err := ui.app.SetRoot(ui.layout, true).Run(); err != nil {
		ui.iter.Release()
		ui.db.Close()
		log.Fatal(err)
	}
}

func (ui *ui) nextKeyBatch() {

	if ui.prev == 1 { // If last page move was backward, move iterator.
		for i := 0; i < ui.keys.GetItemCount()+1; i++ {
			ui.iter.Next()
		}
	}

	// Is there a next page?
	if !ui.iter.Next() {
		return
	}
	ui.iter.Prev() // Return to current key.

	ui.keys.Clear()

	var i int
	for i < ui.max { // Populate list.
		key := ui.iter.Key()
		ui.keys.AddItem(ui.keyOut(key), string(key), 0, nil)
		i++
		if !ui.iter.Next() {
			break
		}
	}

	ui.page++
	ui.keys.SetTitle(fmt.Sprintf(" Keys - page: %d ", ui.page))
	ui.keys.SetCurrentItem(0)
	ui.prev = 0
}

func (ui *ui) previousKeyBatch() {

	if ui.prev == 0 { // If last page move was forward, move iterator.
		for i := 0; i < ui.keys.GetItemCount()+1; i++ {
			ui.iter.Prev()
		}
	}

	ui.keys.Clear()
	var i int
	for i < ui.max { // Populate list.
		key := ui.iter.Key()
		ui.keys.InsertItem(0, ui.keyOut(key), string(key), 0, nil)
		i++
		if !ui.iter.Prev() {
			break
		}
	}

	ui.page--
	ui.keys.SetTitle(fmt.Sprintf(" Keys - page: %d ", ui.page))
	ui.keys.SetCurrentItem(ui.keys.GetItemCount() - 1)
	ui.prev = 1
}

func (ui *ui) keyOut(data []byte) string {
	switch ui.keyFmt {
	case Str:
		return string(data)
	case Hex:
		return hex.EncodeToString(data)
	}
	return ""
}

func (ui *ui) keyIn(data string) []byte {
	switch ui.keyFmt {
	case Hex:
		i, _ := hex.DecodeString(data)
		return i
	default:
		return []byte(data)
	}
}

func (ui *ui) valOut(data []byte) string {
	switch ui.valFmt {
	case Str:
		return string(data)
	case Hex:
		return hex.EncodeToString(data)
	case Num:
		return fmt.Sprintf("%d", data)
	case Bin:
		return fmt.Sprintf("%b", data)
	}
	return ""
}

func (ui *ui) handleInput(input string) {

	if input == "" {
		ui.iter = ui.db.NewIterator(nil, nil)
		ui.iter.First()
		ui.prev, ui.page = 0, 0
		ui.nextKeyBatch()
		return
	}

	if input[0] == '$' {
		params := strings.Split(input[1:], "=")
		if len(params) != 2 {
			return
		}

		switch strings.TrimSpace(params[0]) {

		case "key":
			switch strings.TrimSpace(params[1]) {

			case "hex":
				ui.keyFmt = Hex

			case "string":
				ui.keyFmt = Str
			}

		case "val":
			switch strings.TrimSpace(params[1]) {

			case "hex":
				ui.valFmt = Hex
			case "string":
				ui.valFmt = Str
			case "number":
				ui.valFmt = Num
			case "binary":
				ui.valFmt = Bin
			}

		case "max":
			ui.max, _ = strconv.Atoi(params[1])

		}
		return
	}

	ui.iter = ui.db.NewIterator(util.BytesPrefix([]byte(ui.keyIn(input))), nil)
	ui.iter.First()
	ui.prev, ui.page = 0, 0
	ui.nextKeyBatch()
}

func (ui *ui) shiftRatio(right bool) {

	if right {
		ui.ratio++
	} else {
		ui.ratio--
	}

	// max ration is 4.
	if ui.ratio > 4 {
		ui.ratio = 4
		return
	}
	if ui.ratio < -4 {
		ui.ratio = -4
		return
	}

	if ui.ratio == 0 {
		ui.layout.SetColumns(0, 0)
	} else if ui.ratio < 0 {
		ui.layout.SetColumns(0, ui.ratio)
	} else {
		ui.layout.SetColumns(-ui.ratio, 0)
	}
}
