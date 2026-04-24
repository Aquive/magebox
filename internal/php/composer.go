package php

import (
	"encoding/json"
	"os"
	"regexp"
)

// composerJSON is a minimal representation of composer.json for PHP version detection.
type composerJSON struct {
	Require map[string]string `json:"require"`
	Config  struct {
		Platform map[string]string `json:"platform"`
	} `json:"config"`
}

// DetectVersionFromComposer reads composer.json at the given path and returns
// the PHP version it targets. It prefers an exact pinned version from
// config.platform.php, falling back to parsing a constraint from require.php
// (e.g. "~8.3", "^8.2", ">=8.2", "8.2.*"). The returned version is normalized
// to "MAJOR.MINOR" and is only returned when it matches one of SupportedVersions.
// An empty string is returned when nothing usable is found.
func DetectVersionFromComposer(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	var c composerJSON
	if err := json.Unmarshal(data, &c); err != nil {
		return ""
	}

	if v := extractMajorMinor(c.Config.Platform["php"]); isSupportedVersion(v) {
		return v
	}

	if c.Require != nil {
		if v := extractMajorMinor(c.Require["php"]); isSupportedVersion(v) {
			return v
		}
	}

	return ""
}

// majorMinorRe matches the first "MAJOR.MINOR" occurrence in a version string
// or Composer constraint. Handles values like "8.3", "8.3.10", "~8.3",
// "^8.2", ">=8.2", "8.2.*", ">=8.2 <8.4".
var majorMinorRe = regexp.MustCompile(`(\d+)\.(\d+)`)

func extractMajorMinor(v string) string {
	m := majorMinorRe.FindStringSubmatch(v)
	if len(m) < 3 {
		return ""
	}
	return m[1] + "." + m[2]
}

func isSupportedVersion(v string) bool {
	if v == "" {
		return false
	}
	for _, sv := range SupportedVersions {
		if sv == v {
			return true
		}
	}
	return false
}
