package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"sort"
	"strings"
)

type Package struct {
	Name    string
	Version string
	Arch    string
}

type SortedPackages []*Package

func (l SortedPackages) Len() int           { return len(l) }
func (l SortedPackages) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l SortedPackages) Less(i, j int) bool { return pkgless(l[i], l[j]) }

func pkgless(x, y *Package) bool {
	if x.Name == y.Name {
		if x.Version == y.Version {
			return x.Arch < y.Arch
		} else {
			return x.Version < y.Version
		}
	} else {
		return x.Name < y.Name
	}
}

func getCmdOutputScanner(name string, args ...string) *bufio.Scanner {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		log.Fatalln(err.Error())
	}
	reader := bytes.NewReader(output)
	return bufio.NewScanner(reader)
}

func getPkgsInstalled() []*Package {
	var pkgs []*Package
	scanner := getCmdOutputScanner("dpkg-query", "-f",
		"${Package} ${Version} ${Architecture}\\n", "-W")
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		pkg := &Package{Name: fields[0], Version: fields[1], Arch: fields[2]}
		pkgs = append(pkgs, pkg)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err.Error())
	}
	return pkgs
}

func getPkgsAvailable() []*Package {
	scanner := getCmdOutputScanner("apt-cache", "dumpavail")
	var pkgs []*Package
	var pkg *Package
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" && pkg != nil {
			pkgs = append(pkgs, pkg)
			pkg = nil
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		switch fields[0] {
		case "Package:":
			pkg = &Package{}
			pkg.Name = fields[1]
		case "Version:":
			pkg.Version = fields[1]
		case "Architecture:":
			pkg.Arch = fields[1]
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err.Error())
	}
	if pkg != nil {
		pkgs = append(pkgs, pkg)
	}
	return pkgs
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	pkgsAvail := getPkgsAvailable()
	sort.Sort(SortedPackages(pkgsAvail))
	for _, pkg := range getPkgsInstalled() {
		idx := sort.Search(len(pkgsAvail), func(i int) bool {
			return pkgless(pkg, pkgsAvail[i])
		})
		if idx >= len(pkgsAvail) || *pkgsAvail[idx-1] != *pkg {
			fmt.Printf("%s:%s %s absent\n", pkg.Name, pkg.Arch, pkg.Version)
		} else {
			fmt.Printf("%s:%s %s present\n", pkg.Name, pkg.Arch, pkg.Version)
		}
	}
}
