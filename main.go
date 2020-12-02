package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"errors"
	"github.com/spf13/cobra"
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

func Cleanup(days int, dir string, dry bool) error {
	mode, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("Direnv store (%s) does not exist, create some direnvs!", dir))
	}

	if !mode.IsDir() {
		return errors.New(fmt.Sprintf("Direnv store (%s) exists but is not a directory, aborting.", dir))
	}

	last_modified_limit := time.Now().AddDate(0, 0, -days)

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
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
			return err
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
			return err
		}

		// we have two cases here
		// 1. target exists and is a folder - we delete the environment
		// 2. target didn't exist - we delete the link

		if err == nil && direnv.IsDir() {
			cleaned_up_size += GetDirectorySize(direnvPath)
			cleaned_up += 1

			// delete folder
			if !dry {
				fmt.Printf("Removing %s (%s)\n", direnvPath, entry.Name())
				os.RemoveAll(direnvPath)
			} else {
				fmt.Printf("Would remove %s (%s)\n", direnvPath, entry.Name())
			}
		}

		// delete link
		if !dry {
			os.Remove(path)
		}
	}

	saved_mb := cleaned_up_size / 1024 / 1024
	if !dry {
		fmt.Printf("Cleaned up %d environments, saving a total of %dmb!\n", cleaned_up, saved_mb)
	} else {
		fmt.Printf("Would clean up %d environments, saving a total of %dmb!\n", cleaned_up, saved_mb)
	}
	return nil
}

func PrintHook() error {
	hook, err := Asset("shell/hook.sh")
	if err != nil {
		return errors.New("Failed to find hook, maybe it wasn't included in the build...")
	}
	fmt.Printf(string(hook))
	return nil
}

func GetStoreDirectory() (string, error) {
	if storeDirectory := os.Getenv("DIRENV_STORE"); storeDirectory != "" {
		return storeDirectory, nil
	}

	dataDirectory := os.Getenv("XDG_DATA_HOME")
	if dataDirectory == "" {
		homeDirectory, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dataDirectory = path.Join(homeDirectory, ".local/share/")
	}

	return path.Join(dataDirectory, "direnv/store"), nil
}

func main() {
	var days int
	var storePath string
	var dryRun bool

	storeDirectory, err := GetStoreDirectory()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	rootCmd := &cobra.Command{
		Use:   "direnv-gc",
		Short: "cleans up unused direnvs",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Cleanup(days, storePath, dryRun)
		},
	}

	rootCmd.Flags().IntVarP(&days, "days", "d", 10, "number of days to keep")
	rootCmd.Flags().BoolVarP(&dryRun, "dry-run", "", false, "dry run")
	rootCmd.Flags().StringVarP(&storePath, "store-path", "", storeDirectory, "path to store where all the direnvs are linked to")

	rootCmd.AddCommand(&cobra.Command{
		Use:   "hook",
		Short: "prints the direnv hook",
		RunE: func(cmd *cobra.Command, args []string) error {
			return PrintHook()
		},
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
