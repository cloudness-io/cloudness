package pipeline

import (
	"strings"

	"github.com/cloudness-io/cloudness/app/services/config"
	"github.com/cloudness-io/cloudness/types"
)

type RestartPolicy string

const (
	RestartPolicyNever  RestartPolicy = "Never"
	RestartPolicyAlways RestartPolicy = "Always"
)

type (
	// RunnerContextInput contains the input params to generate runner context.
	RunnerContextInput struct {
		Application        *types.Application
		Variables          map[string]*types.Variable
		Volumes            []*types.Volume
		Deployment         *types.Deployment
		PreviousDeployment *types.Deployment
		Netrc              *types.Netrc
		Config             *config.PipelineConfig
		ServerRestctions   *types.ServerRestrictions
	}

	// RunnerContext holds the information required for runner to execute a deployment.
	RunnerContext struct {
		RunnerName     string
		Deployment     *types.Deployment
		ApplicationUID int64
		Domain         string
		Workspaces     []*Volume
		Steps          []*Step
		InitSteps      []*Step
		Secrets        []*Secret
		Variables      []*Variable
		ConfigFiles    []*ConfigFile // Files to add to ConfigMap
		ResourcesLimit *ResourcesLimit
	}

	Step struct {
		Name             string
		Image            string
		Command          []string
		Args             []string
		ScriptCommands   []string
		WorkingDir       string
		Envs             map[string]string
		VolumeMounts     []*VolumeMount
		Secrets          []*SecretEnv
		Variables        []*VariableEnv
		ConfigFileMounts []*ConfigFileMount // ConfigMap file mounts
		Privileged       bool
		RestartPolicy    RestartPolicy

		//housekeeping
		Liveness *Liveness
	}

	Secret struct {
		Name string
		Data string
		Mask bool
	}

	Variable struct {
		Name  string
		Value string
	}

	SecretEnv struct {
		Key string
	}

	VariableEnv struct {
		Key string
	}

	// ConfigFile represents a file to be mounted from ConfigMap
	ConfigFile struct {
		Key      string // Key in ConfigMap
		Filename string // Filename when mounted
		Content  string // File content
	}

	// ConfigFileMount represents a mounted ConfigMap as files
	ConfigFileMount struct {
		Path string   // Mount path directory
		Keys []string // Keys to mount as files
	}

	Volume struct {
		ID   string
		Size int64
	}

	VolumeMount struct {
		ID       string
		Path     string
		Readonly bool
	}

	//Liveness
	Liveness struct {
		Exec *LivenessExec
	}

	LivenessExec struct {
		Command []string
	}

	ResourcesLimit struct {
		Memory float64
		CPU    float64
	}
)

func (s *Step) AppendArgs(args []string) {
	if len(s.Args) == 0 {
	}
	s.Args = append(s.Args, args...)
}

func (s *Step) AddScriptCmd(cmd string) {
	s.ScriptCommands = append(s.ScriptCommands, cmd)
}

func (s *Step) AddStripCmds(cmds ...string) {
	s.ScriptCommands = append(s.ScriptCommands, cmds...)
}

func (s *Step) GenerateShellScript() string {
	var sb strings.Builder
	for _, cmd := range s.ScriptCommands {
		sb.WriteString(cmd + "\n")
	}
	return sb.String()
}
