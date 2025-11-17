package types

import "github.com/cloudness-io/cloudness/types/enum"

type TemplateSpec struct {
	Icon     string             `json:"icon,omitempty" yaml:"icon,omitempty"`
	Name     string             `json:"name" yaml:"name"`
	Readme   string             `json:"readme,omitempty" yaml:"readme,omitempty"`
	Tags     []string           `json:"tags,omitempty" yaml:"tags,omitempty"`
	Services []*TemplateService `json:"services" yaml:"services"`
}

type TemplateService struct {
	Icon       string                  `json:"icon,omitempty" yaml:"icon,omitempty"`
	Name       string                  `json:"name" yaml:"name"`
	Build      *BuildConfiguration     `json:"build,omitempty" yaml:"build,omitempty" mapstructure:"build"`
	Deploy     *DeployConfiguration    `json:"deploy,omitempty" yaml:"deploy,omitempty" mapstructure:"deploy"`
	Networking *NetworkConfiguration   `json:"networking,omitempty" yaml:"networking,omitempty" mapstructure:"networking"`
	Volumes    []*TemplateVolumeMounts `json:"volumes,omitempty" yaml:"volumes,omitempty" mapstructure:"volumes"`
	Variables  []*TemplateVariable     `json:"variables,omitempty" yaml:"variables,omitempty" mapstructure:"variables"`
}

type TemplateVariable struct {
	Key   string            `json:"key" yaml:"key"`
	Value string            `json:"value" yaml:"value"`
	Type  enum.VariableType `json:"type" yaml:"type"`
}

type TemplateVolumeMounts struct {
	Name      string `json:"name" yaml:"name" mapstructure:"name"`
	MountPath string `json:"mountPath" yaml:"mountPath" mapstructure:"mountPath"`
}

func (s *TemplateService) ToVolumeMounts() []*VolumeMounts {
	if len(s.Volumes) > 0 {
		dst := make([]*VolumeMounts, 0)
		for _, v := range s.Volumes {
			dst = append(dst, &VolumeMounts{
				VolumeName: v.Name,
				MountPath:  v.MountPath,
			})
		}
		return dst
	}
	return nil
}
