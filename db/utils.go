package db

import (
	"bufio"
	"os"
)

type stdOut struct {
	r *os.File
	w *os.File
}

func ConnectToStdout() *stdOut {
	r, w, _ := os.Pipe()
	os.Stdout = w

	return &stdOut{
		r: r,
		w: w,
	}
}

func ReadStdout(stdOut *stdOut) []string {
	result := make([]string, 0)
	scanner := bufio.NewScanner(stdOut.r)
	done := make(chan struct{})

	go func() {
		// Read line by line and process it
		for scanner.Scan() {
			line := scanner.Text()
			result = append(result, line)
		}

		// We're all done, unblock the channel
		done <- struct{}{}
	}()

	stdOut.w.Close()
	<-done

	return result
}
