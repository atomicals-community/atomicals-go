package pythonparse

import (
	"fmt"
	"os/exec"
	"strings"
)

// parse_operation_from_script
func ParseAtomicalsOperation(script string, height int64) (string, string, error) {
	cmd := exec.Command("python3", "atomicals-core/witness/python-parse/parse.py", script, fmt.Sprintf("%d", height))
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing command:", err)
	}
	fmt.Println(string(output))
	op_name, payloadStr, err := parseOperation(string(output))
	if err != nil {
		return "", "", err
	}
	if strings.Contains(string(output), "parents") {
		panic("find it!")
	}

	return op_name, payloadStr, nil
}

func parseOperation(input string) (string, string, error) {
	// Split the input string by spaces
	parts := strings.Split(input, "\n")

	// Check if we have enough parts and the correct format
	if len(parts) == 3 {
		// Return the last part, which should be "nft"
		return parts[0], parts[1], nil
	}

	// If the format doesn't match, return an error
	return "", "", fmt.Errorf("invalid input format")
}
