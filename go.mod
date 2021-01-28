module github.com/Luet-lab/luet-portage-converter

go 1.15

require (
	github.com/MottainaiCI/lxd-compose v0.7.2
	github.com/Sabayon/pkgs-checker v0.8.1-0.20210128171435-eaf496434915
	github.com/mudler/luet v0.0.0-20210125133601-4eab1eb738c5
	github.com/onsi/ginkgo v1.14.2
	github.com/onsi/gomega v1.10.3
	github.com/spf13/cobra v1.1.1
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/containerd/containerd => github.com/containerd/containerd v1.3.1-0.20200227195959-4d242818bf55

replace github.com/docker/docker => github.com/Luet-lab/moby v17.12.0-ce-rc1.0.20200605210607-749178b8f80d+incompatible
