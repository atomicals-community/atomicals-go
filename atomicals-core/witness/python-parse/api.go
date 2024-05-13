package pythonparse

import (
	"fmt"
	"os/exec"
	"strings"
)

// parse_operation_from_script
func ParseAtomicalsOperation(script string) {
	cmd := exec.Command("python3", "witness/python-parse/parse.py", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing command:", err)
	}
	fmt.Println("output", string(output))
	if strings.Contains(string(output), "parents") {
		panic("find it!")
	}
}
