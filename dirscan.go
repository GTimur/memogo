package memogo

import (
	"log"
	"path/filepath"
	"strings"
)

//search files in folder by selected masks
func FindFiles(dir string, mask []string) (files map[string]string, err error) {
	var list []string
	files = make(map[string]string)

	for i := range mask {
		list, err = filepath.Glob(dir + "\\" + strings.ToUpper(mask[i]))
		if err != nil {
			log.Println("findFiles error: ", err)
			return nil, err
		}
		//files = append(files, list...)
		for _, f := range list {
			files[f] = mask[i]
		}
	}
	for i := range mask {
		list, err = filepath.Glob(dir + "\\" + strings.ToLower(mask[i]))
		if err != nil {
			log.Println("findFiles error: ", err)
			return nil, err
		}
		//files = append(files, list...)
		for _, f := range list {
			files[f] = mask[i]
		}
	}
	return files, err
}

//FindAllFiles - search files in all subdirectories by selected masks
func FindAllFiles(rootdir string, mask []string) (files map[string]string, err error) {
	dirs := make(map[string]string)
	files = make(map[string]string)

	dirs, err = FindFiles(rootdir, []string{"*"})
	if err != nil {
		log.Fatalf("FindAllFiles: FindAllFiles error: %v", err)
	}

	for k := range dirs {
		f, err := FindFiles(k, []string{"*.*"})
		if err != nil {
			log.Fatalf("FindAllFiles: FindAllFiles error: %v", err)
		}

		for kk, vv := range f {
			files[kk] = vv
		}
	}
	return files, err
}
