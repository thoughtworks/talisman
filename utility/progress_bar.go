package utility

import (
	"fmt"
	"os"

	"github.com/cheggaaa/pb/v3"
)

func GetProgressBar(out *os.File, title string) progressBar {
	if isTerminal(out) {
		return &defaultProgressBar{title: title}
	} else {
		return &noOpProgressBar{}
	}
}

func isTerminal(out *os.File) bool {
	fileInfo, _ := out.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

type progressBar interface {
	Start(int)
	Increment()
	Finish()
}

type noOpProgressBar struct {
}

func (d *noOpProgressBar) Start(int) {}

func (d *noOpProgressBar) Increment() {}

func (d *noOpProgressBar) Finish() {}

type defaultProgressBar struct {
	bar   *pb.ProgressBar
	title string
}

func (d *defaultProgressBar) Start(total int) {
	template := fmt.Sprintf(`{{ red "%s:" }} {{counters .}} {{ bar . "<" "-" (cycle . "↖" "↗" "↘" "↙" ) "." ">"}} {{percent . | rndcolor }} {{green}} {{blue}}`, d.title)
	bar := pb.ProgressBarTemplate(template).New(total)
	bar.Set(pb.Terminal, true)
	d.bar = bar.Start()
}

func (d *defaultProgressBar) Increment() {
	d.bar.Increment()
}

func (d *defaultProgressBar) Finish() {
	d.bar.Finish()
}
