package templates

import (
	"bytes"
	"embed"
	"text/template"
)

//go:embed static/*
var templateFs embed.FS

var (
	statelessTmpl    *template.Template
	statelessScripts *template.Template
	statefulTmpl     *template.Template
	statefulScripts  *template.Template
)

// latest
var (
	kubeCommon    *template.Template
	kubePVC       *template.Template
	kubeApp       *template.Template
	kubeHttproute *template.Template
	kubeScripts   string
)

func init() {
	statelessTmpl = getTemplate("stateless.yaml")
	statelessScripts = getTemplate("stateless.sh")
	statefulTmpl = getTemplate("stateful.yaml")
	statefulScripts = getTemplate("stateful.sh")

	kubeCommon = getTemplate("1-common.yaml")
	kubePVC = getTemplate("2-pvc.yaml")
	kubeApp = getTemplate("3-app.yaml")
	kubeHttproute = getTemplate("4-httproute.yaml")
}

func getFileContent(fileName string) string {
	file, err := templateFs.ReadFile("static/" + fileName)
	if err != nil {
		panic(err)
	}
	return string(file)
}

func getTemplate(fileName string) *template.Template {
	template, err := template.ParseFS(templateFs, "static/"+fileName)
	if err != nil {
		panic(err)
	}

	return template
}

type (
	TemplateIn struct {
		Identifier    string
		Namespace     string
		HasState      bool
		Image         string
		Command       []string
		Args          []string
		ServicePorts  []int
		ServiceDomain *ServiceDomain
		Volumes       []*Volume
		MaxReplicas   int64
		CPU           int64
		Memory        float64
		Variables     map[string]string
		Secrets       map[string]string
	}

	ServiceDomain struct {
		Websecure bool
		Domain    string
		Port      int
	}

	Volume struct {
		VolumeName string
		Storage    string
		MountPath  string
	}

	DeploymentIn struct {
		Identifier   string
		Namespace    string
		ManifestPath string
		HasHTTPRoute bool
	}
)

func GenerateKubeTemplates(input *TemplateIn) (common string, pvc string, app string, route string, err error) {
	common, err = renderTemplate(kubeCommon, input)
	if err != nil {
		return
	}

	pvc, err = renderTemplate(kubePVC, input)
	if err != nil {
		return
	}

	app, err = renderTemplate(kubeApp, input)
	if err != nil {
		return
	}

	route, err = renderTemplate(kubeHttproute, input)
	return
}

func renderTemplate(tmpl *template.Template, input *TemplateIn) (string, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, input); err != nil {
		return "", err
	}

	return formatTemplate(buf), nil
}

func GenerateKubeStatelessTemplate(input *TemplateIn) (string, error) {
	var buf bytes.Buffer

	if err := statelessTmpl.Execute(&buf, input); err != nil {
		return "", err
	}

	return formatTemplate(buf), nil
}

func GenerateKubeStatefulTemplate(input *TemplateIn) (string, error) {
	var buf bytes.Buffer

	if err := statefulTmpl.Execute(&buf, input); err != nil {
		return "", err
	}

	return formatTemplate(buf), nil
}

func formatTemplate(buf bytes.Buffer) string {
	return buf.String()
}

func GenerateKubeStatelessScript(input *DeploymentIn) (string, error) {
	var buf bytes.Buffer

	if err := statelessScripts.Execute(&buf, input); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func GenerateKubeStatefulScript(input *DeploymentIn) (string, error) {
	var buf bytes.Buffer

	if err := statefulScripts.Execute(&buf, input); err != nil {
		return "", err
	}
	return buf.String(), nil
}
