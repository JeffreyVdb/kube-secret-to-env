package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// Based on work by kelseyhightower
// https://github.com/kelseyhightower/conf2kube

// OutputType is the output type of a key value pair
var OutputType string

// Secret holds a Kubernetes secret.
type Secret struct {
	APIVersion string                 `json:"apiVersion"`
	Data       map[string]string      `json:"data"`
	Kind       string                 `json:"kind"`
	Metadata   map[string]interface{} `json:"metadata"`
	Type       string                 `json:"type"`
}

// KeyValueWriter is an interface describing how a key value is written to stdout
type KeyValueWriter interface {
	WriteKV(key, value string)
}

// EnvKeyValueWriter is a type that can write key value in "{key}={value}" form
type EnvKeyValueWriter struct {
	writer io.Writer
}

// WriteKV writes a key, value pair in "{key}={value}" form
func (kvWriter EnvKeyValueWriter) WriteKV(key, value string) {
	fmt.Fprintf(kvWriter.writer, "%s=%s\n", key, value)
}

// ExportKeyValueWriter is a type that writes a key, value pair in shell export syntax
type ExportKeyValueWriter struct {
	writer io.Writer
}

// WriteKV writes a key, value pair in form "export {key}={value}"
func (kvWriter ExportKeyValueWriter) WriteKV(key, value string) {
	fmt.Fprintf(kvWriter.writer, "export %s=\"%s\"\n", key, strings.Replace(value, "\"", "\\\"", -1))
}

func init() {
	flag.StringVar(&OutputType, "type", "env", "The output type, can be [env, shell]")
}

func main() {
	flag.Parse()

	var kubernetesSecret Secret
	decoder := json.NewDecoder(os.Stdin)
	if err := decoder.Decode(&kubernetesSecret); err != nil {
		fmt.Fprintf(os.Stderr, "Error occurred while reading stdin: %v\n", err)
		os.Exit(1)
	}

	var outputWriter KeyValueWriter
	switch OutputType {
	case "env":
		outputWriter = EnvKeyValueWriter{writer: os.Stdout}
	case "shell":
		outputWriter = ExportKeyValueWriter{writer: os.Stdout}
	default:
		fmt.Fprintf(os.Stderr, "%s is not a supported output type\n", OutputType)
		os.Exit(1)
	}

	for k, v := range kubernetesSecret.Data {
		decodedVal, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred while decoding key %s: %v\n", k, err)
			os.Exit(1)
		}
		outputWriter.WriteKV(k, string(decodedVal))
	}
}
