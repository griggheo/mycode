package main

import (
	"fmt"
//	"os"
	"log"
        "bytes"
	"strings"
	"flag"
//	"sort"
	"os/exec"
	"io/ioutil"
)

func getCookbookVersionInEnvironment(cookbook *string, env string, verbose bool) string {
       cmd := exec.Command("knife", "environment", "show", env)
       var out bytes.Buffer
       cmd.Stdout = &out
       err := cmd.Run()
       if err != nil {
               log.Fatal(err)
       }
       //fmt.Printf("%s\n", out.String())
       lines := strings.Split(out.String(), "\n")
       var c, version string;
       for _, line := range lines {
               if ! strings.Contains(line, "=") {
                       continue
               }
               //fmt.Println(line)
               values := strings.Split(line, "=")

               /*
               for _, value := range values {
                       fmt.Printf("value: %s\n", value)
               }
               */
               c = strings.Trim(values[0], ": ")
               version = strings.Trim(values[1], " ")
	       if c == *cookbook {
       			if verbose {
				fmt.Printf("Cookbook \"%s\" has version \"%s\" in env \"%s\"\n", *cookbook, version, env)
       			}
			return version
		}
       }
       return ""
}

func getLatestCookbookVersionOnChefServer(cookbook *string, verbose bool) string {
       cmd := exec.Command("knife", "cookbook", "show", *cookbook)
       var out bytes.Buffer
       cmd.Stdout = &out
       err := cmd.Run()
       if err != nil {
               log.Fatal(err)
       }
       fields := strings.Fields(out.String())
       var version string;
       version = fields[1]
       if verbose {
		fmt.Printf("Cookbook \"%s\" has version \"%s\" as latest on Chef server\n", *cookbook, version)
       }
       return version
}

func getCookbookVersionInLocalMetadata(cookbook *string, verbose bool) string {
        filename := fmt.Sprintf("/Users/grig.gheorghiu/chef-repo/cookbooks/%s/metadata.rb", *cookbook)

        content, err := ioutil.ReadFile(filename)
        if err != nil {
                log.Fatalln("Error reading file", filename)
        }

        // content returned as []byte not string
	lines := strings.Split(string(content), "\n")
        var version string
	var fields []string
        for _, line := range lines {
                if ! strings.Contains(line, "version") {
                        continue
                }
		fields = strings.Fields(line)
		version = fields[1]
		version = strings.Replace(version, "\"", "", -1)
		break
        }
        if verbose {
		fmt.Printf("Cookbook \"%s\" has version \"%s\" in local metadata.rb\n", *cookbook, version)
        }
        return version
}

func getGitDiff(cookbook *string) ([]byte, error) {
        filename := fmt.Sprintf("/Users/grig.gheorghiu/chef-repo/cookbooks/%s/metadata.rb", *cookbook)
	fmt.Println("\nGetting git diff for ", filename)
	cmd := exec.Command("git", "diff", "--exit-code", filename)
	combined_output, err := cmd.CombinedOutput()
	return combined_output, err
}

func main() {
	cookbook := flag.String("cookbook", "", "the cookbook we are inspecting")
	flag.Parse()
	envs := []string {"prod", "stg", "qa1"}
	for _, env := range envs {
		getCookbookVersionInEnvironment(cookbook, env, true)
		// version := getCookbookVersionInEnvironment(cookbook, env, true)
		// fmt.Println("version:", version)
	}
	getLatestCookbookVersionOnChefServer(cookbook, true)
	getCookbookVersionInLocalMetadata(cookbook, true)
	
	output, err := getGitDiff(cookbook)
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
	} else {
		fmt.Println("Found no difference between GitHub and local version.rb")
	}
}	
