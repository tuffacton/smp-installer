package util

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

func MergeFiles(filenames ...string) (map[string]interface{}, error) {
	if len(filenames) < 1 {
		return nil, fmt.Errorf("number of arguments should be more than 0")
	}
	var master = make(map[string]interface{})
	for _, f := range filenames {
		bs, err := os.ReadFile(f)
		var override = make(map[string]interface{})
		if err != nil {
			log.Err(err).Msgf("could not read file %s", f)
			return nil, err
		}
		if err := yaml.Unmarshal(bs, &override); err != nil {
			log.Err(err).Msgf("could not read yaml %s", f)
			return nil, err
		}
		master = MergeMaps(master, override)
	}
	return master, nil
}

func MergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(map[string]interface{}); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]interface{}); ok {
					out[k] = MergeMaps(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}
