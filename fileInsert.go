package main

import (
    "bufio"
    "io"
  	"fmt"
    "io/ioutil"
    "os"
  	"./filepop"
)

// File2lines function to get reader for file.
func File2lines(filePath string) ([]string, error) {
    f, err := os.Open(filePath)
    check(err)
    defer f.Close()
    return LinesFromReader(f)
}

// LinesFromReader function to read lines from reader to an array of strings.
func LinesFromReader(r io.Reader) ([]string, error) {
    var lines []string
    scanner := bufio.NewScanner(r)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return lines, nil
}

// InsertStringToFile function is used to Insert sting to n-th line of file.
// If you want to insert a line, append newline '\n' to the end of the string.
func InsertStringToFile(path, str string, index int) error {
    lines, err := File2lines(path)
    if err != nil {
        return err
    }

    fileContent := ""
    for i, line := range lines {
        if i == index {
            fileContent += str
        }
        fileContent += line
        fileContent += "\n"
    }

    return ioutil.WriteFile(path, []byte(fileContent), 0644)
}

/*source               = "../../../../modules/aws/compute/ec2"
count_instances      = 2
instance_type        = "t2.medium"*/

func main() {

  filename := "./file2.tf"

  ec2statefullines := []string{
    "module \"ec2_service\" {",
    "  source               = \"../../../../modules/aws/compute/ec2\"",
    "  count_instances      = 2",
    "  instance_type        = \"t2.medium\""}

  // Deleting lines related to autoscaling.
  for i:=0; i<7; i++ {
   filepop.Pop(filename)
  }

  // Adding lines for static instances.
  for j:=3; j>=0; j-- {
        err := InsertStringToFile(file_name, ec2_stateful_lines[j]+"\n", 0)
        if err != nil {
              fmt.Println(err)
        }
   }
}