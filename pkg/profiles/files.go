package profiles

import (
	"embed"
	"os"
	"path"

	"github.com/harness/smp-installer/pkg/util"
	"github.com/rs/zerolog/log"
)

//go:embed small
var SmallProfileFiles embed.FS

//go:embed common
var CommonProfileFiles embed.FS

//go:embed pov
var PovProfileFiles embed.FS

var filesMap = map[string]embed.FS{
	"small":  SmallProfileFiles,
	"common": CommonProfileFiles,
	"pov":    PovProfileFiles,
}

func CopyFiles(profile string, outDir string, filenames []string) ([]string, error) {
	os.MkdirAll(outDir, 0777)
	profileFiles := filesMap[profile]
	filesCopied := make([]string, 0)
	files, err := profileFiles.ReadDir(profile)
	if err != nil {
		log.Err(err).Msgf("could not find %s directory", profile)
		return nil, err
	}
	for _, f := range files {
		skipFile := false
		log.Info().Msgf("file to copy: %s", f.Name())
		if len(filenames) == 0 || !util.Contains(filenames, f.Name()) {
			skipFile = true
		}
		if !skipFile {
			data, err := profileFiles.ReadFile(path.Join(profile, f.Name()))
			if err != nil {
				log.Err(err).Msgf("cannot read file %s", f.Name())
				return nil, err
			}
			outputFileName := path.Join(outDir, profile+"-"+f.Name())
			err = os.WriteFile(outputFileName, data, 0666)
			if err != nil {
				log.Err(err).Msgf("failed to copy file %s to directory %s", f.Name(), outDir)
				return nil, err
			}
			filesCopied = append(filesCopied, outputFileName)
		} else {
			log.Info().Msgf("skipping file %s", f.Name())
		}
	}
	return filesCopied, nil
}

func CopyOverrideFiles(profile string, outDir string) ([]string, error) {
	filesCopied, err := CopyFiles(profile, outDir, []string{"override.yaml"})
	if err != nil {
		return nil, err
	}
	commonFilesCopied, err := CopyFiles("common", outDir, []string{"override.yaml"})
	if err != nil {
		return nil, err
	}
	return append(commonFilesCopied, filesCopied...), nil
}

func CopyInstallerFiles(profile string, outDir string) ([]string, error) {
	filesCopied, err := CopyFiles(profile, outDir, []string{"config.yaml"})
	if err != nil {
		return nil, err
	}
	commonFilesCopied, err := CopyFiles("common", outDir, []string{"config.yaml"})
	if err != nil {
		return nil, err
	}
	return append(commonFilesCopied, filesCopied...), nil
}
