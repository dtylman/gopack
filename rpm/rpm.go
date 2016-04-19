package rpm

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"io"
	"os/exec"
	"bytes"
	"fmt"
	"strings"
)

type Rpm struct {
	Spec* SpecFile
	name string
	workingFolder string
	buildRoot string
}

func New(name, version, revision, arch string) (*Rpm,error){
	r :=new(Rpm)
	r.Spec = newSpec()
	r.name = name
	var err error
	r.workingFolder, err = ioutil.TempDir("",name)
	if err != nil {
		return nil, err
	}
	r.Spec.SetName(r.name)
	r.Spec.SetVersion(version,revision)
	r.buildRoot=filepath.Join(r.workingFolder,"BUILD")
	err=os.MkdirAll(r.buildRoot,0755)
	if err != nil {
		os.RemoveAll(r.workingFolder)
		return nil,err
	}
	r.Spec.Header["BuildRoot"]=r.buildRoot
	r.Spec.AddDefine("_tmppath " + os.TempDir())
	r.Spec.AddDefine("_topdir " + r.workingFolder)
	r.Spec.AddDefine("_sourcedir " + r.workingFolder)
	r.Spec.AddDefine("buildroot " + r.buildRoot )
	// Disable the stupid stuff rpm distros include in the build process by default:
	// Disable any prep shell actions. replace them with simply 'true'
	r.Spec.AddDefine("__spec_prep_post true")
	r.Spec.AddDefine("__spec_prep_pre true")
	r.Spec.AddDefine("__spec_build_post true")
	r.Spec.AddDefine("__spec_build_pre true")
	r.Spec.AddDefine("__spec_install_post true")
	r.Spec.AddDefine("__spec_install_pre true")
	r.Spec.AddDefine("__spec_clean_post true")
	r.Spec.AddDefine("__spec_clean_pre true")

	return r,nil
}

func (r* Rpm) Close() error{
	return os.RemoveAll(r.workingFolder)
}
/*

	//Running rpmbuild {:args=>["rpmbuild", "-bb", "--define", "buildroot /tmp/package-rpm-build20160419-10869-1byo5sr/BUILD", "--define", "_topdir /tmp/package-rpm-build20160419-10869-1byo5sr",
	// "--define", "_sourcedir /tmp/package-rpm-build20160419-10869-1byo5sr", "--define", "_rpmdir /tmp/package-rpm-build20160419-10869-1byo5sr/RPMS", "--define", "_tmppath /tmp",
	// "/tmp/package-rpm-build20160419-10869-1byo5sr/SPECS/demistoserver.spec"], :level=>:info}

 */
func (r* Rpm) Create() error{
	rpms:=filepath.Join(r.workingFolder,"RPMS")
	err:=os.MkdirAll(rpms,0755)
	if err != nil {
		return err
	}
	specFolder:=filepath.Join(r.workingFolder,"SPECS")
	err=os.MkdirAll(specFolder,0755)
	if err != nil {
		return err
	}
	specFile, err:=os.Create(filepath.Join(specFolder,r.name+".spec"))
	if err != nil {
		return err
	}
	defer specFile.Close()
	err= r.Spec.Write(specFile)
	if err != nil {
		return err
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd:=exec.Command("rpmbuild","-bb", specFile.Name())
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err=cmd.Run()
	if err != nil {
		return fmt.Errorf("rpmbuild failed with: %v. Stdout: %v. Stderr: %v",err,stdout.String(),stderr.String())
	}
	return filepath.Walk(rpms, r.movePackageFile)
}

func (r* Rpm) movePackageFile(path string, info os.FileInfo, err error) error{
	if (!info.IsDir()){
		if strings.HasSuffix(strings.ToLower(path),".rpm"){
			return os.Rename(path,filepath.Base(path))
		}
	}
	return nil
}

func (r* Rpm) AddEmptyFolder(name string) error{
	return nil
}

func (r* Rpm) AddFolder(path string, prefix string) error {
	return nil
}

func (r *Rpm) AddFile(path string, prefix string) error {
	srcFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	fileName:=filepath.Base(path)
	destFolder:=filepath.Join(r.buildRoot,prefix)
	err=os.MkdirAll(destFolder,0755)
	if err != nil {
		return err
	}
	destFile, err := os.Create(filepath.Join(destFolder,fileName))
	if err != nil {
		return err
	}
	defer destFile.Close()
	fmt.Printf("%v->%v\n",srcFile.Name(),destFile.Name())
	_, err = io.Copy(destFile,srcFile)
	if err != nil {
		return err
	}
	r.Spec.AddFile(filepath.Join(prefix,fileName))
	return nil
}
