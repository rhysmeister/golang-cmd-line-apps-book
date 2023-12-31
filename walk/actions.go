package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func itemInSlice(slice []string, item string) bool {
	for _, slice_item := range slice {
		if item == slice_item {
			return true
		}
	}
	return false
}

func filterOut(path string, ext []string, minSize int64, info os.FileInfo, days int) bool {
	if info.IsDir() || info.Size() < minSize {
		return true
	}
	if len(ext) > 0 && !itemInSlice(ext, filepath.Ext(path)) {
		return true
	}
	cutoffTime := time.Now().AddDate(0, 0, -days)
	if info.ModTime().Before(cutoffTime) { // suggestion S1008??
		return true
	}
	return false
}

func listFile(path string, out io.Writer) error {
	_, err := fmt.Fprintln(out, path)
	return err
}

func delFile(path string, delLogger *log.Logger) error {
	if err := os.Remove(path); err != nil {
		return err
	}
	delLogger.Println(path)
	return nil
}

func archiveFile(destDir, root, path string) error {
	info, err := os.Stat(destDir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", destDir)
	}
	relDir, err := filepath.Rel(root, filepath.Dir(path))
	if err != nil {
		return err
	}
	dest := fmt.Sprintf("%s.gz", filepath.Base(path))
	targetPath := filepath.Join(destDir, relDir, dest)
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	out, err := os.OpenFile(targetPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	in, err := os.Open(path)
	if err != nil {
		return err
	}
	defer in.Close()

	zw := gzip.NewWriter(out)

	zw.Name = filepath.Base(path)

	if _, err = io.Copy(zw, in); err != nil {
		return err
	}
	if err := zw.Close(); err != nil {
		return err
	}
	return out.Close()
}
