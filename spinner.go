package spinner

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type Spinner struct {
	mu         sync.Mutex
	frames     []string
	index      int
	active     bool
	stop       chan struct{}
	writer     io.Writer
	interval   func() time.Duration
	color      func() string
	hideCursor bool
}

type Option func(*Spinner)

func WithWriter(w io.Writer) Option {
	return func(s *Spinner) {
		s.writer = w
	}
}

func WithInterval(d time.Duration) Option {
	return func(s *Spinner) {
		s.interval = func() time.Duration {
			return d
		}
	}
}

func WithFrames(frames []string) Option {
	return func(s *Spinner) {
		s.frames = frames
	}
}

func WithIntervalFunc(f func() time.Duration) func(*Spinner) {
	return func(s *Spinner) {
		s.interval = f
	}
}

func WithColor(color string) func(*Spinner) {
	return func(s *Spinner) {
		s.color = func() string { return color }
	}
}

func WithColorFunc(f func() string) func(*Spinner) {
	return func(s *Spinner) {
		s.color = f
	}
}

func WithHideCursor(hide bool) func(*Spinner) {
	return func(s *Spinner) {
		s.hideCursor = hide
	}
}

var defaultFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

const (
	hideCursorSeq = "\033[?25l"
	showCursorSeq = "\033[?25h"
)

func New(opts ...Option) *Spinner {
	s := &Spinner{
		frames:     defaultFrames,
		stop:       make(chan struct{}),
		writer:     os.Stderr,
		interval:   func() time.Duration { return 60 * time.Millisecond },
		color:      func() string { return White },
		hideCursor: true,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Spinner) Start() {
	s.mu.Lock()
	if s.active {
		s.mu.Unlock()
		return
	}
	s.active = true
	if s.hideCursor {
		fmt.Fprint(s.writer, hideCursorSeq)
	}
	s.mu.Unlock()

	go func() {
		for {
			select {
			case <-s.stop:
				return
			default:
				s.mu.Lock()
				fmt.Fprintf(s.writer, "\r%s%s%s", s.color(), s.frames[s.index], Reset)
				s.index = (s.index + 1) % len(s.frames)
				s.mu.Unlock()
				time.Sleep(s.interval())
			}
		}
	}()
}

func (s *Spinner) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.active {
		s.active = false
		s.stop <- struct{}{}
		fmt.Fprint(s.writer, "\r \r")
		if s.hideCursor {
			fmt.Fprint(s.writer, showCursorSeq)
		}
	}
}

func Color256(n int) string {
	if n < 0 || n > 255 {
		return ""
	}
	return fmt.Sprintf("\033[38;5;%dm", n)
}

const (
	Black  = "\033[38;5;0m"
	Green  = "\033[38;5;2m"
	Olive  = "\033[38;5;3m"
	Navy   = "\033[38;5;4m"
	Teal   = "\033[38;5;6m"
	Silver = "\033[38;5;7m"
	Grey   = "\033[38;5;8m"
	Red    = "\033[38;5;9m"
	Lime   = "\033[38;5;10m"
	Yellow = "\033[38;5;11m"
	Blue   = "\033[38;5;12m"
	Aqua   = "\033[38;5;14m"
	White  = "\033[38;5;15m"
	Reset  = "\033[0m"
)

