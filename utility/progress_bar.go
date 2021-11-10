package utility

import (
	"os"

	"github.com/cheggaaa/pb/v3"
)

func GetProgressBar() progressBar {
	if isTerminal() {
		return &defaultProgressBar{}
	} else {
		return &noOpProgressBar{}
	}
}

func isTerminal() bool {
	fileInfo, _ := os.Stdout.Stat()
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
	bar *pb.ProgressBar
}

func (d *defaultProgressBar) Start(total int) {
	bar := pb.ProgressBarTemplate(`{{ red "Talisman Scan:" }} {{counters .}} {{ bar . "<" "-" (cycle . "↖" "↗" "↘" "↙" ) "." ">"}} {{percent . | rndcolor }} {{green}} {{blue}}`).New(total)
	bar.Set(pb.Terminal, true)
	d.bar = bar.Start()
}

func (d *defaultProgressBar) Increment() {
	d.bar.Increment()
}

func (d *defaultProgressBar) Finish() {
	d.bar.Finish()
}
