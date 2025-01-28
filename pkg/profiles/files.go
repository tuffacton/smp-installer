package profiles

import (
	"embed"
	"os"
	"path"

	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/util"
	"github.com/rs/zerolog/log"
)

//go:embed small
var SmallProfileFiles embed.FS

var filesMap = map[string]embed.FS{
	"small": SmallProfileFiles,
}

func CopyFiles(profile string, outDir string, filenames []string) error {
	os.MkdirAll(outDir, 0777)
	profileFiles := filesMap[profile]
	files, err := profileFiles.ReadDir(profile)
	if err != nil {
		log.Err(err).Msgf("could not find %s directory", profile)
		return err
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
				return err
			}
			err = os.WriteFile(path.Join(outDir, f.Name()), data, 0666)
			if err != nil {
				log.Err(err).Msgf("failed to copy file %s to directory %s", f.Name(), outDir)
				return err
			}
		} else {
			log.Info().Msgf("skipping file %s", f.Name())
		}
	}
	return nil
}
