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
