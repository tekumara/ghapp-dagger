package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/sync/errgroup"
)

func findDockerPath() (*string, error) {
	cmd := exec.Command("which", "docker")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	dir := filepath.Dir(string(out))
	return &dir, nil
}

// TODO: replace updateOutput func with IO writer?
func execDagger(ctx context.Context, repoUrl, ref, token string, updateOutput func(text string) error) (string, error) {
	g, gctx := errgroup.WithContext(ctx)

	cmd := exec.CommandContext(gctx, "dagger", "do", "build", "--log-format=plain")

	log.Printf("cmd %+v\n", cmd)

	dockerPath, err := findDockerPath()
	if err != nil {
		log.Printf("ERROR: findDockerPath() failed with %s\n", err)
		return err.Error(), err
	}

	cmd.Env = []string{
		// dagger needs ~/.config/dagger/version-check
		fmt.Sprintf("HOME=%s", os.Getenv("HOME")),

		// dagger needs docker on the path
		fmt.Sprintf("PATH=%s", *dockerPath),

		// repo details
		fmt.Sprintf("GITHUB_REPO_URL=%s", repoUrl),
		fmt.Sprintf("GITHUB_REF=%s", ref),
	}

	log.Printf("cmd.Env %+v\n", cmd.Env)

	// add token separately so it's not printed 
	cmd.Env = append(cmd.Env, fmt.Sprintf("GITHUB_TOKEN=%s", token))

	combined := &bytes.Buffer{}

	stderr := &bytes.Buffer{}
	cmd.Stderr = io.MultiWriter(os.Stderr, stderr, combined)

	stdoutr, stdoutw := io.Pipe()
	cmd.Stdout = io.MultiWriter(os.Stdout, stdoutw, combined)

	g.Go(func() error {
		defer stdoutw.Close()
		return cmd.Run()
	})

	// TODO: do we need errgroup or could this just be a goroutine?
	g.Go(func() error {
		scan := bufio.NewScanner(stdoutr)
		for scan.Scan() {
			log.Println(scan.Bytes())
			// TODO: flush every X lines to github
		}
		return nil
	})

	if cmdErr := g.Wait(); cmdErr != nil {
		return StripAnsi(combined.String()), cmdErr
	}

	return StripAnsi(combined.String()), nil
}
