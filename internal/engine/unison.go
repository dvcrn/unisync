package engine

import (
	"fmt"
	"os/exec"
)

// example: unison  '/Users/david/.config/appconfigsync/Dash' ~/Library/Application\ Support/Dash -ignore "Path */docSets/*" -ignore "Name .DS_Store" -path "library.dash" -path "License" -path "User Contributed" -watch -ignore "Name *.docset" -auto -force
// -prefer xxx        choose this replica's version for conflicting changes
// -force xxx         force changes from this replica to the other

type Unison struct {
	binaryPath string
}

func NewUnison() *Unison {
	return &Unison{
		binaryPath: "unison",
	}
}

func (unison *Unison) Sync(pathA, pathB, prefer string, filesToSync []string, filesToIgnore []string) error {
	args := []string{
		pathA,
		pathB,
		"-batch", // batch mode: ask no questions at all
		"-auto",  // auto accept non conflicting
		"-contactquietly",
		"-prefer", // on conflict, prefer A
		prefer,
	}

	for _, file := range filesToSync {
		args = append(args, "-path")
		args = append(args, file)
	}

	for _, file := range filesToIgnore {
		args = append(args, "-ignore")
		args = append(args, file)
	}

	cmd := exec.Command(unison.binaryPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return err
	}

	fmt.Printf("%s\n", output)
	return nil
}

func (unison *Unison) ForceSync(pathA, pathB string, filesToSync []string, filesToIgnore []string, forceDirection string) error {
	args := []string{
		pathA,
		pathB,

		// append force arg
		"-batch", // batch mode: ask no questions at all
		"-auto",
		"-contactquietly",
		"-force", // force changes from this replica to the other
		forceDirection,
	}

	for _, file := range filesToSync {
		args = append(args, "-path")
		args = append(args, file)
	}

	for _, file := range filesToIgnore {
		args = append(args, "-ignore")
		args = append(args, file)
	}

	cmd := exec.Command(unison.binaryPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return err
	}
	fmt.Printf("%s\n", output)
	return nil
}

func (unison *Unison) SyncAToB(pathA, pathB string, filesToSync []string, filesToIgnore []string) error {
	return unison.ForceSync(pathA, pathB, filesToSync, filesToIgnore, pathA)
}

func (unison *Unison) SyncBToA(pathA, pathB string, filesToSync []string, filesToIgnore []string) error {
	return unison.ForceSync(pathA, pathB, filesToSync, filesToIgnore, pathB)
}
