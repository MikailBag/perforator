package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type sshSessionOptions struct {
	username  string
	server    string
	dir       string
	identity  string
	extraOpts []string
}

func escapePath(path string) (string, error) {
	if strings.ContainsAny(path, "*\\'?+[]{}()$") {
		// TODO: actually escape?
		return "", fmt.Errorf("path %q contains invalid characters", path)
	}

	return path, nil
}

type sftpCommandUpload struct {
	src string
	dst string
}

func (c *sftpCommandUpload) serialize() (string, error) {
	s, err := escapePath(c.src)
	if err != nil {
		return "", fmt.Errorf("failed to escape source path: %w", err)
	}
	d, err := escapePath(c.dst)
	if err != nil {
		return "", fmt.Errorf("failed to escape destination path: %w", err)
	}
	return fmt.Sprintf("put '%s' '%s'", s, d), nil
}

type sftpCommandProgess struct {
}

func (c *sftpCommandProgess) serialize() (string, error) {
	return "progress", nil
}

type sftpCommand struct {
	upload   *sftpCommandUpload
	progress *sftpCommandProgess
}

func commonArgs(opts sshSessionOptions) []string {
	var args []string
	if opts.identity != "" {
		args = append(args, "-i", opts.identity)
	}
	args = append(args, "-o", "StrictHostKeyChecking=no")
	for _, opt := range opts.extraOpts {
		args = append(args, "-o", opt)
	}
	target := opts.server
	if opts.username != "" {
		target = fmt.Sprintf("%s@%s", opts.username, target)
	}
	if opts.dir != "" {
		target = fmt.Sprintf("%s:%s", target, opts.dir)
	}
	args = append(args, target)
	return args
}

func runSSH(ctx context.Context, opts sshSessionOptions, commands []string) error {
	args := commonArgs(opts)
	args = append(args, commands...)
	cmd := exec.CommandContext(ctx, "ssh", args...)
	cmd.Stdin = bytes.NewBufferString(strings.Join(commands, "\n"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runSFTP(ctx context.Context, opts sshSessionOptions, commands []sftpCommand) error {
	/*target := opts.server
	if opts.username != "" {
		target = fmt.Sprintf("%s@%s", opts.username, target)
	}
	if opts.dir != "" {
		target = fmt.Sprintf("%s:%s", target, opts.dir)
	}*/
	var rawCommands []string
	for _, c := range commands {
		var raw string
		var err error
		switch {
		case c.upload != nil:
			raw, err = c.upload.serialize()
		case c.progress != nil:
			raw, err = c.progress.serialize()
		default:
			err = fmt.Errorf("unknown command %+v", c)
		}
		if err != nil {
			return fmt.Errorf("failed to serialize command: %w", err)
		}
		rawCommands = append(rawCommands, raw)
	}
	sftpArgs := append([]string{"-b", "-", "-N"}, commonArgs(opts)...)
	cmd := exec.CommandContext(ctx, "sftp", sftpArgs...)
	cmd.Stdin = bytes.NewBufferString(strings.Join(rawCommands, "\n"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute commands: %w", err)
	}
	return nil
}
