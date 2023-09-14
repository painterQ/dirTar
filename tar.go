package tar

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// DeTarDir tar -x
func DeTarDir(srcPath, targetDir string) error {
	f, err := os.OpenFile(srcPath, os.O_RDONLY, 0666)
	if err != nil {

	}
	tarReader := tar.NewReader(f)
	var h *tar.Header
	var tErr error
	for h, tErr = tarReader.Next(); tErr == nil; h, tErr = tarReader.Next() {
		fmt.Println(h.Name)
		//todo
	}

	if tErr != io.EOF {
		return fmt.Errorf("DeTarDir error: %v", tErr)
	}
	return nil
}

// TarDir tar
func TarDir(srcDir, targetPath string) error {
	tarfile, err := os.Create(targetPath)
	if err != nil {

	}
	defer func() {
		_ = tarfile.Close()
	}()

	tarwriter := tar.NewWriter(tarfile)

	sfileInfo, err := os.Stat(srcDir)
	if err != nil {
		return err
	}
	if !sfileInfo.IsDir() {
		return fmt.Errorf("expect a dir")
	}
	return tarFolder(srcDir, tarwriter)
}

func tarFolder(directory string, tarwriter *tar.Writer) error {
	var baseFolder = filepath.Base(directory)
	return filepath.Walk(directory, func(targetpath string, file os.FileInfo, err error) error {
		switch {
		case file == nil:
			//read the file failure
			return err
		case file.IsDir():
			// information of file or dir
			header, errFIH := tar.FileInfoHeader(file, "")
			if errFIH != nil {
				return errFIH
			}
			header.Name = filepath.Join(baseFolder, strings.TrimPrefix(targetpath, directory))
			if err = tarwriter.WriteHeader(header); err != nil {
				return err
			}
			return nil
		default:
			if !file.Mode().IsRegular() {
				return fmt.Errorf("unregular file:%v", file.Name())
			}
			//baseFolder is the tar file path
			var fileFolder = filepath.Join(baseFolder, strings.TrimPrefix(targetpath, directory))
			return tarFile(fileFolder, targetpath, file, tarwriter)

		}
	})
}

func tarFile(directory string, filesource string, sfileInfo os.FileInfo, tarwriter *tar.Writer) error {
	sfile, err := os.Open(filesource)
	if err != nil {
		return err
	}
	defer func() {
		_ = sfile.Close()
	}()

	header, err := tar.FileInfoHeader(sfileInfo, "")
	if err != nil {
		fmt.Println(err)
		return err
	}
	header.Name = directory
	if err = tarwriter.WriteHeader(header); err != nil {
		return err
	}

	buf := make([]byte, 512)
	if _, err = io.CopyBuffer(tarwriter, sfile, buf); err != nil {
		return err
	}
	return nil
}
