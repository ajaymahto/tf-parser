// This code does the following actions:
// - Changes the mentioned file (ec2.tf) to convert the service to stateful.
//   - Not using templates as there is no standard library for the same in golang for HCL.
//   - Steps: 
//   - 1. Comment out the ec2_auto parameters.
//   - 2. Add the stateful parameters with appropriate values.
// - Applies the changes incrementally using terraform apply.

package main

import(
    "fmt"
	"io/ioutil"
	"bytes"
	"os"
	"os/exec"
)
func check(err error) {
  if err != nil {
	  fmt.Println(err)
	  os.Exit(1)
  }
}

// Function to populate stateful params with appropriate values.
func createStatefulParams(srv, ec2_count string) []byte {
	statefulParams := fmt.Sprintf("module \"ec2_%s\" {\n", srv)
	statefulParams += fmt.Sprintf("source = \"../../../../modules/aws/compute/ec2\"\n")
	statefulParams += fmt.Sprintf("count_instances = \"%s\"\n", ec2_count)
	return []byte(statefulParams)
}

// Function to modify tf config to change resource from stateless to stateful.
func statelessToStateful(tfConfig []byte, statefulParams []byte) []byte {
	statelessParams := []string{"source", "asg_name", "launchconfig_prefix", "load_balancers", "min_size", "desired_capacity"}
	configParams := tfConfig

	// Commenting out all lines related to ec2_auto.
	// Commenting module line.
	configParams = bytes.Replace(configParams, []byte("module "), []byte("# module "), -1)

	// Commenting autoscaling(stateless) params.
	for _, param := range statelessParams {
		param += " "
		commented := "# " + param
		configParams = bytes.Replace(configParams, []byte(param), []byte(commented), -1)
	}

	// Add params for ec2.
	result := append(statefulParams, configParams...)
	return result
}

// Function to load tf config from file
func readTfConfig(filepath string) []byte {
	data, err := ioutil.ReadFile("file1.tf")
	check(err)
	return data
}

func writeTfConfig(filepath string, tfConfig []byte) error {
   	// the WriteFile method returns an error if unsuccessful
	err := ioutil.WriteFile(filepath, tfConfig, 0664)
	check(err)
	return nil
}

// Function to run shell commands from go code.
func runCmd(cmd *exec.Cmd) string {
	output, err := cmd.Output()
	check(err)
	return string(output)
}

// Function to apply terraform changes. 
// Need to confirm the steps with Grab TF conventions.
func tfApply() error {

	// To run terraform plan in automated way the following sequence of steps needs to be followed.
	// 1. terraform init -input=false [Verify with Grab TF conventions]
	// 2. terraform plan -input=false -out=tfplan
	// 3. terraform apply -input=false -auto-approve tfplan

	// Specify the exact binary location for terraform.
	app := "terraform"
	
	// terraform init
	cmd := exec.Command(app, "init")
	output := runCmd(cmd)
	fmt.Println(output)

	// terraform plan -out=tfplan -input=false
	cmd = exec.Command(app, "plan", "-out=tfplan", "-input=false")
	output = runCmd(cmd)
	fmt.Println(output)

	// terraform apply -input=false -auto-approve tfplan
	cmd = exec.Command(app, "apply", "-input=false", "-auto-approve", "tfplan")
	output = runCmd(cmd)
	fmt.Println(output)

	return nil
}

func main() {

	srv := "grabpay-osquery"
	ec2_count := "2"

	// Read ec2.tf filepath.
	tfFile := "file1.tf"

	// Read tf file.
	tfConfig := readTfConfig(tfFile)

	// Create stateful params using srv, ec2_count values.
	statefulParams := createStatefulParams(srv, ec2_count)

	// Add modifications to change from stateless to stateful.
	statefulTfConfig := statelessToStateful(tfConfig, statefulParams)

	// Write to tf file.
	writeTfConfig(tfFile, statefulTfConfig)

	// Apply the changes through terraform apply.
	tfApply()
}