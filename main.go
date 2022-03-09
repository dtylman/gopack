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

//Options holds commandline options
var Options struct {
	//OutPath is the output path
	OutPath string
	//ConfigFileName config file name
	ConfigFileName string
	//BuildRPM if true, will build RPM
	BuildRPM bool
	//BuildDeb if true, will build Deb
	BuildDeb bool
	//Version to specify in command line
	Version string
	//Revision to specify in command line
	Revision string
}

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
	fileName, err := pkg.Create(Options.OutPath)
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

	fileName, err := deb.Create(Options.OutPath)
	if err != nil {
		return fmt.Errorf("failed to create package: '%v'", err)
	}

	log.Printf("created: '%v'", fileName)
	return nil
}

func create() error {
	cfg, err := config.Load(Options.ConfigFileName)
	if err != nil {
		return fmt.Errorf("failed to load config file: '%v', error: %v", Options.ConfigFileName, err)
	}
	if Options.Version != "" {
		log.Printf("setting version to %v", Options.Version)
		cfg.Version = Options.Version
	}
	if Options.Revision != "" {
		log.Printf("setting revision to %v", Options.Revision)
	}
	if Options.BuildRPM {
		err = createRPM(cfg)
		if err != nil {
			return fmt.Errorf("failed to create rpm: %v", err)
		}
	}
	if Options.BuildDeb {
		err = createDeb(cfg)
		if err != nil {
			return fmt.Errorf("failed to create deb: %v", err)
		}
	}
	if !Options.BuildDeb && !Options.BuildRPM {
		return errors.New("must specify ether 'rpm' or 'deb'")
	}
	return nil
}

func main() {
	flag.BoolVar(&Options.BuildRPM, "rpm", false, "build rpm package")
	flag.BoolVar(&Options.BuildDeb, "deb", false, "build deb package")
	flag.StringVar(&Options.ConfigFileName, "conf", "pkg.config.json", "config file name")
	flag.StringVar(&Options.OutPath, "output", "", "output path")
	flag.StringVar(&Options.Version, "version", "", "specify package version")
	flag.StringVar(&Options.Revision, "revision", "", "specify package revision")
	flag.Parse()
	err := create()
	if err != nil {
		log.Println(err)
	}
}
