package logger

import (
	"comm/config"
	"fmt"
	"github.com/op/go-logging"
	"os"
	"time"
)

// ============================================================================

var (
	DefaultLogger = logging.MustGetLogger("")
)

// ============================================================================

var (
	logFmt = logging.MustStringFormatter(
		`%{color}%{time:2006-01-02 15:04:05.000} [%{level:.1s}] => %{message} | [%{shortfile:s}] %{color:reset}`,
	)
)

// ============================================================================

var (
	quit chan int

	logFn    string
	logH     *os.File
	backEnd  *logging.LogBackend
	lastDate string
)

// ============================================================================

func Open(fn string) {
	if quit != nil {
		return
	}

	quit = make(chan int)

	logFn = fn

	// prepare backend
	if logFn == "" {
		// output to stdout
		logH = os.Stdout
	} else {
		// output to file
		logH, _ = os.OpenFile(logFn, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.FileMode(0644))
		fi, err := logH.Stat()
		if err == nil {
			lastDate = fi.ModTime().Format("2006-01-02")
		} else {
			lastDate = time.Now().Format("2006-01-02")
		}
	}

	// set backend
	backEnd = logging.NewLogBackend(logH, "", 0)
	logging.SetBackend(
		logging.NewBackendFormatter(backEnd, logFmt),
	)

	// set level
	if config.Common.LogLevel == "Info" {
		logging.SetLevel(logging.INFO, "")
	} else {
		logging.SetLevel(logging.DEBUG, "")
	}

	// log rotation check
	if logFn != "" {
		// open check
		checkRotation()

		go func() {
			ticker := time.NewTicker(1 * time.Minute)
			defer func() {
				ticker.Stop()
				quit <- 0
			}()

			for {
				select {
				case <-ticker.C:
					// runtime check
					checkRotation()

				case <-quit:
					// close check
					checkRotation()
					return
				}
			}
		}()
	}
}

func Close() {
	if quit == nil {
		return
	}

	if logFn != "" {
		quit <- 0
		<-quit
	}
	close(quit)
	quit = nil

	logH.Close()
}

// ============================================================================

func checkRotation() {

	newDate := time.Now().Format("2006-01-02")

	if newDate != lastDate {
		// rename log
		newpath := fmt.Sprintf("log/%s.%s", logFn, lastDate)
		os.MkdirAll("log", os.FileMode(0755))
		os.Rename(logFn, newpath)

		// create new log
		oldH := logH
		logH, _ = os.Create(logFn)
		backEnd.Logger.SetOutput(logH)
		oldH.Close()

		// update last date
		lastDate = newDate
	}
}
