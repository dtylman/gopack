package rpm

import (
	"fmt"
	"io"
	"strings"
)

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

type SpecFile struct {
	Header      map[string]string
	Defines     []string
	Requires    []string
	Description string
	Prep        string
	Build       string
	Pre         string
	Post        string
	PreUn       string
	PostUn      string
	ChangeLog   string
	Files       []string
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

func (s *SpecFile) Depends(requires ...string) {
	s.Requires = requires
}
func (s *SpecFile) AddDefine(value string) {
	s.Defines = append(s.Defines, value)
}

func (s *SpecFile) SetName(name string) {
	s.Header[PkgName] = name
}

func (s *SpecFile) SetVersion(version string, revision string) {
	if revision != "" {
		s.Header[PkgVersion] = version + "_" + revision
	} else {
		s.Header[PkgVersion] = version
	}
}

func (s *SpecFile) AddFile(name string) {
	if strings.Contains(name, " ") {
		s.Files = append(s.Files, fmt.Sprintf(`"%s"`, name))
	} else {
		s.Files = append(s.Files, name)
	}

}

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
	_, err = fmt.Fprintln(writer, "%description")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, s.Description)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(writer, "%prep")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, s.Prep)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(writer, "%build")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, s.Build)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(writer, "%pre")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, s.Pre)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, "%post")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, s.Post)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, "%preun")
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
	_, err = fmt.Fprintln(writer, "%postun")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, s.PostUn)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(writer, "%files")
	if err != nil {
		return err
	}

	for _, fileName := range s.Files {
		_, err = fmt.Fprintln(writer, fileName)
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprintln(writer, "%changelog")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, s.ChangeLog)
	if err != nil {
		return err
	}
	return nil
}

func (s *SpecFile) PackageName() string {
	return fmt.Sprintf("%s-%s-%s.%s.rpm", s.Header[PkgName], s.Header[PkgVersion], s.Header[Release], s.Header[BuildArch])
}
