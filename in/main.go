package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"


	"github.com/DennisDenuto/boshver-resource/models"
	"github.com/DennisDenuto/boshver-resource/version"
	"strings"
	"strconv"
	"errors"
)
const (
	numbers  string = "0123456789"
	alphas          = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-"
	alphanum        = alphas + numbers
)


func main() {
	if len(os.Args) < 2 {
		println("usage: " + os.Args[0] + " <destination>")
		os.Exit(1)
	}

	destination := os.Args[1]

	err := os.MkdirAll(destination, 0755)
	if err != nil {
		fatal("creating destination", err)
	}

	var request models.InRequest
	err = json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		fatal("reading request", err)
	}

	inputVersion, err := Parse(request.Version.Number)
	if err != nil {
		fatal("parsing semantic version", err)
	}

	bumped := version.BumpFromParams(request.Params.Bump).Apply(inputVersion)

	if !bumped.Equals(inputVersion) {
		fmt.Fprintf(os.Stderr, "bumped locally from %s to %s\n", inputVersion, bumped)
	}

	versionFileNames := []string{"number", "version"}

	for _, fileName := range versionFileNames {
		numberFile, err := os.Create(filepath.Join(destination, fileName))
		if err != nil {
			fatal("opening number file", err)
		}

		defer numberFile.Close()

		_, err = fmt.Fprintf(numberFile, "%s", bumped.String())
		if err != nil {
			fatal("writing to number file", err)
		}
	}

	json.NewEncoder(os.Stdout).Encode(models.InResponse{
		Version: request.Version,
		Metadata: models.Metadata{
			{"number", request.Version.Number},
		},
	})
}

func Parse(s string) (version.BoshVersion, error) {
	if len(s) == 0 {
		return version.BoshVersion{}, errors.New("Version string empty")
	}

	// Split into major.minor
	parts := strings.SplitN(s, ".", 2)
	if len(parts) != 3 {
		return version.BoshVersion{}, errors.New("No Major.Minor elements found")
	}

	// Major
	if !containsOnly(parts[0], numbers) {
		return version.BoshVersion{}, fmt.Errorf("Invalid character(s) found in major number %q", parts[0])
	}
	if hasLeadingZeroes(parts[0]) {
		return version.BoshVersion{}, fmt.Errorf("Major number must not contain leading zeroes %q", parts[0])
	}
	major, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return version.BoshVersion{}, err
	}

	// Minor
	if !containsOnly(parts[1], numbers) {
		return version.BoshVersion{}, fmt.Errorf("Invalid character(s) found in minor number %q", parts[1])
	}
	if hasLeadingZeroes(parts[1]) {
		return version.BoshVersion{}, fmt.Errorf("Minor number must not contain leading zeroes %q", parts[1])
	}
	minor, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return version.BoshVersion{}, err
	}

	v := version.BoshVersion{}
	v.Major = major
	v.Minor = minor

	return v, nil
}

func containsOnly(s string, set string) bool {
	return strings.IndexFunc(s, func(r rune) bool {
		return !strings.ContainsRune(set, r)
	}) == -1
}

func hasLeadingZeroes(s string) bool {
	return len(s) > 1 && s[0] == '0'
}



func fatal(doing string, err error) {
	println("error " + doing + ": " + err.Error())
	os.Exit(1)
}

