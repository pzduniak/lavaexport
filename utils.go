package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/bgentry/speakeasy"
)

func prompt(input string, secure bool) (string, error) {
	br := bufio.NewReader(os.Stdin)
	var sl string

	for {
		if !secure {
			fmt.Print(input)

			line, _, err := br.ReadLine()
			if err != nil {
				return "", err
			}

			sl = strings.TrimSpace(string(line))
		} else {
			var err error
			sl, err = speakeasy.Ask(input)
			if err != nil {
				return "", err
			}
		}

		fmt.Print("Is the information above correct? [y/n]: ")

		line, _, err := br.ReadLine()
		if err != nil {
			return "", err
		}

		if len(line) < 1 || line[0] != byte('y') {
			continue
		}

		break
	}

	return sl, nil
}

func saveStruct(path string, val interface{}) error {
	data, err := json.MarshalIndent(val, "", "    ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(path, data, 0666); err != nil {
		return err
	}

	return nil
}
