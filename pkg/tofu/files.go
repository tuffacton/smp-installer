package tofu

import (
	"embed"
	"os"
	"path"

	"github.com/rs/zerolog/log"
)

//go:embed aws
var AWSTofuFiles embed.FS

func CopyFiles(srcDir string, outDir string) error {
	os.MkdirAll(outDir, 0777)
	files, err := AWSTofuFiles.ReadDir(srcDir)
	if err != nil {
		log.Err(err).Msgf("could not find %s directory", srcDir)
		return err
	}
	for _, f := range files {
		if f.IsDir() {
			os.MkdirAll(path.Join(outDir, f.Name()), 0777)
			err := CopyFiles(path.Join(srcDir, f.Name()), path.Join(outDir, f.Name()))
			if err != nil {
				return err
			}
		} else {
			data, err := AWSTofuFiles.ReadFile(path.Join(srcDir, f.Name()))
			if err != nil {
				log.Err(err).Msgf("cannot read file %s", f.Name())
				return err
			}
			err = os.WriteFile(path.Join(outDir, f.Name()), data, 0666)
			if err != nil {
				log.Err(err).Msgf("failed to copy file %s to directory %s", f.Name(), outDir)
				return err
			}
		}
	}
	return nil
}
