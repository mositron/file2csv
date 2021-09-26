package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var VERSION = "0.1.1"

func main() {

	log.SetFlags(0)
	LoadIni()

	var files, file []string
	var err error

	for i := range Conf.path {
		file, err = WalkMatch(Conf.path[i])
		fmt.Println(" - ", Conf.path[i])
		if err == nil {
			files = append(files, file...)
		}
	}

	txt := []string{"filename,artist,title,description,keywords"}
	for i := range files {
		txt = append(txt, files[i]+","+Conf.artist+","+Conf.title+","+Conf.description+",\""+Conf.keywords+"\"")
	}
	fmt.Println(txt)

	if err := writeLines(txt, "./result.csv"); err != nil {
		log.Fatalf("writeLines: %s", err)
	}
}

func WalkMatch(root string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		fmt.Println(" - ", path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		for i := range Conf.name {
			if matched, err := filepath.Match(Conf.name[i], filepath.Base(path)); err != nil {
				fmt.Println(" = ", filepath.Base(path), Conf.name[i], err)
				//return err
			} else if matched {
				fmt.Println(" + ", filepath.Base(path), Conf.name[i])
				matches = append(matches, filepath.Base(path))
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}
func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}