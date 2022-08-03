package utils

import (
	"fmt"
	"io/ioutil"
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
		err := ioutil.WriteFile(filename, out, 0o644)
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
