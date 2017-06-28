package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/resty.v0"
)

func getDistro(distro string) (*DistroShowResult, error) {
	var result DistroShowResult
	resp, err := resty.R().
		SetResult(&result).
		Get("/distro/show/" + distro)
	if err != nil {
		return nil, cacaerr("Cannot query cacus: %s", err)
	}
	if resp.StatusCode() != 200 {
		return nil, cacaerr("Cannot query cacus: got %v", resp.Status())
	}
	if !result.Success {
		return nil, cacaerr("Cannot query cacus: got %v", result.Message)
	}

	return &result, nil
}

func getDistroComponents(distro string) []string {
	var availableComponents []string
	distroInfo, err := getDistro(distro)
	if err != nil {
		fail("Cannot find components for distro %v", err)
	}
	for _, v := range distroInfo.Result {
		if v.Distro == distro {
			availableComponents = v.Components
			break
		}
	}

	return availableComponents
}

func uploadPackage(args []string) {
	if len(args) < 2 {
		fail("USAGE: upload PACKAGE [...]  DISTRO/COMPONENT")
	}
	var packages = args[:len(args)-1]
	var dest = args[len(args)-1]

	// i think that regexp pkg was written  by some brain-damaged psycho
	r := regexp.MustCompile(`^([-_.A-Za-z0-9]+)/([-_a-z0-9]+)$`)
	m := r.FindStringSubmatch(dest)
	if m == nil {
		fail("USAGE: upload PACKAGE [...]  DISTRO/COMPONENT")
	}
	distro := m[1]
	component := m[2]

	/*
		if _, ok := config.Instances[distro]; !ok {
			fail("Unknown distro '%s', please check parameters or config", distro)
		}

		resty.SetHostURL(fmt.Sprintf("%s/api/v1", config.Instances[distro].BaseURL))
		resty.SetHeader("Authorization", "Bearer "+config.Instances[distro].Token)
	*/

	availableComponents := getDistroComponents(distro)
	for _, c := range availableComponents {
		if component == c {
			// i fucking REALLY doubt why i'm using this FUCKED language
			goto ok
		}
	}
	fail("Unknown component '%s'. Available components: %v", component, availableComponents)
ok:

	//finally, we can upload the packages
	for _, file := range packages {
		fmt.Printf("Uploading %s... ", file)
		f, err := os.OpenFile(file, os.O_RDONLY, 0644)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			continue
		}
		var status CacusStatus
		resp, err := resty.R().
			SetBody(f).
			SetResult(&status).
			Put(fmt.Sprintf("package/upload/%s/%s", distro, component))
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			continue
		}
		if resp.StatusCode() != 201 {
			fmt.Printf("ERROR: got %v", resp.Status())
			continue
		}
		if !status.Success {
			fmt.Printf("ERROR: %s\n", status.Message)
			continue
		}
		fmt.Printf("SUCCESS: %s\n", status.Message)
	}
}

func showDistro(args []string) {
	for _, distro := range args {
		fmt.Printf("Distro '%s': ", distro)
		info, err := getDistro(distro)
		if err != nil {
			fmt.Printf("ERROR: %v", err)
			continue
		}
		if len(info.Result) < 1 {
			fmt.Printf("Error: not found\n")
			continue
		}
		dumpDistro(info.Result[0])
	}
}

func dumpDistro(distro DistroInfo) {
	if len(distro.Origin) < 1 {
		distro.Origin = "N/A"
	}
	fmt.Print("\n")
	fmt.Printf("\tName: %v\n", distro.Distro)
	fmt.Printf("\tDescription: %v\n", distro.Description)
	fmt.Printf("\tComponents: %v\n", distro.Components)
	fmt.Printf("\tNumber of packages: %v\n", distro.Packages)
	fmt.Printf("\tType: %v\n", distro.Type)
	fmt.Printf("\tOrigin: %v\n", distro.Origin)
	fmt.Printf("\tLast updated at: %v\n", distro.Lastupdated)
	fmt.Print("\n\n")
}

func searchPackages(args []string) {
	terms := flag.NewFlagSet("Search Options", flag.ExitOnError)
	distro := terms.String("distro", "", "Distro name")
	pkg := terms.String("pkg", "", "Package regex")
	ver := terms.String("ver", "", "Package version")
	comp := terms.String("comp", "", "Package component")
	descr := terms.String("descr", "", "Package description")
	terms.Parse(args)

	var url string
	if len(*distro) > 0 {
		url = fmt.Sprintf("/package/search/%s", *distro)
	} else {
		url = "/package/search"
	}

	var result PkgSearchResult
	resp, err := resty.R().
		SetBody(PkgSearchParams{*pkg, *ver, *comp, *descr}).
		SetResult(&result).
		Post(url)

	if err != nil {
		fail("Search failed: %v", err)
	}
	if resp.StatusCode() != 200 {
		fail("Search failed: got %v", resp.Status())
	}
	for distro, entries := range result.Result {
		fmt.Printf("\033[32m==== Results for distro \033[33m%s\033[32m ====\033[0m\n ", distro)
		dumpSearchResult(entries)
		fmt.Print("\n\n")
	}
}

func dumpSearchResult(entries []PkgSearchResultEntry) {
	if len(entries) < 1 {
		fmt.Print("\nNothing found\n")
	} else {
		para := func(f, s interface{}) {
			fmt.Printf("\t\033[1m%v:\033[0m %v\n", f, s)
		}
		for _, entry := range entries {
			para("Package", entry.Package)
			para("Version", entry.Version)
			para("Maintainer", entry.Maintainer)
			para("Architecure", entry.Architecture)
			para("Components", entry.Components)
			para("Description", strings.Replace(entry.Description, "\n", "\n\t\t", -1))
			fmt.Print("\n")
		}
	}
}
