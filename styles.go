package main

import (
	"fmt"
	"sort"
	"strings"
)

type styleVariant struct {
	prefix  string
	pre     string
	post    string
	symbol  string
	spacing string
}

type styleStruct struct {
	stdout styleVariant
	stderr styleVariant
}

var style styleStruct

var styles = map[string]styleStruct{
	"normal": {
		stdout: styleVariant{prefix: "[stdout]"},
		stderr: styleVariant{prefix: "[stderr]"},
	},
	"dim": {
		stdout: styleVariant{prefix: "[stdout]", pre: "\033[2m", post: "\033[0m"},
		stderr: styleVariant{prefix: "[stderr]", pre: "\033[2m", post: "\033[0m"},
	},
	"color": {
		stdout: styleVariant{prefix: "", pre: "\033[32m", post: "\033[0m"},
		stderr: styleVariant{prefix: "", pre: "\033[31m", post: "\033[0m"},
	},
	"dimu": {
		stdout: styleVariant{prefix: "", pre: "\033[2m", post: "\033[0m", symbol: "✅"},
		stderr: styleVariant{prefix: "", pre: "\033[2m", post: "\033[0m", symbol: "❌"},
	},
	"dimred": {
		stdout: styleVariant{prefix: "", pre: "\033[2m", post: "\033[0m"},
		stderr: styleVariant{prefix: "", pre: "\033[2;31m", post: "\033[0m"},
	},
}

func (variant *styleVariant) Println(txt string) {
	fmt.Println(strings.Join([]string{variant.prefix, variant.symbol, variant.spacing, variant.pre, txt, variant.post}, ""))
}

func (style *styleStruct) applyConfig() {
	if cfg.Prefix != nil {
		style.stdout.prefix = *cfg.Prefix
		style.stderr.prefix = *cfg.Prefix
	}
	if cfg.StdoutPrefix != nil {
		style.stdout.prefix = *cfg.StdoutPrefix
	}
	if cfg.StderrPrefix != nil {
		style.stderr.prefix = *cfg.StderrPrefix
	}
	if cfg.Spacing != nil {
		style.stdout.spacing = strings.Repeat(" ", *cfg.Spacing)
		style.stderr.spacing = style.stdout.spacing
	} else {
		if style.stdout.prefix != "" || style.stderr.prefix != "" {
			style.stdout.spacing = " "
			style.stderr.spacing = " "
		}
	}
}

func demo() {
	var names = make([]string, 0, len(styles))
	for name := range styles {
		names = append(names, name)
	}
	sort.SliceStable(names, func(i, j int) bool {
		if names[i] == "normal" || names[j] == "normal" {
			return names[i] == "normal"
		} else {
			return strings.Compare(names[i], names[j]) < 0
		}
	})

	fmt.Println("logcli - all modes of styling:")
	for _, name := range names {
		variant := styles[name]
		variant.applyConfig()
		fmt.Printf("\n'%s'\n", name)
		variant.stdout.Println("normal print (stdout)")
		variant.stderr.Println("error-message (stderr)")
	}
}
