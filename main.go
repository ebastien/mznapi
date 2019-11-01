package main

import (
	"encoding/json"
	"fmt"
	"io"
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
	cmd.Stdin = strings.NewReader("var int: a; constraint a >= 1; constraint a <= 2; solve satisfy;")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

func solve(flatzinc string) {
	solve := exec.Command("minizinc",
		"--solver", "org.gecode.gecode",
		"--input-from-stdin",
		"--output-mode", "json",
		"--solution-separator", "",
		"-a",
	)
	solve.Stdin = strings.NewReader(flatzinc)
	stdout, err := solve.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := solve.Start(); err != nil {
		log.Fatal(err)
	}
	type Solution struct {
		Variable int `json:"a"`
	}
	dec := json.NewDecoder(stdout)
	for {
		var s Solution
		if err := dec.Decode(&s); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("a = %d\n", s.Variable)
	}
	if err := solve.Wait(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println("Hello, World!")

	flatzinc := compile()
	fmt.Print(flatzinc)

	solve(flatzinc)
}
