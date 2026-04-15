package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/AatirNadim/getMe/commons"
	logger "github.com/AatirNadim/getMe/utils"
)

func ValidateJSONAndFilePath(jsonFilePath string) error {
	info, err := os.Stat(jsonFilePath)
	if err != nil {
		return fmt.Errorf("failed to stat JSON file '%s': %w", jsonFilePath, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("JSON path '%s' is not a regular file", jsonFilePath)
	}
	if info.Size() == 0 {
		return fmt.Errorf("JSON file '%s' is empty", jsonFilePath)
	}
	if info.Size() > commons.MaxJSONFileSizeBytes {
		return fmt.Errorf("JSON file '%s' size %d bytes exceeds the limit of %d bytes", jsonFilePath, info.Size(), commons.MaxJSONFileSizeBytes)
	}

	return nil
}

func GetStringFromJSONFile(jsonFilePath string) (string, error) {
	err := ValidateJSONAndFilePath(jsonFilePath)
	if err != nil {
		return "", fmt.Errorf("JSON file validation failed for file '%s': %w", jsonFilePath, err)
	}

	fileContent, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read JSON file '%s': %w", jsonFilePath, err)
	}

	if !json.Valid(fileContent) {
		return "", fmt.Errorf("file '%s' does not contain valid JSON", jsonFilePath)
	}

	var compacted bytes.Buffer
	if err := json.Compact(&compacted, fileContent); err != nil {
		return "", fmt.Errorf("failed to compact JSON from file '%s': %w", jsonFilePath, err)
	}
	value := compacted.String()
	logger.Info("Compacted JSON value: ", value)

	return value, nil
}

func StoreJSONInFile(data []byte, outputPath string) error {
	if !json.Valid(data) {
		return fmt.Errorf("data is not valid JSON, cannot store in file '%s'", outputPath)
	}

	var pretty bytes.Buffer
	if err := json.Indent(&pretty, data, "", "  "); err != nil {
		return fmt.Errorf("failed to format JSON data for file '%s': %w", outputPath, err)
	}

	if err := os.WriteFile(outputPath, pretty.Bytes(), 0o644); err != nil {
		return fmt.Errorf("failed to write JSON to file '%s': %w", outputPath, err)
	}
	fmt.Println("JSON value written to", outputPath)
	return nil

}

func ParseCommandLine(s string) []string {
	var args []string
	var current strings.Builder
	inSingleQuote := false
	inDoubleQuote := false
	hasChar := false

	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
		} else if c == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
		} else if (c == ' ' || c == '\t') && !inSingleQuote && !inDoubleQuote {
			if hasChar {
				args = append(args, current.String())
				current.Reset()
				hasChar = false
			}
		} else {
			current.WriteByte(c)
			hasChar = true
		}
	}
	if hasChar {
		args = append(args, current.String())
	}
	return args
}
