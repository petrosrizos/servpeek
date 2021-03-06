// Package file allows to query various properties of operating system files
package file

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"time"

	"github.com/milosgajdos83/servpeek/utils"
)

// Bit mask for regular files
const ModeRegular = ^os.ModeType

func withFileReader(f Filer, fn func(r io.Reader) error) error {
	file, err := os.Open(f.Path())
	if err != nil {
		return err
	}
	defer file.Close()
	return fn(file)
}

func isTrueOrError(ok bool, err error) error {
	if ok {
		return nil
	}
	return err
}

// IsRegular checks if provided file is a regular file
// It returns error if os.Stat returns error
func IsRegular(f Filer) error {
	return IsMode(f, ModeRegular)
}

// IsDirectory checks if provided file is a directory
// It returns error if os.Stat returns error
func IsDirectory(f Filer) error {
	return IsMode(f, os.ModeDir)
}

// IsBlockDevice checks if provided file is a block device file
// It returns error if os.Stat returns error
func IsBlockDevice(f Filer) error {
	return IsMode(f, os.ModeDevice)
}

// IsCharDevice checks if provided file is a character device file
// It returns error if os.Stat returns error
func IsCharDevice(f Filer) error {
	return IsMode(f, os.ModeCharDevice)
}

// IsPipe checks if provided file is a named pipe file
// It returns error if os.Stat returns error
func IsPipe(f Filer) error {
	return IsMode(f, os.ModeNamedPipe)
}

// IsSocket checks if provided file is a socket
// It returns error if os.Stat returns error
func IsSocket(f Filer) error {
	return IsMode(f, os.ModeSocket)
}

// IsSymlink checks if provided file is a symbolic link
// It returns error if os.Stat returns error
func IsSymlink(f Filer) error {
	return IsMode(f, os.ModeSymlink)
}

var verboseNames = map[os.FileMode]string{
	os.ModeSymlink:   "symlink",
	os.ModeSocket:    "socket",
	os.ModeDir:       "directory",
	os.ModeNamedPipe: "named pipe",
	ModeRegular:      "regular file",
}

// IsMode checks if the provided file has the same mode as the one passed in via paramter
// It returns error if os.Stat returns error
func IsMode(f Filer, mode os.FileMode) error {
	fi, err := f.Info()
	if err != nil {
		return err
	}

	if verboseMode, ok := verboseNames[mode]; ok {
		err = fmt.Errorf("%s file is not a %s", f, verboseMode)
	} else {
		err = fmt.Errorf("%s file is not mode %v", f, mode)
	}
	return isTrueOrError(fi.Mode()&mode != 0, err)
}

// IsOwnedBy checks if the provided file is owned by username user
// It returns an error if os.Stat returns error
func IsOwnedBy(f Filer, username string) error {
	return isOwnedBy(f, username)
}

// IsGrupedInto checks if the provided file is owned by groupname group
// It returns an error if os.Stat returns error
func IsGrupedInto(f Filer, groupname string) error {
	return isGrupedInto(f, groupname)
}

// LinksTo checks if the provided file is a symlink which links to path
// It returns error if the link can't be read
func LinksTo(f Filer, path string) error {
	dst, err := os.Readlink(path)
	if err != nil {
		return err
	}
	if dst == f.Path() {
		return nil
	}
	return fmt.Errorf("%s does not link to %s", f, path)
}

// IsSize checks if the provided file byte size is the same as the one passed in as paramter
// It returns error if os.Stat returns error
func IsSize(f Filer, size int64) error {
	fi, err := f.Info()
	if err != nil {
		return err
	}
	if fi.Size() == size {
		return nil
	}
	return fmt.Errorf("%s size is different from %d", f, size)
}

// ModTimeAfter checks if the file modification time of the file passed in as argument
// is older than the one passed in as paramter. It returns error if os.Stat returns error
// or if the file has been modified before the time supplied as argument
func ModTimeAfter(f Filer, mtime time.Time) error {
	fi, err := f.Info()
	if err != nil {
		return err
	}
	if fi.ModTime().After(mtime) {
		return nil
	}
	return fmt.Errorf("%s has been modifed before: %s", f, mtime)
}

// Contains checks if the provided file content can be matched by any of the RegExps
// passed in as paramters. It returns error if the provided file can't be open
func Contains(f Filer, contents ...*regexp.Regexp) error {
	return withFileReader(f, func(r io.Reader) error {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			for _, content := range contents {
				if content.Match(scanner.Bytes()) {
					return nil
				}
			}
		}
		return fmt.Errorf("%s does not match any provided regular expression", f)
	})
}

// MD5Equal checks if md5 checksum of the file is the same as the one passed in as paramter
// It returns error if the provided file can't be opened
func MD5Equal(f Filer, sum string) error {
	return withFileReader(f, func(r io.Reader) error {
		md5sum, err := utils.HashSum("md5", r)
		if err != nil {
			return err
		}
		if md5sum == sum {
			return nil
		}
		return fmt.Errorf("Expected md5 sum: %s Calculated sum: %s", sum, md5sum)
	})
}

// SHA256Equal checks if sha256 checksum of the file is the same as the one passed in as paramter
// It returns error if the provided file can't be opened
func SHA256Equal(f Filer, sum string) error {
	return withFileReader(f, func(r io.Reader) error {
		sha256sum, err := utils.HashSum("sha256", r)
		if err != nil {
			return err
		}
		if sha256sum == sum {
			return nil
		}
		return fmt.Errorf("Expected sha256 sum: %s Calculated sum: %s", sum, sha256sum)
	})
}
