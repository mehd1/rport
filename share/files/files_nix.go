//+build !windows

package files

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"

	"github.com/pkg/errors"
)

func ChangeOwner(path, owner, group string) error {
	if owner == "" && group == "" {
		return nil
	}

	targetUserUID := os.Getuid()
	if owner != "" {
		usr, err := user.Lookup(owner)
		if err != nil {
			return err
		}
		targetUserUID, err = strconv.Atoi(usr.Uid)
		if err != nil {
			return err
		}
	}

	targetGroupGUID := os.Getgid()
	if group != "" {
		gr, err := user.LookupGroup(group)
		if err != nil {
			return err
		}
		targetGroupGUID, err = strconv.Atoi(gr.Gid)
		if err != nil {
			return err
		}
	}

	err := os.Chown(path, targetUserUID, targetGroupGUID)
	if err == nil {
		return nil
	}

	if os.IsPermission(err) {
		return ChangeOwnerExecWithSudo(path, owner, group)
	}

	return err
}

func ChangeOwnerExecWithSudo(path, owner, group string) error {
	if owner == "" && group == "" {
		return nil
	}

	args := []string{
		"sudo",
		"-n",
		"chown",
		fmt.Sprintf("%s:%s", owner, group),
		path,
	}

	cmd := exec.Command(args[0], args[1:]...)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return errors.Wrapf(err, "failed to execute %s: %s", cmd.String(), string(output))
	}

	return nil
}

func Rename(oldPath, newPath string) error {
	err := os.Rename(oldPath, newPath)
	if os.IsPermission(err) {
		return MoveExecWithSudo(oldPath, newPath)
	}
	if err != nil {
		return err
	}

	return nil
}

func MoveExecWithSudo(sourcePath, targetPath string) error {
	args := []string{
		"sudo",
		"-n",
		"mv",
		sourcePath,
		targetPath,
	}

	cmd := exec.Command(args[0], args[1:]...)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return errors.Wrapf(err, "failed to execute %s: %s", cmd.String(), string(output))
	}

	return nil
}
