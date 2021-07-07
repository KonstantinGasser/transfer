package diff

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/mod/sumdb/dirhash"
)

var (
	ErrDifferntHash = fmt.Errorf("src and dst hash mismatch")
	ErrDataMismatch = fmt.Errorf("compared data differs")
)

type Compareable struct {
	root string
	dst  string
}

// Compare sets the root directory any comparision will be performed
// against with
func Compare(root, dst string) Compareable {
	return Compareable{
		root: root,
		dst:  dst,
	}
}

// Updatable allows to compare any directory with the Comparables root directory
// to see if the indeed are the same - returning the updatable files
func (c Compareable) Updatable() ([]string, error) {

	var updatable []string
	if err := filepath.Walk(c.root, c.traverse(&updatable)); err != nil {
		return nil, err
	}
	return updatable, nil
}

func (c Compareable) traverse(updatable *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == c.root {
			return nil
		}
		if info.IsDir() {
			if err := c.hashDir(path, info.Name()); err != nil {
				if err == ErrDataMismatch {
					fmt.Println("update for: ", path[len(c.root):])
					*updatable = append(*updatable, path[len(c.root):])
				}

			}
		}
		if info.Mode().IsRegular() {
			fmt.Println("regular: ", path)
			dir, _ := filepath.Split(path)
			files, err := dirhash.DirFiles(dir, "")
			if err != nil {
				return err
			}
			fmt.Println("files: ", files)
		}
		return nil
	}
}

func (c Compareable) hashDir(path, name string) error {
	subdir := path[len(c.root):] // cut-off root dir keep only sub-dirs
	fmt.Println("sub-dir: ", subdir)
	srcFiles, err := dirhash.DirFiles(filepath.Join(c.root, subdir), name)
	if err != nil {
		return err
	}
	dstFiles, err := dirhash.DirFiles(filepath.Join(c.dst, subdir), name)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrDataMismatch
		}
		return err
	}

	srcHash, err := hash(srcFiles, c.root, subdir)
	if err != nil {
		return err
	}
	dstHash, err := hash(dstFiles, c.dst, subdir)
	if err != nil {
		return err
	}
	if srcHash != dstHash {
		return ErrDataMismatch
	}
	return nil
}

func hash(files []string, dir, subdir string) (string, error) {
	return dirhash.Hash1(files, func(name string) (io.ReadCloser, error) {
		return os.Open(filepath.Join(dir, subdir, name))
	})
}

// func (c Compareable) hashFiles(srcFiles, dstFiles []string, subdir string) error {
// 	srcHash, err := hash(srcFiles, c.root, subdir)
// 	if err != nil {
// 		return err
// 	}
// 	dstHash, err := hash(dstFiles, c.dst, subdir)
// 	if err != nil {
// 		return err
// 	}
// 	if srcHash != dstHash {
// 		return ErrDataMismatch
// 	}
// 	return nil
// }

// func (c Compareable) hashDir(dir, folder string) (string, error) {
// 	files, err := dirhash.DirFiles(filepath.Join(dir, folder), "")
// 	if err != nil {
// 		return "", err
// 	}
// 	return dirhash.Hash1(files, func(name string) (io.ReadCloser, error) {
// 		return os.Open(filepath.Join(dir, folder, name))
// 	})
// }

// func TMPCompare(src string) Compareable {
// 	return Compareable{src: src}
// }

// func (c Compareable) TMPWith(dst string) ([]string, error) {

// 	var updatable []string
// 	err := filepath.Walk(c.src, func(path string, info os.FileInfo, err error) error {
// 		if path == c.src {
// 			return nil
// 		}
// 		if err != nil {
// 			return err
// 		}
// 		if info.IsDir() {
// 			srcHash, err := c.hashDir(c.src, info.Name())
// 			if err != nil {
// 				return err
// 			}
// 			dstHash, err := c.hashDir(dst, info.Name())
// 			if err != nil {
// 				if os.IsNotExist(err) {
// 					fmt.Println("===============================================================")
// 					fmt.Printf("Folder %s does not exists - coping full folder\n", path)
// 					fmt.Println("===============================================================")
// 					updatable = append(updatable, path)
// 					return nil
// 				}
// 				return err
// 			}

// 			if srcHash != dstHash {
// 				fmt.Println("===============================================================")
// 				fmt.Println("For: ", info.Name())
// 				fmt.Println("want hash: ", srcHash)
// 				fmt.Println("have hash: ", dstHash)
// 				fmt.Println("===============================================================")
// 				updatable = append(updatable, path)
// 			}
// 			return nil
// 		}

// 		if set := strings.Split(path, "/"); len(set) == 2 && info.Mode().IsRegular() {
// 			fmt.Println("===============================================================")
// 			fmt.Printf("File %s check - coping full folder\n", path)
// 			fmt.Println("===============================================================")
// 			srcFile, err := dirhash.Hash1([]string{info.Name()}, func(name string) (io.ReadCloser, error) {
// 				return os.Open(filepath.Join(c.src, name))
// 			})
// 			if err != nil {
// 				return err
// 			}
// 			dstFile, err := dirhash.Hash1([]string{info.Name()}, func(name string) (io.ReadCloser, error) {
// 				return os.Open(filepath.Join(dst, name))
// 			})
// 			if err != nil {
// 				if os.IsNotExist(err) {
// 					fmt.Println("===============================================================")
// 					fmt.Printf("File %s does not exists - coping full folder\n", path)
// 					fmt.Println("===============================================================")
// 					updatable = append(updatable, path)
// 					return nil
// 				}
// 				return err
// 			}
// 			if srcFile != dstFile {
// 				fmt.Println("===============================================================")
// 				fmt.Println("For: ", info.Name())
// 				fmt.Println("want hash: ", srcFile)
// 				fmt.Println("have hash: ", dstFile)
// 				fmt.Println("===============================================================")
// 				updatable = append(updatable, path)
// 			}
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return updatable, nil
// }
