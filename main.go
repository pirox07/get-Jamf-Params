package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	jamf "github.com/pirox07/jamf-pro-go"
	yaml "github.com/goccy/go-yaml"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

const (
	OutDir = "out_conf"
	configFileName = "conf.yml"
)

type Output struct {
	Policy  jamf.Policy  `yaml:"policy" xml:"policy"`
	Scripts jamf.Scripts `yaml:"scripts" json:"scripts"`
}


func main() {
	// server connection information
	url := os.Getenv("JAMF_BASE_URL")
	userName := os.Getenv("JAMF_USER")
	password := os.Getenv("JAMF_USER_PASSWORD")

	// create http client
	conf, err := jamf.NewConfig(url, userName, password)
	if err != nil {
		fmt.Println("err: ", err.Error())
	}
	client := jamf.NewClient(conf)

	// get policies: size (total count), policy ID, policy name
	policies, err := client.GetPolicies()
	if err != nil {
		fmt.Println("err: ", err.Error())
	}

	// create directory
	if _, err := os.Stat(OutDir); !os.IsNotExist(err)  {
		// directory is exist
		if err = os.RemoveAll(OutDir); err != nil {
			fmt.Println(err)
		}
	}
	if err := os.Mkdir(OutDir, 0775); err != nil {
		fmt.Println(err)
	}

	// output policy and script parameters to a YAML file
	for j := 0; j < int(policies.Size); j++ {
		var output Output

		policyID := policies.Policy[j].ID
		policyName := policies.Policy[j].Name
		// escape "/" in file name
		cnvPolicyName := strings.Replace(policyName, "/", "-", -1)
		// create directory for output
		dirPolicy := "policyID_" + fmt.Sprint(policyID) + "_" + cnvPolicyName
		targetDir := path.Join(OutDir, dirPolicy)
		if err := os.Mkdir(targetDir, 0775); err != nil {
			fmt.Println(err)
		}

		// get policy parameters
		policy, err := client.GetPolicy(uint32(policyID))
		if err != nil {
			fmt.Println("err: ", err.Error())
		}
		policyXML, err := xml.Marshal(&policy)
		err = xml.Unmarshal(policyXML, &output.Policy)
		//err = xml.Unmarshal(policyXML, &output)
		if err != nil {
			fmt.Println("error:", err)
		}

		output.Scripts.TotalCount = output.Policy.Scripts.Size
		// get script
		for i := 0; i < int(output.Policy.Scripts.Size); i++ {
			script, err := client.GetScript(output.Policy.Scripts.PolicyScript[i].ID)
			if err != nil {
				fmt.Println("err: ", err.Error())
			}
			scriptjson, err := json.Marshal(&script)

			var ss = jamf.Script{}
			err = json.Unmarshal(scriptjson, &ss)
			if err != nil {
				fmt.Println("error:", err)
			}

			err = WriteScriptContent(targetDir, output.Policy.Scripts.PolicyScript[i].Name, ss.ScriptContents)
			if err != nil {
				fmt.Println("[err] ", err)
			}

			ss.ScriptContents = "(Look at the script file.)"
			output.Scripts.Results = append(output.Scripts.Results, ss)
		}
		err = WriteConfig(targetDir, output)
	}
}


func WriteConfig (dirName string, config Output) error {
	// output YAML file
	f, err := os.OpenFile(path.Join(dirName, configFileName) , os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0664)
	if err != nil {
		fmt.Println("err: ", err.Error())
	}
	defer f.Close()
	d := yaml.NewEncoder(f)
	if err := d.Encode(config); err != nil {
		log.Fatal(err)
	}
	d.Close()

	return nil
}

func WriteScriptContent (dirName, scriptName, scriptContent string) error {
	// output YAML file
	arr := []string{}
	arr = strings.Split(scriptContent,"")

	b := []byte{}
	for _, line := range arr {
		ll := []byte(line)
		for _, l := range ll {
			b = append(b, l)
		}
	}

	err := ioutil.WriteFile(dirName + "/" + scriptName, b , 0666)
	if err != nil {
		fmt.Println(os.Stderr, err)
		os.Exit(1)
	}

	return nil
}