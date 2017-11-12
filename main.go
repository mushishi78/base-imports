package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("cannot get current working directory: %v\n", err)
	}

	re := regexp.MustCompile(`\sfrom\s["']\..+["']`)

	filepath.Walk(baseDir, func(path string, f os.FileInfo, err error) error {

		if strings.Contains(path, "node_modules") {
			return nil
		}
		if !strings.HasSuffix(path, ".ts") && !strings.HasSuffix(path, ".tsx") {
			return nil
		}

		dir := filepath.Dir(path[len(baseDir)+1:])

		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}

		text := string(bytes)

		for re.MatchString(text) {
			indexes := re.FindStringIndex(text)
			start := indexes[0] + 7
			end := indexes[1] - 1
			text = text[:start] +
				filepath.Join(dir, text[start:end]) +
				text[end:]
		}

		err = ioutil.WriteFile(path, []byte(text), 0777)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	})
}
