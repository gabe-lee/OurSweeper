package main

// import (
// 	"errors"
// 	"fmt"
// 	"log"
// 	"os"
// 	"os/exec"
// 	"strings"
// )

// func inputYes(input string) bool {
// 	return strings.HasPrefix(strings.ToLower(strings.TrimSpace(input)), "y")
// }

// type installer struct {
// 	cmd     string
// 	found   bool
// 	updated bool
// }

// var installerList = []string{
// 	"nala",
// 	"apt",
// }

// func (i *installer) findInstaller() error {
// 	if i.found {
// 		return nil
// 	}
// 	for _, inst := range installerList {
// 		cmd, err := exec.LookPath(inst)
// 		if err == nil {
// 			i.cmd = cmd
// 			i.found = true
// 			return nil
// 		}
// 	}
// 	return errors.New(fmt.Sprintf("No installer found from list: %v", installerList))
// }

// func (i *installer) install(pkgs ...string) error {
// 	var err error
// 	if !i.found {
// 		err = i.findInstaller()
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	if !i.updated {
// 		updateCmd := exec.Command(i.cmd, "update")
// 		updateCmd.Stdout = os.Stdout
// 		updateCmd.Stderr = os.Stderr
// 		fmt.Printf("RUNNING COMMAND: %s", updateCmd.String())
// 		err = updateCmd.Run()
// 		if err != nil {
// 			return err
// 		}
// 		i.updated = true
// 	}
// 	pkgsWithInst := append([]string{"install"}, pkgs...)
// 	instCmd := exec.Command(i.cmd, pkgsWithInst...)
// 	instCmd.Stdout = os.Stdout
// 	instCmd.Stderr = os.Stderr
// 	fmt.Printf("RUNNING COMMAND: %s", instCmd)
// 	return instCmd.Run()
// }

// func (i *installer) installCmd(pkgs ...string) *exec.Cmd {
// 	fullArgs := append([]string{"install"}, pkgs...)
// 	return exec.Command(i.cmd, fullArgs...)
// }

// func (i *installer) updateCmd() *exec.Cmd {
// 	return exec.Command(i.cmd, "update")
// }

// const UPDATE = "%s update"
// const INSTALL = "%s install"
// const PSQL_INSTALL = "sudo apt update && sudo apt install postgresql postgresql-contrib"

// func logPsqlImpossibleIfErr(err error) {
// 	if err != nil {
// 		log.Fatalf("%w\nDependancy `psql` (PostgresQL database) not found on system PATH, please download PostgresQl for your system:\n\tsudo apt update && sudo apt install postgresql postgresql-contrib\n`apt` may be swapped for another compatible package manager, else install according to the official instructions: https://www.postgresql.org/download/\n", err)
// 	}
// }

// func init_local_project() int {
// 	var input string
// 	inst := installer{}
// 	fmt.Println("Checking system dependencies...")
// 	psql, err := exec.LookPath("psql")
// 	if err != nil {
// 		err = inst.findInstaller()
// 		logPsqlImpossibleIfErr(err)
// 		cmd := inst.installCmd("postgresql", "postgresql-contrib")
// 		fmt.Printf("Dependancy `psql` (PostgresQL database) not found on system PATH\nInstall command:\n\t%s\nWould you like to install automatically? (y/n):", cmd.String())
// 		fmt.Scanln(&input)
// 		if inputYes(input) {
// 			logPsqlImpossibleIfErr(cmd.Run())
// 		} else {
// 			logPsqlImpossibleIfErr(errors.New("Declined automatic dependency install"))
// 		}
// 	}
// 	goose, err := exec.LookPath("goose")
// 	if err != nil {
// 		log.Fatalf("command `goose` does not exist on the system, please download goose for your system:\n\tgo install github.com/pressly/goose/v3/cmd/goose@latest\n")
// 	}
// }
