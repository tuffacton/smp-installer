package client

import (
	"context"
	"path"
	"strconv"
	"time"

	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/render"
	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/store"
	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/tofu"
	"github.com/rs/zerolog/log"
)

type waitForClusterClient struct {
	clientConfig ClientConfig
	configStore  store.DataStore
	outputStore  store.DataStore
}

// Exec implements ResourceClient.
func (w *waitForClusterClient) Exec(ctx context.Context) error {
	isManaged := w.configStore.GetBool(ctx, "kubernetes.manage")
	if !isManaged {
		return nil
	}
	clusterStatus := w.outputStore.GetString(ctx, "kubernetes.cluster_status")
	numberOfActiveNodes := w.outputStore.GetInt(ctx, "kubernetes.instance_count")
	for clusterStatus != "ACTIVE" || numberOfActiveNodes == 0 {
		log.Info().Msgf("waiting for cluster to be ACTIVE")
		time.Sleep(30 * time.Second)
		tofu.ExecuteCommand(ctx, w.clientConfig.ContextDirectory, tofu.InitCommand)
		err := tofu.ExecuteCommand(ctx, w.clientConfig.ContextDirectory, tofu.ApplyCommand)
		if err != nil {
			log.Err(err).Msgf("unable to find cluster status")
		}
		clusterStatus, err = tofu.GetOutput(ctx, w.clientConfig.ContextDirectory, "cluster_status", ".")
		if err != nil {
			log.Err(err).Msgf("unable to retrieve cluster status from tofu output")
			return err
		}
		instanceCount, err := tofu.GetOutput(ctx, w.clientConfig.ContextDirectory, "instance_count", ".")
		if err != nil {
			log.Err(err).Msgf("unable to retrieve instance count from tofu output")
			return err
		}
		numberOfActiveNodes, err = strconv.Atoi(instanceCount)
		if err != nil {
			log.Err(err).Msgf("unable to convert instance count to int")
			return err
		}
	}
	return nil
}

// PostExec implements ResourceClient.
func (w *waitForClusterClient) PostExec(ctx context.Context) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

// PreExec implements ResourceClient.
func (w *waitForClusterClient) PreExec(ctx context.Context) error {
	err := tofu.CopyFiles(path.Join(w.clientConfig.Provider.Name(), w.clientConfig.ResourceName), w.clientConfig.ContextDirectory)
	if err != nil {
		return err
	}
	renderer := render.NewTemplateRenderer(w.configStore, w.outputStore)
	return renderer.Render(ctx, map[string]interface{}{}, w.clientConfig.ContextDirectory)
}

func NewWaitForClusterClient(clientConfig ClientConfig, configStore store.DataStore, outputStore store.DataStore) ResourceClient {
	return &waitForClusterClient{
		clientConfig: clientConfig,
		configStore:  configStore,
		outputStore:  outputStore,
	}
}
