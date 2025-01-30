package client

import (
	"context"
	"path"

	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/render"
	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/store"
	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/tofu"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var secretsConfig = `
mongo:
  - MONGO_USER
  - MONGO_PASSWORD
  - MONGO_REPLICA_SET_KEY
minio:
  - S3_USER
  - S3_PASSWORD
postgres:
  - POSTGRES_PASSWORD
timescaledb:
  - TIMESCALEDB_ADMIN_PASSWORD
  - TIMESCALEDB_PASSWORD
  - TIMESCALEDB_REPLICATION_PASSWORD
`

var externalSecretTemplate = `
{{ .service }}:
  secrets:
    secretManagement:
      externalSecretsOperator:
        - secretStore:
            name: harness-secret-store
            kind: SecretStore
          remote_keys:
          {{- range .secrets }}
            - {{ . }}:
                name: {{ . }}
                property: ""
          {{- end }}
`

var k8sSecrets = []string{"harness-secrets", "minio", "mongodb-replicaset-chart", "postgres"}
var k8sSecretKeyToExternalSecretKey = map[string]string{
	"mongodb-root-password":        "MONGO_PASSWORD",
	"mongodb-replica-set-key":      "MONGO_REPLICA_SET_KEY",
	"mongodbUsername":              "MONGO_USER",
	"root-password":                "S3_PASSWORD",
	"root-user":                    "S3_USER",
	"postgres-password":            "POSTGRES_PASSWORD",
	"timescaledbAdminPassword":     "TIMESCALEDB_ADMIN_PASSWORD",
	"timescaledbPostgresPassword":  "TIMESCALEDB_PASSWORD",
	"PATRONI_REPLICATION_PASSWORD": "TIMESCALEDB_REPLICATION_PASSWORD",
}

type secretsClient struct {
	clientConfig ClientConfig
	configStore  store.DataStore
	outputStore  store.DataStore
}

// Exec implements ResourceClient.
func (s *secretsClient) Exec(ctx context.Context) error {
	if !s.clientConfig.IsManaged {
		return nil
	}
	tofu.ExecuteCommand(ctx, s.clientConfig.ContextDirectory, tofu.InitCommand)
	return tofu.ExecuteCommand(ctx, s.clientConfig.ContextDirectory, tofu.ApplyCommand)
}

// PostExec implements ResourceClient.
func (s *secretsClient) PostExec(ctx context.Context) (map[string]interface{}, error) {
	serviceToSecrets := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(secretsConfig), &serviceToSecrets)
	if err != nil {
		log.Err(err).Msgf("unable to unmarshal secrets config")
		return nil, err
	}
	outputData := make(map[string]interface{})
	for service, secrets := range serviceToSecrets {
		data := make(map[string]interface{})
		data["service"] = service
		data["secrets"] = secrets
		extSecretValue, err := render.RenderString(ctx, data, externalSecretTemplate)
		if err != nil {
			log.Err(err).Msgf("unable to render external secret template")
			return nil, err
		}
		outputData[service] = extSecretValue
	}
	return outputData, nil
}

// PreExec implements ResourceClient.
func (s *secretsClient) PreExec(ctx context.Context) error {
	err := tofu.CopyFiles(path.Join(s.clientConfig.Provider.Name(),
		s.clientConfig.ResourceName),
		s.clientConfig.ContextDirectory)
	if err != nil {
		return err
	}
	secret_keys := make([]string, 0)
	for key := range k8sSecretKeyToExternalSecretKey {
		secret_keys = append(secret_keys, key)
	}
	selfData := make(map[string]interface{})
	selfData["secrets_in_k8s"] = k8sSecrets
	selfData["secret_keys"] = secret_keys

	renderer := render.NewTemplateRenderer(s.configStore, s.outputStore)
	return renderer.Render(ctx, selfData, s.clientConfig.ContextDirectory)
}

func NewSecretsClient(clientConfig ClientConfig, configStore store.DataStore, outputStore store.DataStore) ResourceClient {
	return &secretsClient{
		clientConfig: clientConfig,
		configStore:  configStore,
		outputStore:  outputStore,
	}
}
