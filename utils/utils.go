package utils

import (
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func LogPretty(msg string, val interface{}) {
	outp, err := yaml.Marshal(val)
	if err != nil {
		log.Errorf("error: %v", err)
	}
	fmt.Printf("msg\n%s\n", outp)
}

func LoadFile(filen string) string {
	file, err := os.Open(filen)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	return string(b)
}
