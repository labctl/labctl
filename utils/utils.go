package utils

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// Log yaml data with optional message
func LogYAML(data interface{}, msg ...string) {
	outp, err := yaml.Marshal(data)
	if err != nil {
		log.Errorf("error: %v", err)
	}
	if len(msg) == 0 {
		fmt.Printf("%s\n", outp)
	} else {
		fmt.Printf("%s\n%s\n", strings.Join(msg, ", "), outp)
	}
}

// Write YAML to a file (or stdout if no filename)
func WriteYAML(filename string, data interface{}) error {
	out, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("invalid result: %v", err)
	}
	if filename == "" {
		fmt.Printf("%s\n", out)
	} else {
		err := os.WriteFile(filename, out, 0o644)
		if err != nil {
			return err
		}
		log.Infof("Written to %s", filename)
	}
	return nil
}

// Ensure the input is a valid 'dict'
func Mapify(i interface{}) (map[string]interface{}, bool) {
	value := reflect.ValueOf(i)
	if value.Kind() == reflect.Map {
		m := map[string]interface{}{}
		for _, k := range value.MapKeys() {
			m[fmt.Sprintf("%v", k)] = value.MapIndex(k).Interface()
		}
		return m, true
	}
	return map[string]interface{}{}, false
}

// Ensure the input is an array of 'dicts'
func ArrayMapify(in interface{}) ([]map[string]interface{}, error) {
	res := []map[string]interface{}{}
	items := reflect.ValueOf(in)
	for i := 0; i < items.Len(); i++ {
		item := items.Index(i).Interface()
		mapitem, ok := Mapify(item)
		if !ok {
			return nil, fmt.Errorf("expected a dictionary like map[string]interface{}, not %v", item)
		}
		res = append(res, mapitem)
	}
	return res, nil
}

// Does a slice contain a value?
func Contains[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

// Partition a string similar to Python's partition. Returns (before, sep, after)
func Partition(s string, sep string) (string, string, string) {
	parts := strings.SplitN(s, sep, 2)
	if len(parts) == 1 {
		return parts[0], "", ""
	}
	return parts[0], sep, parts[1]
}

// Removes duplicates from a slice
func Unique[T comparable](elems []T) []T {
	inResult := make(map[T]bool)
	var result []T
	for _, str := range elems {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return result
}

func Must[T comparable](v T, err error) T {
	if err != nil {
		log.Fatalln(err)
	}
	return v
}

// Check and warn if root
func CheckRoot() {
	if os.Getegid() == 0 {
		log.Warn("You are running as root user.")
	}
}

// Try find a local topo file
func EnsureTopo(topo string) (string, error) {
	if topo != "" {
		return topo, nil
	}
	files, err := filepath.Glob("*.clab.y*ml")
	if err != nil {
		return "", fmt.Errorf("could not find local clab files: %v", err)
	}
	if len(files) == 0 {
		return "", fmt.Errorf("no local topo files")
	}
	if len(files) > 1 {
		return "", fmt.Errorf("multiple topology files found: %v", strings.Join(files, ", "))
	}
	return files[0], nil
}
