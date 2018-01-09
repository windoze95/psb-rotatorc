package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/subcommands"
	"github.com/orange-lightsaber/psb-rotatord/rotator"
	"github.com/orange-lightsaber/psb-rotatord/sockets"
)

type Run struct {
	name    string
	compkey string
	freq    int
	delay   int
	year    int
	month   int
	day     int
	initial int
}

func (cmd *Run) SetFlags(f *flag.FlagSet) {
	f.StringVar(&cmd.name, "name", "", "")
	f.StringVar(&cmd.compkey, "compkey", "", "")
	f.IntVar(&cmd.freq, "freq", -1, "")
	f.IntVar(&cmd.delay, "delay", -1, "")
	f.IntVar(&cmd.year, "year", -1, "")
	f.IntVar(&cmd.month, "month", -1, "")
	f.IntVar(&cmd.day, "day", -1, "")
	f.IntVar(&cmd.initial, "initial", -1, "")
}

func (cmd *Run) CheckFlags() (r bool) {
	switch {
	case cmd.name == "":
	case cmd.compkey == "":
	case cmd.freq < 0:
	case cmd.delay < 0:
	case cmd.year < 0:
	case cmd.month < 0:
	case cmd.day < 0:
	case cmd.initial < 0:
	default:
		r = true
	}
	return
}

type preRunCmd struct {
	Run
}

func (*preRunCmd) Name() string     { return "prerun" }
func (*preRunCmd) Synopsis() string { return "Start pre-run." }
func (*preRunCmd) Usage() string {
	return `prerun [-name] <string> [-compkey] <string> [-freq] <int> [-delay] <int> [-year] <int> [-month] <int> [-day] <int> [-initial] <int>:
  Starts pre-run operation.
`
}

func (cmd *preRunCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if ok := cmd.CheckFlags(); !ok {
		fmt.Print("invalid arguments")
		os.Exit(2)
	}
	req := sockets.Request{
		sockets.PreRun_Req,
		rotator.RunConfigData{
			CompatibilityKey: cmd.compkey,
			Name:             cmd.name,
			Frequency:        cmd.freq,
			RotationDelay:    cmd.delay,
			Year:             rotator.Year{Duration: cmd.year},
			Month:            rotator.Month{Duration: cmd.month},
			Day:              rotator.Day{Duration: cmd.day},
			Initial:          rotator.Initial{Duration: cmd.initial},
		}}
	res, err := req.NewRequest()
	if err != nil {
		fmt.Print(err)
		os.Exit(2)
	}
	if res.Error != "" {
		fmt.Print(res.Error)
		os.Exit(2)
	}
	// fmt.Print(res.Response)
	return subcommands.ExitSuccess
}
