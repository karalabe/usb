package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var dummy = `
// +build dummy

// This Go file is part of a workaround for "go mod vendor".
package %q
`

func updateHidApi() error {
	if err := os.Rename("./hidapi", "./hidapi_old"); err != nil {
		return err
	}
	if _, err := exec.Command("git", "clone", "https://github.com/libusb/hidapi").Output(); err != nil {
		return err
	}
	os.Chdir("./hidapi")
	var version string
	// Parse commit
	if v, err := exec.Command("git", "rev-parse", "HEAD").Output(); err != nil {
		return err
	} else {
		version = string(v)
	}
	fmt.Printf("Version: %v\n", version)
	// Strip git metadata
	if err := os.RemoveAll(".git"); err != nil {
		return err
	}
	// traverse folders, dump dummy.go-file in them
	filepath.Walk("./", func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}
		if strings.HasPrefix(path, ".") {
			return nil
		}
		if strings.Contains(path, "test") {
			return nil
		}

		fmt.Printf("walking path %v : %v\n", path, info.Name())
		return os.WriteFile(fmt.Sprintf("%v/dummy.go", path),
			[]byte(fmt.Sprintf(dummy, info.Name())), 0777)
	})
	os.WriteFile("dummy.go", []byte(fmt.Sprintf(dummy, "hidapi")), 0777)
	os.Chdir("../")
	os.WriteFile("hidapi_version.txt", []byte(version), 0777)
	return nil
}

func updateLibusb() error {
	if err := os.Rename("./libusb", "./libusb_old"); err != nil {
		return err
	}
	if _, err := exec.Command("git", "clone", "https://github.com/libusb/libusb").Output(); err != nil {
		return err
	}
	os.Chdir("./libusb")
	var version string
	// Parse commit
	if v, err := exec.Command("git", "rev-parse", "HEAD").Output(); err != nil {
		return err
	} else {
		version = string(v)
	}
	fmt.Printf("Version: %v\n", version)
	// Strip git metadata
	if err := os.RemoveAll(".git"); err != nil {
		return err
	}
	// traverse folders, dump dummy.go-file in them
	filepath.Walk("./", func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}
		if strings.HasPrefix(path, ".") {
			return nil
		}
		if strings.Contains(path, "test") {
			return nil
		}

		fmt.Printf("walking path %v : %v\n", path, info.Name())
		return os.WriteFile(fmt.Sprintf("%v/dummy.go", path),
			[]byte(fmt.Sprintf(dummy, info.Name())), 0777)
	})
	os.WriteFile("dummy.go", []byte(fmt.Sprintf(dummy, "libusb")), 0777)
	os.Chdir("../")
	os.WriteFile("libusb_version.txt", []byte(version), 0777)
	return nil
}

func main() {
	//if err := updateHidApi(); err != nil {
	//	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	//	os.Exit(1)
	//}
	if err := updateLibusb(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

//update_hidapi() {
//mv ./hidapi ./hidapi_old #stash away the old code
//git clone https://github.com/libusb/hidapi #clone the new
//( cd hidapi
//ver=$(git rev-parse HEAD) # remember git commit
//rm -rf .git/              # remove git metadata
//# traverse folders, dump a dummy.go-file in them
//ls -d */ | xargs -I {} bash -c "cp ../hidapi_old/dummy.go '{}'"
//)
//echo "hidapi at version $ver" >> ./hidapi.version
//}
//
//update_hidapi
