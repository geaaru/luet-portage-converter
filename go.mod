module github.com/Luet-lab/luet-portage-converter

go 1.16

replace github.com/mudler/luet => github.com/geaaru/luet v0.22.1-geaaru

replace github.com/containerd/containerd => github.com/containerd/containerd v1.3.1-0.20200227195959-4d242818bf55

replace github.com/hashicorp/go-immutable-radix => github.com/tonistiigi/go-immutable-radix v0.0.0-20170803185627-826af9ccf0fe

replace github.com/jaguilar/vt100 => github.com/tonistiigi/vt100 v0.0.0-20190402012908-ad4c4a574305

replace github.com/opencontainers/runc => github.com/opencontainers/runc v1.0.0-rc9.0.20200221051241-688cf6d43cc4

replace github.com/docker/docker => github.com/Luet-lab/moby v17.12.0-ce-rc1.0.20200605210607-749178b8f80d+incompatible

require (
	github.com/MottainaiCI/lxd-compose v0.16.1
	github.com/geaaru/pkgs-checker v0.11.0
	github.com/mudler/luet v0.0.0-00010101000000-000000000000
	github.com/onsi/ginkgo/v2 v2.0.0
	github.com/onsi/gomega v1.17.0
	github.com/spf13/cobra v1.2.1
	gopkg.in/yaml.v2 v2.4.0
)
