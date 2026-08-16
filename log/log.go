package log

import (
	"log"
	"os"

	c "github.com/TwiN/go-color"
)

type Tlog string

func (g Tlog) Warn() {

	warninfo := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	warninfo.Printf("\t%sWARN%s\t%s", c.Yellow, c.Reset, g)
}
func (g Tlog) Error() {

	warninfo := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	warninfo.Printf("\t%sERROR%s\t%s", c.Red, c.Reset, g)
}
func (g Tlog) Info() {

	warninfo := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	warninfo.Printf("\t%sINFO%s\t%s", c.Cyan, c.Reset, g)
}
