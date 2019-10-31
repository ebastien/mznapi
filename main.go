package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func compile() string {
	cmd := exec.Command("minizinc",
		"--compile",
		"--input-from-stdin",
		"--output-to-stdout",
		"--no-output-ozn")
	cmd.Stdin = strings.NewReader("var int: a; solve satisfy;")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

func solve(flatzinc string) string {
	solve := exec.Command("minizinc",
		"--solver", "org.gecode.gecode",
		"--input-from-stdin")
	solve.Stdin = strings.NewReader(flatzinc)
	out, err := solve.Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

func main() {
	fmt.Println("Hello, World!")

	flatzinc := compile()
	fmt.Print(flatzinc)

	solution := solve(flatzinc)
	fmt.Print(solution)
}
