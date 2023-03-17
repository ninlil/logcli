package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/alexflint/go-arg"
)

var cfg struct {
	StdoutPrefix *string `arg:"-o,--" placeholder:"PREFIX" help:"prefix for <stdout>"`
	StderrPrefix *string `arg:"-e,--" placeholder:"PREFIX" help:"prefix for <stderr>"`
	Prefix       *string `arg:"-p,--" placeholder:"PREFIX" help:"prefix for any line"`
	Mode         string  `arg:"-m,--" default:"normal" help:"mode/style of output"`
	Spacing      *int    `arg:"-s,--" help:"spacing between prefix and output"`
	Demo         bool    `arg:"--demo" help:"prints a demo of all modes"`
}

type line struct {
	isErr bool
	line  string
}

func addLine(outputChan chan<- line, isErr bool, text string) {
	row := line{
		isErr: isErr,
		line:  text,
	}
	outputChan <- row
}

func parseArgs() []string {
	var cmdLine []string
	for i, arg := range os.Args {
		if arg == "--" {
			cmdLine = os.Args[i+1:]
			os.Args = os.Args[:i]
			break
		}
	}
	if len(cmdLine) < 1 {
		fmt.Fprintf(os.Stderr, "no command to execute - add '--' last and whatever command you want to run after\n")
	}
	arg.MustParse(&cfg)

	if cfg.Demo {
		demo()
		os.Exit(0)
	}

	if len(cmdLine) < 1 {
		os.Exit(1)
	}

	if _, ok := styles[cfg.Mode]; ok {
		style = styles[cfg.Mode]
	} else {
		style = styles["normal"]
	}

	return cmdLine
}

func main() {
	cmdLine := parseArgs()

	style.applyConfig()

	cmd := exec.Command(cmdLine[0], cmdLine[1:]...)

	cmd.Stdin = os.Stdin

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating stdout pipe: %v\n", err)
		return
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating stderr pipe: %v\n", err)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting command: %v\n", err)
		return
	}

	outputChan := make(chan line, 100)

	var wgCapture = new(sync.WaitGroup)
	wgCapture.Add(2)
	go captureOutput(stdoutPipe, false, outputChan, wgCapture)
	go captureOutput(stderrPipe, true, outputChan, wgCapture)

	var wgPrint = new(sync.WaitGroup)
	wgPrint.Add(1)
	go printOutput(outputChan, wgPrint)

	err = cmd.Wait()
	wgCapture.Wait()

	close(outputChan)
	wgPrint.Wait()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Command finished with error: %v\n", err)
	} else {
		fmt.Println("Command finished successfully")
	}
	if cmd.ProcessState != nil {
		os.Exit(cmd.ProcessState.ExitCode())
	}
}

func captureOutput(pipe io.ReadCloser, isErr bool, outputChan chan<- line, wg *sync.WaitGroup) {
	defer func() {
		pipe.Close()
		wg.Done()
	}()
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		addLine(outputChan, isErr, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading output: %v\n", err)
	}
}

func printOutput(outputChan <-chan line, wg *sync.WaitGroup) {
	defer wg.Done()
	for line := range outputChan {
		var variant *styleVariant
		if line.isErr {
			variant = &style.stderr
		} else {
			variant = &style.stdout
		}

		variant.Println(line.line)

		//		fmt.Println(strings.Join([]string{variant.prefix, variant.symbol, variant.spacing, variant.pre, line.line, variant.post}, ""))
	}
}
