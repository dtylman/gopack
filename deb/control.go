package deb

type Control struct {
	Package      string
	Version      string
	Depends      string
	Architecture string
	Maintainer   string
	Conflicts    string
	Section      string
	Homepage     string
	Description  string
}
