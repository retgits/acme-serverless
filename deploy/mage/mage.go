//+build mage

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

var (
	stage       = envWithFallback("STAGE", "dev")
	service     = envWithFallback("SERVICE", "payment")
	serviceType = envWithFallback("TYPE", "sqs")
	author      = envWithFallback("AUTHOR", "retgits")
	team        = envWithFallback("TEAM", "vcs")
	s3bucket    = envWithFallback("AWS_S3_BUCKET", "myS3Bucket")
	workingDir  = getwd()
)

var lambdas = map[string][]string{
	"payment-eventbridge":  []string{"lambda-payment-eventbridge"},
	"payment-sqs":          []string{"lambda-payment-sqs"},
	"shipment-eventbridge": []string{"lambda-shipment-eventbridge"},
	"shipment-sqs":         []string{"lambda-shipment-sqs"},
	"cart":                 []string{"lambda-cart-additem", "lambda-cart-all", "lambda-cart-clear", "lambda-cart-itemmodify", "lambda-cart-itemtotal", "lambda-cart-modify", "lambda-cart-total", "lambda-cart-user"},
	"catalog":              []string{"lambda-catalog-all", "lambda-catalog-get", "lambda-catalog-newproduct"},
	"order-eventbridge":    []string{"lambda-order-all", "lambda-order-users", "lambda-order-eventbridge-add", "lambda-order-eventbridge-ship", "lambda-order-eventbridge-update"},
	"order-sqs":            []string{"lambda-order-all", "lambda-order-users", "lambda-order-sqs-add", "lambda-order-sqs-ship", "lambda-order-sqs-update"},
	"user":                 []string{"lambda-user-all", "lambda-user-get", "lambda-user-login", "lambda-user-login", "lambda-user-refreshtoken", "lambda-user-register", "lambda-user-verifytoken"},
}

// envWithFallback retrieves the value of the environment variable named by the key.
// If no value is found, the fallback value is returned.
func envWithFallback(key string, fallback string) string {
	val, found := os.LookupEnv(key)
	if !found {
		return fallback
	}
	return val
}

// getwd returns a rooted path name corresponding to the current directory.
// If an error occurs, the program is halted.
func getwd() string {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Unable to get working directory: %s", err.Error())
	}

	return workingDir
}

// gitVersion collects the current commit version of the service
func gitVersion() string {
	env := make(map[string]string)
	env["WD"] = path.Join(workingDir, "..", "cloudformation", service, fmt.Sprintf("acme-serverless-%s", service))

	v, _ := runCmd(env, "git", "describe", "--tags", "--always", "--dirty=-dev")
	if len(v) == 0 {
		v = "dev"
	}
	return v
}

// runCmd starts the specified command and waits for it to complete.
func runCmd(env map[string]string, cmd string, args ...string) (string, error) {
	buf := &bytes.Buffer{}

	c := exec.Command(cmd, args...)

	if val, ok := env["WD"]; ok {
		c.Dir = val
		delete(env, "WD")
	}

	c.Env = os.Environ()
	for k, v := range env {
		c.Env = append(c.Env, k+"="+v)
	}

	c.Stderr = os.Stderr
	c.Stdout = buf
	c.Stdin = os.Stdin

	err := c.Run()

	return strings.TrimSuffix(buf.String(), "\n"), err
}

// Get performs a git clone of the source code from GitHub for the service specified.
func Get() error {
	repo := fmt.Sprintf("https://github.com/retgits/acme-serverless-%s", service)

	env := make(map[string]string)
	env["WD"] = path.Join(workingDir, "..", "cloudformation", service)

	res, err := runCmd(env, "git", "clone", repo)
	log.Println(res)
	return err
}

// Deps resolves and downloads dependencies to the current development module and then builds and installs them.
// Deps will rely on the Go environment variable GOPROXY (go env GOPROXY) to determine from where to obtain the
// sources for the build.
func Deps() error {
	env := make(map[string]string)
	env["WD"] = path.Join(workingDir, "..", "cloudformation", service, fmt.Sprintf("acme-serverless-%s", service))

	res, err := runCmd(env, "go", "get", "./...")
	log.Println(res)
	return err
}

