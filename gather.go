// Copyright (c) 2021 Yandex LLC. All rights reserved.
// Author: Martynov Pavel <covariance@yandex-team.ru>

package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"licenser/pkg/util"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	gatherCmd = &cobra.Command{
		Use:   "gather <path to module> <output file>",
		Short: "Transitively gathers licenses from vendor modules",
		Args:  cobra.MinimumNArgs(2),
		RunE:  gather,
	}

	license = regexp.MustCompile("\\b(license|LICENSE|notice|NOTICE|copying|COPYING)(\\.*)?\\b")
)

func init() {
	rootCmd.AddCommand(gatherCmd)
}

func gather(_ *cobra.Command, args []string) error {
	modulePath := filepath.Clean(args[0])
	outPath := filepath.Clean(args[1])

	log.Infof("Gathering licenses for module %s to %s", modulePath, outPath)

	moduleInfo, err := util.Exists(modulePath)
	if err != nil {
		return err
	}
	if moduleInfo == nil || !moduleInfo.IsDir() {
		return fmt.Errorf("module directory does not exist")
	}

	vendor := filepath.Join(modulePath, "vendor")
	vendorInfo, err := util.Exists(vendor)
	if err != nil {
		return err
	}

	if vendorInfo != nil && !vendorInfo.IsDir() {
		return fmt.Errorf("vendor is not a directory")
	}

	if vendorInfo == nil {
		log.Info("Vendor folder not found, trying to execute \"go mod vendor\"...")

		cmd := exec.Command("go", "mod", "vendor")
		cmd.Dir = modulePath
		cmdout, err := cmd.Output()
		defer func() {
			log.Info("cleaning up vendor directory")
			err = os.RemoveAll(vendor)
			if err != nil {
				log.Errorf("error while cleaning up: %v", err)
			}
			log.Info("successfully cleaned up")
		}()
		if err != nil {
			return err
		}
		if len(cmdout) > 0 {
			log.Info(string(cmdout))
		}
		log.Info("Successfully got vendor dependencies, continuing...")
	}

	out, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	err = filepath.Walk(vendor,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if license.MatchString(filepath.Base(path)) {
				log.Infof("license found in: %s", path)
				_, err = out.WriteString(strings.TrimPrefix(strings.TrimSuffix(path, "/"+filepath.Base(path)), vendor+"/") + ": " + filepath.Base(path) + "\n\n")
				if err != nil {
					return err
				}
				licenseContents, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				_, err = out.Write(licenseContents)
				if err != nil {
					return err
				}
				_, err = out.WriteString("\n\n")
				if err != nil {
					return err
				}
			}
			return err
		},
	)
	if err != nil {
		return err
	}
	log.Info("successfully gathered licenses")
	return nil
}