// Spinner styles
var (
	Dots1               = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	Dots2               = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
	Dots3               = []string{"⠋", "⠙", "⠚", "⠞", "⠖", "⠦", "⠴", "⠲", "⠳", "⠓"}
	Dots4               = []string{"⠄", "⠆", "⠇", "⠋", "⠙", "⠸", "⠰", "⠠", "⠰", "⠸", "⠙", "⠋", "⠇", "⠆"}
	Dots5               = []string{"⠋", "⠙", "⠚", "⠒", "⠂", "⠂", "⠒", "⠲", "⠴", "⠦", "⠖", "⠒", "⠐", "⠐", "⠒", "⠓", "⠋"}
	Dots6               = []string{"⠁", "⠉", "⠙", "⠚", "⠒", "⠂", "⠂", "⠒", "⠲", "⠴", "⠤", "⠄", "⠄", "⠤", "⠴", "⠲", "⠒", "⠂", "⠂", "⠒", "⠚", "⠙", "⠉", "⠁"}
	Dots7               = []string{"⠈", "⠉", "⠋", "⠓", "⠒", "⠐", "⠐", "⠒", "⠖", "⠦", "⠤", "⠠", "⠠", "⠤", "⠦", "⠖", "⠒", "⠐", "⠐", "⠒", "⠓", "⠋", "⠉", "⠈"}
	Dots8               = []string{"⠁", "⠁", "⠉", "⠙", "⠚", "⠒", "⠂", "⠂", "⠒", "⠲", "⠴", "⠤", "⠄", "⠄", "⠤", "⠠", "⠠", "⠤", "⠦", "⠖", "⠒", "⠐", "⠐", "⠒", "⠓", "⠋", "⠉", "⠈", "⠈"}
	Dots9               = []string{"⢹", "⢺", "⢼", "⣸", "⣇", "⡧", "⡗", "⡏"}
	Dots10              = []string{"⢄", "⢂", "⢁", "⡁", "⡈", "⡐", "⡠"}
	Dots11              = []string{"⠁", "⠂", "⠄", "⡀", "⢀", "⠠", "⠐", "⠈"}
	Dots12              = []string{"⢀⠀", "⡀⠀", "⠄⠀", "⢂⠀", "⡂⠀", "⠅⠀", "⢃⠀", "⡃⠀", "⠍⠀", "⢋⠀", "⡋⠀", "⠍⠁", "⢋⠁", "⡋⠁", "⠍⠉", "⠋⠉", "⠋⠉", "⠉⠙", "⠉⠙", "⠉⠩", "⠈⢙", "⠈⡙", "⢈⠩", "⡀⢙", "⠄⡙", "⢂⠩", "⡂⢘", "⠅⡘", "⢃⠨", "⡃⢐", "⠍⡐", "⢋⠠", "⡋⢀", "⠍⡁", "⢋⠁", "⡋⠁", "⠍⠉", "⠋⠉", "⠋⠉", "⠉⠙", "⠉⠙", "⠉⠩", "⠈⢙", "⠈⡙", "⠈⠩", "⠀⢙", "⠀⡙", "⠀⠩", "⠀⢘", "⠀⡘", "⠀⠨", "⠀⢐", "⠀⡐", "⠀⠠", "⠀⢀", "⠀⡀"}
	Line                = []string{"-", "\\", "|", "/"}
	Pipe                = []string{"┤", "┘", "┴", "└", "├", "┌", "┬", "┐"}
	SimpleDots          = []string{".  ", ".. ", "...", "   "}
	SimpleDotsScrolling = []string{".  ", ".. ", "...", " ..", "  .", "   "}
	Star                = []string{"✶", "✸", "✹", "✺", "✹", "✷"}
	Flip                = []string{"_", "_", "_", "-", "`", "`", "'", "´", "-", "_", "_", "_"}
	Hamburger           = []string{"☱", "☲", "☴"}
	GrowVertical        = []string{"▁", "▃", "▄", "▅", "▆", "▇", "▆", "▅", "▄", "▃"}
	GrowHorizontal      = []string{"▏", "▎", "▍", "▌", "▋", "▊", "▉", "▊", "▋", "▌", "▍", "▎"}
	Balloon             = []string{" ", ".", "o", "O", "@", "*", " "}
	Noise               = []string{"▓", "▒", "░"}
	Bounce              = []string{"⠁", "⠂", "⠄", "⠂"}
	BoxBounce           = []string{"▖", "▘", "▝", "▗"}
	BoxBounce2          = []string{"▌", "▀", "▐", "▄"}
	Triangle            = []string{"◢", "◣", "◤", "◥"}
	Arc                 = []string{"◜", "◠", "◝", "◞", "◡", "◟"}
	Circle              = []string{"◡", "⊙", "◠"}
	SquareCorners       = []string{"◰", "◳", "◲", "◱"}
	CircleQuarters      = []string{"◴", "◷", "◶", "◵"}
	CircleHalves        = []string{"◐", "◓", "◑", "◒"}
	Moon                = []string{"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"}
	Smiley              = []string{"😄 ", "😝 "}
	Monkey              = []string{"🙈 ", "🙈 ", "🙉 ", "🙊 "}
	Hearts              = []string{"💛 ", "💙 ", "💜 ", "💚 ", "❤️ "}
	Clock               = []string{"🕛 ", "🕐 ", "🕑 ", "🕒 ", "🕓 ", "🕔 ", "🕕 ", "🕖 ", "🕗 ", "🕘 ", "🕙 ", "🕚 "}
	Earth               = []string{"🌍 ", "🌎 ", "🌏 "}
	Material            = []string{"█▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁", "██▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁", "███▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁", "████▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁", "██████▁▁▁▁▁▁▁▁▁▁▁▁▁▁", "██████▁▁▁▁▁▁▁▁▁▁▁▁▁▁", "███████▁▁▁▁▁▁▁▁▁▁▁▁▁", "████████▁▁▁▁▁▁▁▁▁▁▁▁", "█████████▁▁▁▁▁▁▁▁▁▁▁", "█████████▁▁▁▁▁▁▁▁▁▁▁", "██████████▁▁▁▁▁▁▁▁▁▁", "███████████▁▁▁▁▁▁▁▁▁", "█████████████▁▁▁▁▁▁▁", "██████████████▁▁▁▁▁▁", "██████████████▁▁▁▁▁▁", "▁██████████████▁▁▁▁▁", "▁██████████████▁▁▁▁▁", "▁██████████████▁▁▁▁▁", "▁▁██████████████▁▁▁▁", "▁▁▁██████████████▁▁▁", "▁▁▁▁█████████████▁▁▁", "▁▁▁▁██████████████▁▁", "▁▁▁▁██████████████▁▁", "▁▁▁▁▁██████████████▁", "▁▁▁▁▁██████████████▁", "▁▁▁▁▁██████████████▁", "▁▁▁▁▁▁██████████████", "▁▁▁▁▁▁██████████████", "▁▁▁▁▁▁▁█████████████", "▁▁▁▁▁▁▁█████████████", "▁▁▁▁▁▁▁▁████████████", "▁▁▁▁▁▁▁▁████████████", "▁▁▁▁▁▁▁▁▁███████████", "▁▁▁▁▁▁▁▁▁███████████", "▁▁▁▁▁▁▁▁▁▁██████████", "▁▁▁▁▁▁▁▁▁▁██████████", "▁▁▁▁▁▁▁▁▁▁▁▁████████", "▁▁▁▁▁▁▁▁▁▁▁▁▁███████", "▁▁▁▁▁▁▁▁▁▁▁▁▁▁██████", "▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁█████", "▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁█████", "█▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁████", "██▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁███", "██▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁███", "███▁▁▁▁▁▁▁▁▁▁▁▁▁▁███", "████▁▁▁▁▁▁▁▁▁▁▁▁▁▁██", "█████▁▁▁▁▁▁▁▁▁▁▁▁▁▁█", "█████▁▁▁▁▁▁▁▁▁▁▁▁▁▁█", "██████▁▁▁▁▁▁▁▁▁▁▁▁▁█", "████████▁▁▁▁▁▁▁▁▁▁▁▁", "█████████▁▁▁▁▁▁▁▁▁▁▁", "█████████▁▁▁▁▁▁▁▁▁▁▁", "█████████▁▁▁▁▁▁▁▁▁▁▁", "█████████▁▁▁▁▁▁▁▁▁▁▁", "███████████▁▁▁▁▁▁▁▁▁", "████████████▁▁▁▁▁▁▁▁", "████████████▁▁▁▁▁▁▁▁", "██████████████▁▁▁▁▁▁", "██████████████▁▁▁▁▁▁", "▁██████████████▁▁▁▁▁", "▁██████████████▁▁▁▁▁", "▁▁▁█████████████▁▁▁▁", "▁▁▁▁▁████████████▁▁▁", "▁▁▁▁▁████████████▁▁▁", "▁▁▁▁▁▁███████████▁▁▁", "▁▁▁▁▁▁▁▁█████████▁▁▁", "▁▁▁▁▁▁▁▁█████████▁▁▁", "▁▁▁▁▁▁▁▁▁█████████▁▁", "▁▁▁▁▁▁▁▁▁█████████▁▁", "▁▁▁▁▁▁▁▁▁▁█████████▁", "▁▁▁▁▁▁▁▁▁▁▁████████▁", "▁▁▁▁▁▁▁▁▁▁▁████████▁", "▁▁▁▁▁▁▁▁▁▁▁▁███████▁", "▁▁▁▁▁▁▁▁▁▁▁▁███████▁", "▁▁▁▁▁▁▁▁▁▁▁▁▁███████", "▁▁▁▁▁▁▁▁▁▁▁▁▁███████", "▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁█████", "▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁████", "▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁████", "▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁████", "▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁███", "▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁███", "▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁██", "▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁██", "▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁██", "▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁█", "▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁█", "▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁█", "▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁", "▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁", "▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁", "▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁"}
)

// Helpers

func GreyPulse(interval time.Duration) func() string {
	return ColorPulse(238, 255, interval)
}

func ColorPulse(start, end int, duration time.Duration) func() string {
	t := time.Now()
	direction := 1
	color := start
	return func() string {
		if time.Since(t) > duration {
			t = time.Now()
			color += direction
			if color > end {
				color = end
				direction = -1
			}
			if color < start {
				color = start
				direction = 1
			}
		}
		return Color256(color)
	}
}

func SpeedupInterval(start, end, duration time.Duration) func() time.Duration {
	var t time.Time
	return func() time.Duration {
		if t.IsZero() {
			t = time.Now()
		}
		x := time.Since(t).Microseconds()
		y := duration.Microseconds()
		if x > y {
			return end
		}
		progress := float64(x) / float64(y)
		return time.Duration(float64(start.Nanoseconds())*(1-progress) + float64(end.Nanoseconds())*progress)
	}
}