// 'Go test' automates testing the packages named by the import paths. go:test compiles and tests each of the
// packages listed on the command line. If a package test passes, go test prints only the final 'ok' summary
// line.
func Test() error {
	env := make(map[string]string)
	env["WD"] = path.Join(workingDir, "..", "cloudformation", service, fmt.Sprintf("acme-serverless-%s", service))

	res, err := runCmd(env, "go", "test", "./...")
	log.Println(res)
	return err
}

// Vuln uses Snyk to test for any known vulnerabilities in go.mod. The command relies on access to the Snyk.io
// vulnerability database, so it cannot be used without Internet access.
func Vuln() error {
	env := make(map[string]string)
	env["WD"] = path.Join(workingDir, "..", "cloudformation", service, fmt.Sprintf("acme-serverless-%s", service))

	res, err := runCmd(env, "snyk", "test")
	log.Println(res)
	return err
}

// Build compiles the individual commands in the cmd folder, along with their dependencies. All built executables
// are stored in the 'bin' folder. Specifically for deployment to AWS Lambda, GOOS is set to linux and GOARCH is
// set to amd64.
func Build() error {
	env := make(map[string]string)
	env["GOOS"] = "linux"
	env["GOARCH"] = "amd64"
	env["WD"] = path.Join(workingDir, "..", "cloudformation", service, fmt.Sprintf("acme-serverless-%s", service))

	funcs := lambdas[fmt.Sprintf("%s-%s", service, serviceType)]
	if len(funcs) == 0 {
		funcs = lambdas[service]
	}

	for _, lambda := range funcs {
		res, err := runCmd(env, "go", "build", "-o", path.Join("..", "bin", lambda), fmt.Sprintf("./cmd/%s", lambda))
		if err != nil {
			log.Printf("error building %s: %s", lambda, err.Error())
		}
		log.Println(res)
	}
	return nil
}

// Clean removes object files from package source directories.
func Clean() error {
	env := make(map[string]string)
	env["WD"] = path.Join(workingDir, "..", "cloudformation", service)

	res, err := runCmd(env, "rm", "-rf", "bin")
	log.Println(res)
	return err
}

// Deploy packages, deploys, and returns all outputs of your stack. Packages the local artifacts (local paths) that your
// AWS CloudFormation template references and uploads  local  artifacts to an S3 bucket. The command returns a copy of your
// template, replacing references to local artifacts with the S3 location where the command uploaded the artifacts. Deploys
// the specified AWS CloudFormation template by creating and then executing a change set. The command terminates after AWS
// CloudFormation executes  the change set. Returns the description for the specified stack.
func Deploy() error {
	version := gitVersion()

	template := fmt.Sprintf("lambda-%s-template.yaml", serviceType)

	if _, err := os.Stat(path.Join(workingDir, "..", "cloudformation", service, template)); err != nil {
		template = "lambda-template.yaml"
	}

	env := make(map[string]string)
	env["WD"] = path.Join(workingDir, "..", "cloudformation", service)

	res, err := runCmd(env, "aws", "cloudformation", "package", "--template-file", template, "--output-template-file", "lambda-packaged.yaml", "--s3-bucket", s3bucket)
	log.Println(res)
	if err != nil {
		return err
	}

	res, err = runCmd(env, "aws", "cloudformation", "deploy", "--template-file", "lambda-packaged.yaml", "--stack-name", fmt.Sprintf("%s-%s", service, stage), "--capabilities", "CAPABILITY_IAM", "--parameter-overrides", fmt.Sprintf("Version=%s", version), fmt.Sprintf("Author=%s", author), fmt.Sprintf("Team=%s", team))
	log.Println(res)
	if err != nil {
		return err
	}

	res, err = runCmd(env, "aws", "cloudformation", "describe-stacks", "--stack-name", fmt.Sprintf("%s-%s", service, stage), "--query", "'Stacks[].Outputs'")
	log.Println(res)
	if err != nil {
		return err
	}

	return nil
}
