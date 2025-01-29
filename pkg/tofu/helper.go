package tofu

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

type InfraCommand string

const (
	InitCommand  InfraCommand = "init"
	PlanCommand  InfraCommand = "plan"
	ApplyCommand InfraCommand = "apply"
)

func ExecuteCommand(ctx context.Context, outDir string, icmd InfraCommand) error {
	currentDir, _ := os.Getwd()
	os.Chdir(outDir)
	extraArgs := []string{"-auto-approve"}
	cmd := exec.Command("tofu", fmt.Sprintf("%s", icmd), "-var-file=tf.vars", "-no-color")
	if icmd != InitCommand && icmd != PlanCommand {
		cmd.Args = append(cmd.Args, extraArgs...)
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Err(err).Msgf("infra %s command failed with output: %s", icmd, err.(*exec.ExitError).Stderr)
		os.Chdir(currentDir)
		return err
	}

	log.Info().Msgf("infra %s command output", icmd)
	// log.Info().Msg(string(stdout))

	os.Chdir(currentDir)
	return nil
}

func GetOutput(ctx context.Context, outDir, outname, jspath string) (string, error) {
	currentDir, _ := os.Getwd()
	os.Chdir(outDir)
	tmpfilename := fmt.Sprintf("tmpout-%s.json", outname)
	tofucmd := exec.Command("tofu", "output", "-json", outname)
	tofuout, err := tofucmd.Output()
	if err != nil {
		log.Err(err).Msgf("tofu output failed")
	}
	err = os.WriteFile(tmpfilename, tofuout, 0666)
	if err != nil {
		log.Err(err).Msgf("failed to create tofu json output file")
		os.Chdir(currentDir)
		return "", nil
	}
	jqcmd := exec.Command("jq", "-r", jspath, tmpfilename)
	jqout, err := jqcmd.Output()
	if err != nil {
		log.Err(err).Msgf("jq command wait failed: %s", err.(*exec.ExitError).Stderr)
		os.Chdir(currentDir)
		return "", err
	}
	os.Remove(tmpfilename)
	os.Chdir(currentDir)
	return strings.TrimSpace(string(jqout)), nil
}
