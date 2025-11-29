package parser

import (
	"bufio"
	"io"
	"strings"
)

// Parser parses Vim profile output files
type Parser struct {
	profiles []*Profile
}

// NewParser creates a new Parser instance
func NewParser() *Parser {
	return &Parser{
		profiles: make([]*Profile, 0),
	}
}

// parseScriptHeader parses the SCRIPT header section
func (p *Parser) parseScriptHeader(reader io.Reader) (*Script, error) {
	scanner := bufio.NewScanner(reader)
	var scriptPath string

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "SCRIPT") {
			scriptPath = p.extractScriptPath(line)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	script := &Script{
		Path:      scriptPath,
		Lines:     make(map[int]*Line),
		Functions: make(map[string]*Function),
	}

	return script, nil
}

// extractScriptPath extracts the script path from a SCRIPT line
func (p *Parser) extractScriptPath(line string) string {
	// Format: "SCRIPT  /path/to/script.vim" or "SCRIPT /path/to/script.vim"
	parts := strings.Fields(line)
	if len(parts) >= 2 {
		return parts[1]
	}
	return ""
}
