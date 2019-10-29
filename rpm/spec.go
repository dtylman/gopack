package rpm

import (
	"fmt"
	"io"
	"strings"
)

//PkgName rpm contants
const (
	PkgName     = "Name"
	PkgVersion  = "Version"
	Release     = "Release"
	Summary     = "Summary"
	Group       = "Group"
	License     = "License"
	Vendor      = "Vendor"
	URL         = "URL"
	Packager    = "Packager"
	Requires    = "Requires"
	Conflicts   = "Conflicts"
	Prefix      = "Prefix"
	AutoReqProv = "AutoReqProv"
	BuildRoot   = "BuildRoot"
	BuildArch   = "BuildArch"
)

//SpecFile describes the spec file structure
type SpecFile struct {
	Header      map[string]string `json:"header"`
	Defines     []string          `json:"defines"`
	Requires    []string          `json:"requires"`
	Description string            `json:"description"`
	Prep        string            `json:"prep"`
	Build       string            `json:"build"`
	Pre         string            `json:"-"`
	Post        string            `json:"-"`
	PreUn       string            `json:"-"`
	PostUn      string            `json:"-"`
	ChangeLog   string            `json:"changelog"`
	Files       []string          `json:"files"`
}

func newSpec() *SpecFile {
	s := new(SpecFile)
	s.Header = make(map[string]string)
	s.Header[Release] = "1"
	s.Header[Group] = "default"
	s.Header[License] = "unknown"
	s.Header[Prefix] = "/"
	s.Header[AutoReqProv] = "no"
	return s
}

//Depends adds package dependencies
func (s *SpecFile) Depends(requires ...string) {
	s.Requires = requires
}

//AddDefine adds a define value
func (s *SpecFile) AddDefine(value string) {
	s.Defines = append(s.Defines, value)
}

//SetName sets packgae name
func (s *SpecFile) SetName(name string) {
	s.Header[PkgName] = name
}

//SetVersion sets package version
func (s *SpecFile) SetVersion(version string, revision string) {
	if revision != "" {
		s.Header[PkgVersion] = version + "_" + revision
	} else {
		s.Header[PkgVersion] = version
	}
}

//AddFile adds a file to the package
func (s *SpecFile) AddFile(name string) {
	if strings.Contains(name, " ") {
		s.Files = append(s.Files, fmt.Sprintf(`"%s"`, name))
	} else {
		s.Files = append(s.Files, name)
	}

}

//Write writes the rpm to writer
func (s *SpecFile) Write(writer io.Writer) error {
	var err error
	for _, define := range s.Defines {
		_, err = fmt.Fprintln(writer, "%define "+define)
		if err != nil {
			return err
		}
	}
	for name, value := range s.Header {
		_, err = fmt.Fprintln(writer, name+": "+value)
		if err != nil {
			return err
		}
	}
	for _, requires := range s.Requires {
		_, err = fmt.Fprintln(writer, Requires+": "+requires)
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprintln(writer, "%"+"description")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, s.Description)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(writer, "%"+"prep")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, s.Prep)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(writer, "%"+"build")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, s.Build)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(writer, "%"+"pre")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, s.Pre)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, "%"+"post")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, s.Post)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, "%"+"preun")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, s.PreUn)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, "%"+"postun")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, s.PostUn)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(writer, "%"+"files")
	if err != nil {
		return err
	}

	for _, fileName := range s.Files {
		_, err = fmt.Fprintln(writer, fileName)
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprintln(writer, "%"+"changelog")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, s.ChangeLog)
	if err != nil {
		return err
	}
	return nil
}

//PackageName returns a package file name
func (s *SpecFile) PackageName() string {
	return fmt.Sprintf("%s-%s-%s.%s.rpm", s.Header[PkgName], s.Header[PkgVersion], s.Header[Release], s.Header[BuildArch])
}
