package main

import (
	"encoding/json"
	//"flag"
	"fmt"
	"io/ioutil"
	"os"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/printer"
	jsonParser "github.com/hashicorp/hcl/json/parser"
)

type tag struct {
	Name string	`json:Name`
	TagComponentName string	`json:tag_component_name`
	TagCostCenter string	`json:tag_cost_center`
	TagEnv string	`json:tag_env`
	TagJira string	`json:tag_jira`
}

type srv struct {
	CountInstances string	`json:count_instances`
	IAMInstanceProfile string	`json:iam_instance_profile`
	InstanceSubnets []string	`json:instance_subnets`
	InstanceType string	`json:instance_type`
	KeyPair string	`json:key_pair`
	SGIds []string	`json:sg_ids`
	Source string	`json:source`
	Tags []tag	`json:tags`
}

type module struct {
	Service []srv
}

type tf struct {
	Module map[string]module	`json:module` 
}

func main() {
	tfFile := "file1.tf"
	toJSON(tfFile)

	jsonFile := "file1.json"
	toHCL(jsonFile)
}

func toJSON(tfFile string) error {
	input, err := ioutil.ReadFile(tfFile)
	//fmt.Println(input)

	if err != nil {
		return fmt.Errorf("unable to read from stdin: %s", err)
	}

	var v interface{}
	err = hcl.Unmarshal(input, &v)
	if err != nil {
		return fmt.Errorf("unable to parse HCL: %s", err)
	}

	//fmt.Println(v)
	
	json1, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal json: %s", err)
	}

	var v1 tf
	json.Unmarshal(json1, &v1)
	
	fmt.Println(string(json1))
	fmt.Println("Post Marshalling!")
	fmt.Println(v1)
	return nil
}

func toHCL(jsonFile string) error {
	input, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return fmt.Errorf("unable to read from stdin: %s", err)
	}

	ast, err := jsonParser.Parse([]byte(input))
	if err != nil {
		return fmt.Errorf("unable to parse JSON: %s", err)
	}

	err = printer.Fprint(os.Stdout, ast)
	if err != nil {
		return fmt.Errorf("unable to print HCL: %s", err)
	}

	return nil
}