package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/paulmach/orb/encoding/mvt"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		fmt.Println("Usage: pbftile decode")
		fmt.Println("Reads a .pbf map tile from stdin and outputs pretty-printed JSON")
		fmt.Println("Example: pmtiles tile ./static/colorado.pmtiles 12 849 1550 | pbftile decode")
		return
	}

	// Read the entire input from stdin
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
		os.Exit(1)
	}

	if len(data) == 0 {
		fmt.Fprintf(os.Stderr, "No data received from stdin\n")
		os.Exit(1)
	}

	// Check if data is gzipped and decompress if necessary
	data, err = handleGzipData(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error handling gzipped data: %v\n", err)
		os.Exit(1)
	}

	// Parse the MVT tile
	layers, err := mvt.Unmarshal(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing MVT tile: %v\n", err)
		os.Exit(1)
	}

	// Convert to a JSON-serializable format
	result := make(map[string]interface{})

	for _, layer := range layers {
		layerData := map[string]interface{}{
			"features": make([]map[string]interface{}, 0),
		}

		for _, feature := range layer.Features {
			featureData := map[string]interface{}{
				"type":       "Feature",
				"geometry":   feature.Geometry,
				"properties": feature.Properties,
			}

			layerData["features"] = append(layerData["features"].([]map[string]interface{}), featureData)
		}

		result[layer.Name] = layerData
	}

	// Pretty print as JSON
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling to JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonData))
}

// handleGzipData checks if the data is gzipped and decompresses it if necessary
func handleGzipData(data []byte) ([]byte, error) {
	// Check for gzip magic number (0x1f, 0x8b)
	if len(data) < 2 {
		return data, nil
	}

	if data[0] == 0x1f && data[1] == 0x8b {
		// Data is gzipped, decompress it
		reader, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, fmt.Errorf("failed to create gzip reader: %v", err)
		}
		defer reader.Close()

		decompressed, err := io.ReadAll(reader)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress gzipped data: %v", err)
		}

		return decompressed, nil
	}

	// Data is not gzipped, return as is
	return data, nil
}
