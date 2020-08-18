package util

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CmdFunc func(reader *bufio.Reader) error

type Cmd struct {
	Name    string
	Desc    string
	CmdFunc CmdFunc

	Cmds []*Cmd
}

func NewCmd(name, desc string, cmdFunc CmdFunc, cmds ...*Cmd) *Cmd {
	return &Cmd{
		Name:    name,
		Desc:    desc,
		CmdFunc: cmdFunc,
		Cmds:    append([]*Cmd{}, cmds...),
	}
}

func (cmd *Cmd) Append(cmds ...*Cmd) {
	cmd.Cmds = append(cmd.Cmds, cmds...)
}

func (cmd *Cmd) Run() error {
	reader := bufio.NewReader(os.Stdin)
	cmds := cmd.getCmds()

	for {
		for i := 0; i < len(cmd.Name)+len(cmd.Desc)+7; i++ {
			fmt.Printf("#")
		}
		fmt.Printf("\n")
		fmt.Printf("# %s | %s #\n", cmd.Name, cmd.Desc)
		for i := 0; i < len(cmd.Name)+len(cmd.Desc)+7; i++ {
			fmt.Printf("#")
		}
		fmt.Printf("\n")

		// display cmds
		for i, c := range cmd.Cmds {
			fmt.Printf("%d. %s", i+1, c.Name)

			if c.Desc != "" {
				fmt.Printf(" - %s\n", c.Desc)
			}
		}

		fmt.Printf("%d. %s\n", 0, "exit")

		// user input
		data, err := GetInput(reader, true)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if data == "0" || data == "exit" {
			break
		}

		val, exists := cmds[data]
		if !exists {
			fmt.Printf("'%s' not proper command\n\n", data)
			continue
		}

		fmt.Printf("\n")

		if val.CmdFunc == nil && len(val.Cmds) != 0 {
			val.Run()
		} else if val.CmdFunc != nil {
			err = val.CmdFunc(reader)
			if err != nil {
				fmt.Println(err)
			}
		}

		fmt.Printf("\n")
	}

	return nil
}

func (cmd *Cmd) getCmds() map[string]*Cmd {
	result := map[string]*Cmd{}
	n := 1

	for _, c := range cmd.Cmds {
		nStr := strconv.Itoa(n)
		n += 1

		result[nStr] = c
		result[c.Name] = c
	}

	return result
}

func GetInput(reader *bufio.Reader, guide bool) (string, error) {
	if guide {
		fmt.Printf("\n> ")
	}

	data, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	data = strings.ToLower(data)
	data = strings.TrimRight(data, "\r\n")

	return data, nil
}
