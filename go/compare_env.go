package main

import (
	"fmt"
//	"os"
	"log"
        "bytes"
	"strings"
	"flag"
	"sort"
	"os/exec"
	"io/ioutil"
)

func getChefEnvironment(environment *string, verbose bool) map[string]string {
       fmt.Println("Getting Chef cookbook versions for environment:", *environment)
       cmd := exec.Command("knife", "environment", "show", *environment)
       var out bytes.Buffer
       cmd.Stdout = &out
       err := cmd.Run()
       if err != nil {
               log.Fatal(err)
       }
       //fmt.Printf("%s\n", out.String())
       lines := strings.Split(out.String(), "\n")
       m := make(map[string]string)
       keys := make([]string, 0, len(lines))
       var cookbook, version string;
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
               cookbook = strings.Trim(values[0], ": ")
               version = strings.Trim(values[1], " ")
               m[cookbook] = version
               keys = append(keys, cookbook)
       }
       if verbose {
               sort.Strings(keys)
               for _, key := range keys {
                       fmt.Printf("cookbook \"%s\" has version \"%s\"\n", key, m[key])
               }
       }
       return m
}

func getLocalEnvironment(environment *string, verbose bool) map[string]string {
	fmt.Println("Getting local cookbook versions for environment:", *environment)
        filename := fmt.Sprintf("/Users/grig.gheorghiu/chef-repo/environments/%s.rb", *environment)

        content, err := ioutil.ReadFile(filename)
        if err != nil {
                log.Fatalln("Error reading file", filename)
        }

        // content returned as []byte not string
	lines := strings.Split(string(content), "\n")
        m := make(map[string]string)
        keys := make([]string, 0, len(lines))
        var cookbook, version string;
        for _, line := range lines {
                if ! strings.Contains(line, "=") {
                        continue
                }
		line = strings.Replace(line, "\"", "", -1)
		line = strings.Replace(line, ",", "", -1)
		line = strings.Replace(line, "cookbook", "", -1)
                //fmt.Println(line)
                values := strings.Split(line, "=")
		/*
                for _, value := range values {
                        fmt.Printf("value: %s\n", value)
                }
		*/
                cookbook = strings.Trim(values[0], ": ")
                version = strings.Trim(values[1], " ")
                m[cookbook] = version
                keys = append(keys, cookbook)
        }
        if verbose {
                sort.Strings(keys)
                for _, key := range keys {
                        fmt.Printf("cookbook \"%s\" has version \"%s\"\n", key, m[key])
                }
        }
        return m
}

func getGitDiff(environment *string) ([]byte, error) {
        filename := fmt.Sprintf("/Users/grig.gheorghiu/chef-repo/environments/%s.rb", *environment)
	fmt.Println("\nGetting git diff for ", filename)
	cmd := exec.Command("git", "diff", "--exit-code", filename)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	combined_output, err := cmd.CombinedOutput()
	return combined_output, err
}
		

func main() {
	environment := flag.String("env", "prod", "the environment we are inspecting")
	flag.Parse()
	from_chef := getChefEnvironment(environment, false)
	/*
	for cookbook, version := range from_chef {
		fmt.Printf("Chef cookbook \"%s\" has version \"%s\"\n", cookbook, version)
	}
	*/
	from_local := getLocalEnvironment(environment, false)
	/*
	for cookbook, version := range from_local {
		fmt.Printf("Local cookbook \"%s\" has version \"%s\"\n", cookbook, version)
	}
	*/
	found := false
	for key := range from_local {
		if from_chef[key] != from_local[key] {
			fmt.Printf("Found a difference for %s! Local: %s Chef: %s\n", key, from_local[key], from_chef[key])
			found = true
		}
	}
	if !found {
		fmt.Println("Found no difference between Chef and local.")
	}
	output, err := getGitDiff(environment)
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
	} else {
		fmt.Println("Found no difference between GitHub and local.")
	}
}	
