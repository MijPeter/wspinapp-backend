package test_utils

import (
	"log"
	"os"
)

// Golden If golden already exists then it is loaded and returned.
// Otherwise, create new golden and return its value.
func Golden(name string, actual string) string {
	name = "goldens/" + name + ".json"
	if fileExists(name) {
		return string(load(name))
	} else {
		save(name, []byte(actual))
		return actual
	}
}

func load(name string) []byte {
	dat, err := os.ReadFile(name)
	check(err)
	return dat
}

func save(name string, data []byte) {
	f, err := os.Create(name)
	check(err)
	defer f.Close()

	_, err = f.Write(data)
	check(err)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func check(e error) {
	if e != nil {
		log.Fatal("Error in test utils: " + e.Error())
	}
}
