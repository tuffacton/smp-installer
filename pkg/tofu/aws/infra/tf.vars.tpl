region = {{ .Region | quote }}
cluster_version = {{ .KubernetesConfig.Version | quote }}
machine_type = {{ .NodesConfig.MachineType | quote }}
minimum_nodes = {{ .NodesConfig.MinimumNodes }}
maximum_nodes = {{ .NodesConfig.MaximumNodes }}
initial_nodes = {{ .NodesConfig.InitialNodes }}