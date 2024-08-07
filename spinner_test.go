package spinner_test

import (
	"fmt"
	"os"
	"time"

	"github.com/tmc/spinner"
)

func ExampleSpinner_basic() {
	s := spinner.New()
	s.Start()
	time.Sleep(time.Second)
	s.Stop()
	// output:
}

func ExampleSpinner_withCustomFrames() {
	s := spinner.New(spinner.WithFrames(spinner.Dots12))
	s.Start()
	time.Sleep(3 * time.Second)
	s.Stop()
	// output:
}
func ExampleSpinner_withCustomFramesAllDots() {
	for i, f := range [][]string{
		spinner.Dots1,
		spinner.Dots2,
		spinner.Dots3,
		spinner.Dots4,
		spinner.Dots5,
		spinner.Dots6,
		spinner.Dots7,
		spinner.Dots8,
		spinner.Dots9,
		spinner.Dots10,
		spinner.Dots11,
		spinner.Dots12,
	} {
		s := spinner.New(spinner.WithFrames(f))
		fmt.Fprintf(os.Stderr, "spinner.Dots%v\n", i+1)
		s.Start()
		time.Sleep(2 * time.Second)
		s.Stop()
		fmt.Fprintf(os.Stderr, "\n")
	}
	// output:
}

func ExampleSpinner_withAdvancedOptions() {
	s := spinner.New(
		spinner.WithFrames(spinner.Dots8),
		spinner.WithIntervalFunc(
			spinner.SpeedupInterval(90*time.Millisecond, 40*time.Millisecond, time.Second*5),
		),
		spinner.WithColorFunc(spinner.GreyPulse(15*time.Millisecond)),
	)
	s.Start()
	time.Sleep(5 * time.Second)
	s.Stop()
	// output:
}
