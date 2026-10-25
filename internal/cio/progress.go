package cio

import (
	"fmt"
	"math"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// ProgressBar ports Symfony\Component\Console\Helper\ProgressBar as used by
// Composer (default formats, 28 char bar, 1/25s redraw throttle). On
// non-decorated output every redraw is written on a new line and the final
// line-break is left to the caller, exactly like Symfony.
type ProgressBar struct {
	io             *IO
	max, step      int
	stepWidth      int
	percent        float64
	barWidth       int
	format         string
	overwrite      bool
	previous       *string
	lastWrite      time.Time
	start          time.Time
	minRedraw      time.Duration
	maxRedraw      time.Duration
	writeCount     int
	terminalWidth  int
	explicitFormat bool
}

// NewProgressBar is ConsoleIO::getProgressBar().
func (o *IO) NewProgressBar(max int) *ProgressBar {
	b := &ProgressBar{io: o, barWidth: 28, overwrite: o.Decorated, minRedraw: time.Second / 25, maxRedraw: time.Second, start: time.Now(), terminalWidth: TerminalWidth()}
	b.setMax(max)
	return b
}

func (b *ProgressBar) setMax(max int) {
	b.format = ""
	if max < 0 {
		max = 0
	}
	b.max = max
	if max > 0 {
		b.stepWidth = len(strconv.Itoa(max))
	} else {
		b.stepWidth = 4
	}
}

// Start displays the bar at step 0.
func (b *ProgressBar) Start(max int) {
	b.start = time.Now()
	b.step = 0
	b.percent = 0
	b.setMax(max)
	b.display()
}

// SetProgress moves the bar, redrawing like Symfony's throttled setProgress.
func (b *ProgressBar) SetProgress(step int) {
	if b.max > 0 && step > b.max {
		b.max = step
	} else if step < 0 {
		step = 0
	}
	freq := float64(b.max) / 10
	if b.max == 0 {
		freq = 1
	}
	prevPeriod := int(float64(b.step) / freq)
	currPeriod := int(float64(step) / freq)
	b.step = step
	if b.max > 0 {
		b.percent = float64(step) / float64(b.max)
	} else {
		b.percent = 0
	}
	interval := time.Since(b.lastWrite)
	if b.max == step {
		b.display()
		return
	}
	if interval < b.minRedraw {
		return
	}
	if prevPeriod != currPeriod || interval >= b.maxRedraw {
		b.display()
	}
}

// MaxSteps returns the current maximum.
func (b *ProgressBar) MaxSteps() int { return b.max }

// Progress is ProgressBar::getProgress.
func (b *ProgressBar) Progress() int { return b.step }

// Finish completes the bar.
func (b *ProgressBar) Finish() {
	if b.max == 0 {
		b.max = b.step
	}
	if b.step == b.max && !b.overwrite {
		return
	}
	b.SetProgress(b.max)
}

// Clear removes the bar when overwriting is supported.
func (b *ProgressBar) Clear() {
	if !b.overwrite {
		return
	}
	b.resolveFormat()
	b.write("")
}

func (b *ProgressBar) resolveFormat() {
	if b.format != "" {
		return
	}
	nomax := b.max == 0
	switch {
	case b.io.Verbosity >= Debug:
		b.format = pick(nomax, " %current% [%bar%] %elapsed:6s% %memory:6s%", " %current%/%max% [%bar%] %percent:3s%% %elapsed:6s%/%estimated:-6s% %memory:6s%")
	case b.io.Verbosity >= VeryVerbose:
		b.format = pick(nomax, " %current% [%bar%] %elapsed:6s%", " %current%/%max% [%bar%] %percent:3s%% %elapsed:6s%/%estimated:-6s%")
	case b.io.Verbosity >= Verbose:
		b.format = pick(nomax, " %current% [%bar%] %elapsed:6s%", " %current%/%max% [%bar%] %percent:3s%% %elapsed:6s%")
	default:
		b.format = pick(nomax, " %current% [%bar%]", " %current%/%max% [%bar%] %percent:3s%%")
	}
}

func pick(nomax bool, a, b string) string {
	if nomax {
		return a
	}
	return b
}

func (b *ProgressBar) display() {
	if b.io.Verbosity <= Quiet {
		return
	}
	b.resolveFormat()
	b.write(b.buildLine())
}

var rePlaceholder = regexp.MustCompile(`(?i)%([a-z\-_]+)(?::([^%]+))?%`)

func (b *ProgressBar) buildLine() string {
	line := b.render()
	if w := Width(StripTags(line)); w > b.terminalWidth {
		// setBarWidth: max(1, $size)
		b.barWidth = max(1, b.barWidth-w+b.terminalWidth)
		line = b.render()
	}
	return line
}

func (b *ProgressBar) render() string {
	return rePlaceholder.ReplaceAllStringFunc(b.format, func(m string) string {
		sm := rePlaceholder.FindStringSubmatch(m)
		var text string
		switch sm[1] {
		case "bar":
			text = b.bar()
		case "elapsed":
			// time() - $startTime: whole seconds of the wall clock
			text = formatTime(time.Now().Unix() - b.start.Unix())
		case "estimated":
			est := 0.0
			if b.step > 0 {
				est = math.Round(float64(time.Now().Unix()-b.start.Unix()) / float64(b.step) * float64(b.max))
			}
			text = formatTime(int64(est))
		case "memory":
			var ms runtime.MemStats
			runtime.ReadMemStats(&ms)
			text = formatMemory(int64(ms.Sys))
		case "current":
			text = fmt.Sprintf("%*d", b.stepWidth, b.step)
		case "max":
			text = strconv.Itoa(b.max)
		case "percent":
			text = strconv.Itoa(int(math.Floor(b.percent * 100)))
		default:
			return m
		}
		if sm[2] != "" {
			text = phpSprintfS(sm[2], text)
		}
		return text
	})
}

// phpSprintfS applies a "%<spec>" string conversion such as "3s" or "-6s".
func phpSprintfS(spec, text string) string {
	if !strings.HasSuffix(spec, "s") {
		return text
	}
	n, err := strconv.Atoi(strings.TrimSuffix(spec, "s"))
	if err != nil {
		return text
	}
	return fmt.Sprintf("%*s", n, text)
}

func (b *ProgressBar) bar() string {
	var complete int
	if b.max > 0 {
		complete = int(math.Floor(b.percent * float64(b.barWidth)))
	} else {
		complete = int(math.Min(5, float64(b.barWidth)/15)*float64(b.writeCount)) % b.barWidth
	}
	barChar := "="
	if b.max == 0 {
		barChar = "-"
	}
	s := strings.Repeat(barChar, complete)
	if complete < b.barWidth {
		empty := b.barWidth - complete - 1
		if empty < 0 {
			empty = 0
		}
		s += ">" + strings.Repeat("-", empty)
	}
	return s
}

func (b *ProgressBar) write(msg string) {
	if b.previous != nil && *b.previous == msg {
		return
	}
	orig := msg
	if b.overwrite {
		if b.previous != nil {
			// Cursor::moveToColumn(1) + clearLine()
			msg = "\x1b[1G\x1b[2K" + msg
		}
	} else if b.step > 0 {
		msg = "\n" + msg
	}
	b.previous = &orig
	b.lastWrite = time.Now()
	b.io.WriteErrorRaw(msg)
	b.writeCount++
}

func formatTime(secs int64) string {
	formats := []struct {
		min  int64
		text string
		div  int64
	}{{0, "< 1 sec", 0}, {1, "1 sec", 0}, {2, "secs", 1}, {60, "1 min", 0}, {120, "mins", 60}, {3600, "1 hr", 0}, {7200, "hrs", 3600}, {86400, "1 day", 0}, {172800, "days", 86400}}
	for i, f := range formats {
		if secs >= f.min && (i == len(formats)-1 || secs < formats[i+1].min) {
			if f.div == 0 {
				return f.text
			}
			return strconv.FormatInt(secs/f.div, 10) + " " + f.text
		}
	}
	return ""
}

func formatMemory(m int64) string {
	switch {
	case m >= 1<<30:
		return fmt.Sprintf("%.1f GiB", float64(m)/(1<<30))
	case m >= 1<<20:
		return fmt.Sprintf("%.1f MiB", float64(m)/(1<<20))
	case m >= 1<<10:
		return fmt.Sprintf("%d KiB", m/1024)
	}
	return fmt.Sprintf("%d B", m)
}
