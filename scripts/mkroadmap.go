// Package main provides a tool to generate a Markdown-formatted API support status document
// for PagerDuty's REST API implementation. It processes the OpenAPI specification and creates
// a structured overview of available endpoints grouped by tags.
//
// The tool can read API specifications either from a local file or directly from PagerDuty's
// API reference URL. It generates a formatted document that shows the implementation status
// of different HTTP methods (GET, POST, PUT, DELETE, PATCH) for each endpoint.
//
// Usage:
//
//	mkroadmap [-i input_file] [-o output_file] [-header=bool]
//
// Flags:
//
//	-i string
//		Input file path (defaults to fetching from PagerDuty API if not specified)
//	-o string
//		Output file path (defaults to "website/docs/guides/pagerduty_api_support_status.html.markdown")
//	-header
//		Include file header in output (default true)
//
// The output format includes:
//   - Grouping of endpoints by tags
//   - Color-coded HTTP methods
//   - Markdown-compatible checkboxes for tracking implementation status
//   - Escaped path parameters for proper Markdown rendering
//
// The generated document serves as both a reference and a roadmap for
// PagerDuty API implementation status.
package main

// Run something like this to include jira cloud and slack connections
// ```
// go run ./scripts/mkroadmap.go
// go run ./scripts/mkroadmap.go -header=0 -i ~/Downloads/jiracloud.json -o - >> website/docs/guides/pagerduty_api_support_status.html.markdown
// go run ./scripts/mkroadmap.go -header=0 -i ~/Downloads/slack.json -o - >> website/docs/guides/pagerduty_api_support_status.html.markdown
// ```

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"sort"
	"strings"
)

type List = []any
type Map = map[string]any
type Set = map[string]struct{}

const (
	defaultOutputFilename = "website/docs/guides/pagerduty_api_support_status.html.markdown"
	apiReferenceURL       = "https://stoplight.io/api/v1/projects/pagerduty-upgrade/api-schema/nodes/reference/REST/openapiv3.json?fromExportButton=true&snapshotType=http_service&deref=optimizedBundle"
)

var flagInFilename, flagOutFilename string
var flagHeaderWrite bool

func main() {
	flag.StringVar(&flagOutFilename, "o", defaultOutputFilename, "output file")
	flag.StringVar(&flagInFilename, "i", "", "input file")
	flag.BoolVar(&flagHeaderWrite, "header", true, "include file header in output")
	flag.Parse()

	if flagInFilename == "" {
		args := flag.Args()
		if len(args) > 0 {
			flagInFilename = args[0]
		}
	}

	var fOutput *os.File
	if flagOutFilename == "" || flagOutFilename == "-" {
		fOutput = os.Stdout
	} else {
		f, err := os.Create(flagOutFilename)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		fOutput = f
	}

	var r io.ReadCloser
	if flagInFilename == "" {
		res, err := http.Get(apiReferenceURL)
		if err != nil {
			panic(err)
		}
		r = res.Body
	} else if flagInFilename == "-" {
		r = os.Stdin
	} else {
		f, err := os.Open(flagInFilename)
		if err != nil {
			panic(err)
		}
		r = f
	}

	var apiDefinition Map
	if err := json.NewDecoder(r).Decode(&apiDefinition); err != nil {
		r.Close()
		panic(err)
	}
	r.Close()

	db := make(endpointsByTag)

	paths, _ := apiDefinition["paths"].(Map)
	for pKey, p := range paths {
		cleanPath := strings.ReplaceAll(pKey, `{type_id_or_name}`, `{id}`)
		cleanPath = strings.ReplaceAll(string(cleanPath), `_`, `\_`)

		for opKey, op := range p.(Map) {
			switch opKey {
			case "get", "post", "put", "patch", "delete":
				operation, _ := op.(Map)
				tags, _ := operation["tags"].(List)
				for _, t := range tags {
					tag, ok := t.(string)
					if !ok {
						continue
					}

					if v, ok := db[tag]; ok && v != nil {
						db[tag] = v
					} else {
						db[tag] = make(map[string]Set)
					}

					if v, ok := db[tag][cleanPath]; ok && v != nil {
						db[tag][cleanPath] = v
					} else {
						db[tag][cleanPath] = make(Set)
					}

					db[tag][cleanPath][opKey] = struct{}{}
				}
			}
		}
	}

	GenerateOut(fOutput, db)
}

type endpoint map[string]Set
type endpointsByTag map[string]endpoint

func GenerateOut(fOutput *os.File, db endpointsByTag) {
	tagKeys := make([]string, len(db))
	i := 0
	for k := range db {
		tagKeys[i] = k
		i++
	}
	sort.Strings(tagKeys)

	if flagHeaderWrite {
		fmt.Fprint(fOutput, `---
page_title: "PagerDuty API Support Status"
---

# PagerDuty API Support Status

This document tracks the implementation status of PagerDuty operations REST API across the Resources and Data Sources of the Terraform Provider.
It serves as both a quick reference for available functionality and a development roadmap.
The status is updated as new features are implemented and tested.
Features marked with a “prohibited 🚫” emoji are not relevant for this provider.
`)
	}

	styleFn := func(m string) string {
		s := "position:relative;top:-2px;color:white;font-family:monospace;padding:3px;font-size:.7em;border-radius:4px;background-color:"
		switch m {
		case "get":
			s += "#05b870"
		case "post":
			s += "#19abff"
		case "delete":
			s += "#f05151"
		case "put":
			s += "#f46d2a"
		case "patch":
			s += "#f46d2a"
		default:
			s += "black"
		}
		return s + ";}"
	}

	for _, tag := range tagKeys {
		fmt.Fprintln(fOutput, "\n-", tag)
		endp := db[tag]

		paths := make([]string, len(endp))
		i := 0
		for k := range endp {
			paths[i] = k
			i++
		}
		sort.Strings(paths)

		for _, path := range paths {
			options := endp[path]

			i := 0
			optKeys := make([]string, len(options))
			for k := range options {
				optKeys[i] = k
				i++
			}

			slices.SortFunc(optKeys, methodSortFunc)

			for _, key := range optKeys {
				method := key
				fmt.Fprintf(fOutput, "    - [ ]")
				fmt.Fprintf(
					fOutput, ` <span style="%s">%s</span>`,
					styleFn(method), strings.ToUpper(method),
				)
				fmt.Fprintf(fOutput, ` <span style="font-size:.9em">%s</span>`+"\n", path)
			}
		}
	}
}

var methodOrder = map[string]int{
	"get": 0, "post": 1, "put": 2, "delete": 3,
}

func methodSortFunc(a, b string) int {
	na, ok := methodOrder[a]
	if !ok {
		na = 4
	}

	nb, ok := methodOrder[b]
	if !ok {
		nb = 4
	}

	return na - nb
}
