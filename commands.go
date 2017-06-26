package main

import (
	"fmt"
	"os"
	"regexp"

	"gopkg.in/resty.v0"
)

func getDistroComponents(distro string) []string {
	var distroInfo DistroInfo
	resp, err := resty.R().
		SetResult(&distroInfo).
		Get("/distro/show/" + distro)
	if err != nil {
		fail("Cannot query cacus: %s", err)
	}
	if resp.StatusCode() != 200 {
		fail("Cannot query cacus: got %v", resp.Status())
	}
	if !distroInfo.Success {
		fail("Cannot query cacus: got %v", distroInfo.Message)
	}
	var availableComponents []string
	for _, v := range distroInfo.Result {
		if v.Distro == distro {
			availableComponents = v.Components
			break
		}
	}
	if availableComponents == nil {
		fail("Cannot find components for distro '%s'", distro)
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

	if _, ok := config.Instances[distro]; !ok {
		fail("Unknown distro '%s', please check parameters or config", distro)
	}

	resty.SetHostURL(fmt.Sprintf("%s/api/v1", config.Instances[distro].BaseURL))
	resty.SetHeader("Authorization", "Bearer "+config.Instances[distro].Token)

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
