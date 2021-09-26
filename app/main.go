package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var VERSION = "0.1.4"

func main() {
	var cur_path string
	ex, er := os.Executable()
	if er == nil {
		cur_path = filepath.Dir(ex)
	} else {
		exReal, er := filepath.EvalSymlinks(ex)
		if er != nil {
			panic(er)
		}
		cur_path = filepath.Dir(exReal)
	}

	fmt.Println("Directory:", cur_path)

	log.SetFlags(0)
	LoadIni(cur_path)

	var files, file []string

	for i := range Conf.path {
		file, er = WalkMatch(Conf.path[i])
		fmt.Println(" - ", Conf.path[i])
		if er == nil {
			files = append(files, file...)
		}
	}

	txt := []string{"filename,artist,title,description,keywords"}
	for i := range files {
		txt = append(txt, files[i]+","+Conf.artist+","+Conf.title+","+Conf.description+",\""+Conf.keywords+"\"")
	}
	if err := writeLines(txt, cur_path+"/keywords.csv"); err != nil {
		log.Fatalf("writeLines: %s", err)
	} else {
		fmt.Println("finished - ", len(txt)-1, "files")
	}
}

func WalkMatch(root string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		//fmt.Println(" - ", path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		for i := range Conf.name {
			if matched, err := filepath.Match(Conf.name[i], filepath.Base(path)); err != nil {
				//fmt.Println(" = ", filepath.Base(path), Conf.name[i], err)
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
