package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

func GetDirectorySize(path string) (size int64) {
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		size += info.Size()
		return nil
	})
	return
}

func main() {
	days, err := strconv.Atoi(os.Args[1])

	if err != nil {
		log.Fatalln("invalid argument", os.Args[1])
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}

	dir := path.Join(usr.HomeDir, "/.config/direnv/store/")
	last_modified_limit := time.Now().AddDate(0, 0, -days)

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	var cleaned_up int64
	var cleaned_up_size int64
	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())

		// check if entry is a symlink
		if entry.IsDir() {
			continue
		}

		if (entry.Mode() & os.ModeSymlink) == 0 {
			continue
		}

		// check if sha256 path match
		direnvPath, err := os.Readlink(path)
		if err != nil {
			log.Fatalln(err)
		}

		hash := sha256.Sum256([]byte(direnvPath))
		if hex.EncodeToString(hash[:]) != entry.Name() {
			continue
		}

		// check if last modified time is before want we want
		if entry.ModTime().After(last_modified_limit) {
			continue
		}

		direnv, err := os.Stat(direnvPath)
		if err != nil && !os.IsNotExist(err) {
			log.Fatalln(err)
		}

		// we have two cases here
		// 1. target exists and is a folder - we delete the environment
		// 2. target didn't exist - we delete the link

		if err == nil && direnv.IsDir() {
			cleaned_up_size += GetDirectorySize(direnvPath)
			cleaned_up += 1

			// delete folder
			fmt.Printf("Removing %s (%s)\n", direnvPath, entry.Name())
			os.RemoveAll(direnvPath)
		}

		// delete link
		fmt.Printf("Removing link %s\n", entry.Name())
		os.Remove(path)
	}

	fmt.Printf("Cleaned up %d environments, saving a total of %dmb!", cleaned_up, cleaned_up_size/1024/1024)
}
