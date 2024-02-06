package main

// set GOOS=linux
// set GOOS=windows
// set GOARCH=amd64 go build
// ./parseconf reposyncDirr insert_position_bblayer.conf insert_position_local.conf
// ./sep reposyncDir/tmp_dir/deploy/imgDir/spdx/
import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
)

func main() {
	fmt.Println(strings.Repeat("=", 10), "Start of file", path.Base(os.Args[0]), strings.Repeat("=", 10))
	var wg sync.WaitGroup
	wg.Add(2)

	repoAddText := "  ${TOPDIR}/../meta-doubleopen \\"
	inheritText := "INHERIT += \"doubleopen\""

	workDir := "workdir"
	confDir := workDir + string(os.PathSeparator) + "conf"

	bblConfFile := confDir + string(os.PathSeparator) + "bblayers.conf"
	localConfFile := confDir + string(os.PathSeparator) + "local.conf"

	var folderPath string
	var posToInsertBBLayerConfFile, posToInsertLocalConfFile int
	if len(os.Args) == 4 {
		folderPath = os.Args[1]
		posToInsertBBLayerConfFile, _ = strconv.Atoi(os.Args[2])
		posToInsertLocalConfFile, _ = strconv.Atoi(os.Args[3])
	} else {
		fmt.Println("Please provide reposync folder to process that contains workdir , insert position to bblayerconf file and local conf file as an argument")
		fmt.Println("check https://doubleopen-project.github.io/tldr.html")
		os.Exit(1)
	}
	// fmt.Println(folderPath)

	// go func(folderPath string) {
	// 	defer wg.Done()
	// 	doubleOpenRepo := folderPath + string(os.PathSeparator) + "meta-doubleopen"
	// 	doubleOpenURL := "git@gitext.elektrobitautomotive.com:oss/meta-doubleopen.git"
	// 	if _, err := os.Stat(doubleOpenRepo); errors.Is(err, fs.ErrNotExist) {

	// 		fmt.Println(doubleOpenRepo, " git Repo doesnt exist in", folderPath, ". Cloning from git@gitext.elektrobitautomotive.com:oss/meta-doubleopen.git")
	// 		_, err := git.PlainClone(doubleOpenRepo, false, &git.CloneOptions{
	// 			URL:      doubleOpenURL,
	// 			Progress: os.Stdout,
	// 		})
	// 		if err != nil {
	// 			fmt.Println(" Error while cloning the meta-doubleopen repo")
	// 			fmt.Println(err)
	// 			os.Exit(1)
	// 		}

	// 	} else {
	// 		fmt.Println("Repo already exists so Skipping clone!")
	// 	}
	// }(folderPath)

	go func(folderPath string) {
		defer wg.Done()
		bblConfFile = folderPath + string(os.PathSeparator) + bblConfFile
		if _, err := os.Stat(bblConfFile); errors.Is(err, fs.ErrNotExist) {
			fmt.Println(bblConfFile, "doesnt exist")
			os.Exit(1)
		}
		if !IsExist("meta-doubleopen", bblConfFile) {
			InsertStringToFile(bblConfFile, repoAddText+"\n", posToInsertBBLayerConfFile)
		} else {
			fmt.Println("Skipping modifying", bblConfFile)
		}
	}(folderPath)

	go func(folderPath string) {
		defer wg.Done()
		localConfFile = folderPath + string(os.PathSeparator) + localConfFile
		if _, err := os.Stat(localConfFile); errors.Is(err, fs.ErrNotExist) {
			fmt.Println(localConfFile, "doesnt exist")
			os.Exit(1)
		}
		if !IsExist("doubleopen", localConfFile) {
			InsertStringToFile(localConfFile, inheritText+"\n", posToInsertLocalConfFile)
		} else {
			fmt.Println("Skipping modifying", localConfFile)
		}
	}(folderPath)

	wg.Wait()
	fmt.Println(strings.Repeat("=", 10), "End of file", path.Base(os.Args[0]), strings.Repeat("=", 10))
}
