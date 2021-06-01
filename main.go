package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func validateFileExists(filePath string) error {
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist", filePath)
	}
	if fileInfo.IsDir() {
		return fmt.Errorf("%s is a directory", filePath)
	}
	return nil
}

func validateFileArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return errors.New("<file1> <file2> not provided")
	}
	file1, file2 := args[0], args[1]
	if err := validateFileExists(file1); err != nil {
		return err
	}
	if err := validateFileExists(file2); err != nil {
		return err
	}
	return nil
}

type setOperation func(file1, file2 string) error

func buildCobraCmd(setOperationFunc setOperation) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		file1, file2 := args[0], args[1]
		return setOperationFunc(file1, file2)
	}
}

func getLineLookup(filePath string) (map[string]struct{}, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	lookup := make(map[string]struct{})
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lookup[scanner.Text()] = struct{}{}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lookup, nil
}

func intersection(file1, file2 string) error {
	file2Lines, err := getLineLookup(file2)
	if err != nil {
		return err
	}

	f, err := os.Open(file1)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if _, ok := file2Lines[scanner.Text()]; ok {
			fmt.Println(scanner.Text())
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func difference(file1, file2 string) error {
	file2Lines, err := getLineLookup(file2)
	if err != nil {
		return err
	}

	f, err := os.Open(file1)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if _, ok := file2Lines[scanner.Text()]; !ok {
			fmt.Println(scanner.Text())
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setop",
		Short: "Perform set operations on files.",
	}
	intersectCmd := &cobra.Command{
		Use:     "intersection <file1> <file2>",
		Aliases: []string{"intersect"},
		Short:   "(A & B) - Output lines in both file1 and file2",
		Args:    validateFileArgs,
		RunE:    buildCobraCmd(intersection),
	}
	diffCmd := &cobra.Command{
		Use:     "difference <file1> <file2>",
		Aliases: []string{"diff"},
		Short:   "(A - B) - Output lines in file1 but not file2",
		Args:    validateFileArgs,
		RunE:    buildCobraCmd(difference),
	}
	cmd.AddCommand(intersectCmd)
	cmd.AddCommand(diffCmd)
	cmd.SilenceUsage = true
	return cmd
}

func main() {
	if err := NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
