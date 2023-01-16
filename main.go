package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/manifoldco/promptui"
)

const ConfName = ".gitswitch.json"

var (
	// flag prama
	Add    = flag.Bool("n", false, "Add a new git user")
	Del    = flag.Bool("d", false, "Delete an existing user")
	Output = flag.Bool("o", false, "Print shell exec output info")

	// check email regexp
	emailRexp = regexp.MustCompile(`^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,64}$`)
)

type User struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	SSHKeyFilePath string `json:"ssh_key_file_path"`
}

func GetConfig() (users []User, err error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}

	path := filepath.Join(homeDir, ConfName)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = json.Unmarshal(content, &users)
	return users, err
}

func SaveConfig(users []User) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	path := filepath.Join(homeDir, ConfName)

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	content, err := json.MarshalIndent(users, "", "\t")
	if err != nil {
		return err
	}

	if _, err := file.Write(content); err != nil {
		return err
	}

	return nil
}

func UsersFormat(users []User) []string {
	items := make([]string, len(users))

	nameMax := 0
	emailMax := 0
	for _, usr := range users {
		if len(usr.Name) > nameMax {
			nameMax = len(usr.Name)
		}
		if len(usr.Email) > emailMax {
			emailMax = len(usr.Email)
		}
	}

	// "%-20s %-20s %s" content align left
	for i, usr := range users {
		items[i] = fmt.Sprintf("%-"+strconv.Itoa(nameMax+4)+"s %-"+strconv.Itoa(emailMax+4)+"s %s", usr.Name, usr.Email, usr.SSHKeyFilePath)
	}

	return items
}

func ShellExec(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if *Output {
		fmt.Printf("> %v\n", cmd)
	}
	if err := cmd.Run(); err != nil {
		fmt.Printf("err: %v", err)
		return err
	}

	return nil
}

func SwitchUser() error {
	users, err := GetConfig()
	if err != nil {
		return err
	}

	if len(users) == 0 {
		fmt.Println("no users found")
	}

	action := promptui.Select{
		Label: "Select a switch git user",
		Items: UsersFormat(users),
	}

	index, _, err := action.Run()
	if err != nil {
		return err
	}

	option := "--global"

	if err := ShellExec("ssh-add", "-D"); err != nil {
		return err
	}
	if err := ShellExec("ssh-add", users[index].SSHKeyFilePath); err != nil {
		return err
	}
	if err := ShellExec("git", "config", option, "user.name", users[index].Name); err != nil {
		return err
	}
	if err := ShellExec("git", "config", option, "user.email", users[index].Email); err != nil {
		return err
	}

	return nil
}

func AddUser() error {
	inputName := promptui.Prompt{
		Label: "Enter git user name",
	}

	name, err := inputName.Run()
	if err != nil {
		return err
	}

	inputEmail := promptui.Prompt{
		Label: "Enter git user email",
		Validate: func(email string) error {
			if !emailRexp.MatchString(email) {
				return errors.New("invalid email address")
			}
			return nil
		},
	}

	email, err := inputEmail.Run()
	if err != nil {
		return err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	files, err := os.ReadDir(fmt.Sprintf("%v/.ssh/", homeDir))
	if err != nil {
		return err
	}

	var fs []string
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		fs = append(fs, f.Name())
	}

	action := promptui.Select{
		Label: "select ssh key file",
		Items: fs,
	}
	_, sshKeyPath, err := action.Run()
	if err != nil {
		return err
	}

	usr := User{
		Name:           name,
		Email:          email,
		SSHKeyFilePath: fmt.Sprintf("%v/.ssh/%v", homeDir, sshKeyPath),
	}

	users, err := GetConfig()
	if os.IsNotExist(err) {
		users = []User{}
	}

	users = append(users, usr)
	if err = SaveConfig(users); err != nil {
		return err
	}

	return nil
}

func DelUser() error {
	users, err := GetConfig()
	if err != nil {
		return err
	}

	if len(users) == 0 {
		return errors.New("no users found")
	}

	action := promptui.Select{
		Label: "Select a delete git user",
		Items: UsersFormat(users),
	}

	index, _, err := action.Run()
	if err != nil {
		return err
	}

	users = append(users[:index], users[index+1:]...)
	if err = SaveConfig(users); err != nil {
		return err
	}

	return nil
}

func main() {
	flag.Parse()

	switch {
	case *Add:
		if err := AddUser(); err != nil {
			fmt.Printf("add a git user, err: %v\n", err)
		}
	case *Del:
		if err := DelUser(); err != nil {
			fmt.Printf("delete a git user, err: %v\n", err)
		}
	default:
		if err := SwitchUser(); err != nil {
			fmt.Printf("could not select user, err: %v\n", err)
		}
	}
}
