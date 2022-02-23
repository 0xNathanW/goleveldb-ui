# Goleveldb-ui
A terminal interface for github.com/syndtr/goleveldb

Slightly rough around the edges.

## Installation

git clone repo: `git clone https://github.com/0xNathanW/goleveldb-ui`

build: `go build -o={executable name} main.go` 

## Usage

Run the executable with the path to the database as a value to the `--db` flag.

Eg. `ui --db=./testdb`

Other flags:

<img width="320" alt="flags" src="https://user-images.githubusercontent.com/86011312/155250313-ebd2a6a5-e813-44fa-9271-cb47fd74149c.png">

Where key format is in {"string", "hex"}

Value format is in {"string", "hex", "number", "binary"}

Pages have been implemented for efficiency, such that a monolith of keys are not stored in memory at one time.

## Navigation 

Use up and down keys to maneuver keys, optionally using pgUp and pgDown to quickly change pages.

Ctrl+O and Ctrl+P can be used to change the relative size of the two columns.

### Search

Use Ctrl+S to enter the search input field and the Escape key to exit.

Along with conventional searching, inputs with a preceding `$` can be used to change ui elements while using the program, rather than using flags.

- `$key={hex, string}` changes key formatting.
- `$val={hex, string, number, binary}` changes value formatting.
- `$max={int}` change max number of keys displayed on a page.



