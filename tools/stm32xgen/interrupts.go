package main

import (
	"fmt"
	"strconv"
	"strings"
)

type IRQ struct {
	Name  string
	Num   int
	Descr string
}

func interrupts(r *scanner) []*IRQ {
	var irqs []*IRQ
	for r.Scan() {
		line := strings.TrimSpace(r.Text())
		if strings.HasPrefix(line, "typedef") && strings.Contains(line, "IRQn") {
			irqs = make([]*IRQ, 0, 248)
			continue
		}
		if strings.HasPrefix(line, "}") && strings.Contains(line, "IRQn_Type") {
			break
		}
		if irqs == nil {
			continue
		}
		n := strings.IndexByte(line, '=')
		if n < 0 {
			continue
		}
		name := strings.TrimSuffix(strings.TrimSpace(line[:n]), "_IRQn")
		line = strings.TrimSpace(line[n+1:])
		var descr string
		if n = strings.Index(line, "/*!<"); n > 0 {
			descr = line[n+4:]
			descr = strings.TrimSpace(strings.TrimSuffix(descr, "*/"))
			line = strings.TrimSpace(line[:n])
		}
		num, err := strconv.ParseInt(strings.TrimSuffix(line, ","), 0, 0)
		checkErr(err)
		if num < -14 || num > 247 {
			die("Bad IRQ number", num, "for", name, "interrupt.")
		}
		irqs = append(irqs, &IRQ{name, int(num), descr})
	}
	checkErr(r.Err())
	return irqs
}

func saveIRQs(irqs []*IRQ) {
	mkdir("irq")
	chdir("irq")
	defer chdir("..")
	w := create("irq.go")
	defer w.Close()
	fmt.Fprintln(
		w, "// Package irq provides list of all defined external interrupts.",
	)
	fmt.Fprintln(w, "package irq")
	fmt.Fprintln(w)
	fmt.Fprintln(w, `import "arch/cortexm/nvic"`)
	fmt.Fprintln(w, "const (")
	for _, irq := range irqs {
		if irq.Num < 0 {
			continue
		}
		fmt.Fprintf(w, "\t%s nvic.IRQ = %d // %s.\n", irq.Name, irq.Num, irq.Descr)
	}
	fmt.Fprintln(w, ")")
}
