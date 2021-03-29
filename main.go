package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/dtylman/gopack/config"
	"github.com/dtylman/gopack/rpm"

	"github.com/dtylman/gopack/deb"
)

//output folder
var outputPath string

func addScript(sourceFileName string, value *string) error {
	if sourceFileName == "" {
		return nil
	}
	data, err := ioutil.ReadFile(sourceFileName)
	if err != nil {
		return err
	}
	if value == nil {
		return errors.New("script value is null")
	}
	log.Printf("Adding script '%v' (%v bytes)", sourceFileName, len(data))
	*value = string(data)
	return nil
}

func createRPM(cfg *config.PackageOptions) error {
	log.Println("Creating rpm...")
	pkg, err := rpm.New(cfg.Name, cfg.Version, cfg.Revision, cfg.Arch)
	if err != nil {
		return err
	}
	pkg.Spec.Header[rpm.Summary] = cfg.Name
	pkg.Spec.Header[rpm.Packager] = cfg.Maintainer
	pkg.Spec.Header[rpm.URL] = cfg.Homepage
	pkg.Spec.Depends(strings.Split(cfg.Depends, " ")...)
	pkg.Spec.Description = cfg.Description
	for path, prefix := range cfg.Folders {
		if prefix == "" {
			err = pkg.AddEmptyFolder(path)
		} else {
			err = pkg.AddFolder(path, prefix)
		}
		if err != nil {
			return fmt.Errorf("failed to add folder: %v", err)
		}
	}
	for source, target := range cfg.Files {
		err = pkg.AddFile(source, target)
		if err != nil {
			return fmt.Errorf("failed to ad file: '%v'", err)
		}
	}
	fileName, err := pkg.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create package: '%v'", err)
	}
	log.Printf("created: '%v'", fileName)
	return nil
}

func createDeb(cfg *config.PackageOptions) error {
	log.Println("Creating deb...")
	deb, err := deb.New(cfg.Name, cfg.Version, cfg.Revision, cfg.Arch)
	if err != nil {
		return err
	}
	deb.Info.Description = cfg.Description
	deb.Info.Homepage = cfg.Homepage
	deb.Info.Depends = cfg.Depends
	deb.Info.Section = cfg.Section
	deb.Info.Maintainer = cfg.Maintainer

	for path, prefix := range cfg.Folders {
		if prefix == "" {
			log.Printf("Adding empty folder '%v'", path)
			err = deb.AddEmptyFolder(path)
		} else {
			log.Printf("Adding folder '%v'->'%v'", path, prefix)
			err = deb.AddFolder(path, prefix)
		}
		if err != nil {
			return fmt.Errorf("failed to add folder: %v", err)
		}
	}
	for source, target := range cfg.Files {
		log.Printf("Adding file '%v'->'%v'", source, target)
		err = deb.AddFile(source, target)
		if err != nil {
			return fmt.Errorf("failed to add file: '%v'", err)
		}
	}

	err = addScript(cfg.Script.PostInst, &deb.PostInst)
	if err != nil {
		return fmt.Errorf("failed to add script '%v'", err)
	}
	err = addScript(cfg.Script.PreInst, &deb.PreInst)
	if err != nil {
		return fmt.Errorf("failed to add script '%v'", err)
	}

	err = addScript(cfg.Script.PostUnInst, &deb.PostRm)
	if err != nil {
		return fmt.Errorf("failed to add script '%v'", err)
	}

	err = addScript(cfg.Script.PreUnInst, &deb.PreRm)
	if err != nil {
		return fmt.Errorf("failed to add script '%v'", err)
	}

	fileName, err := deb.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create package: '%v'", err)
	}

	log.Printf("created: '%v'", fileName)
	return nil
}

func create(configFile string, rpm, deb bool) error {
	cfg, err := config.Load(configFile)
	if err != nil {
		return fmt.Errorf("failed to load config file: '%v', error: %v", configFile, err)
	}
	if rpm {
		err = createRPM(cfg)
		if err != nil {
			return fmt.Errorf("failed to create rpm: %v", err)
		}
	}
	if deb {
		err = createDeb(cfg)
		if err != nil {
			return fmt.Errorf("failed to create deb: %v", err)
		}
	}
	return nil
}

func main() {
	rpm := flag.Bool("rpm", false, "build rpm pacakge")
	deb := flag.Bool("deb", false, "build deb package")
	conf := flag.String("conf", "pkg.config.json", "config file name")
	flag.StringVar(&outputPath, "output", "", "output path")
	flag.Parse()
	err := create(*conf, *rpm, *deb)
	if err != nil {
		log.Println(err)
	}
}
