module github.com/Luet-lab/luet-portage-converter

go 1.16

require (
	github.com/MottainaiCI/lxd-compose v0.13.0
	github.com/MottainaiCI/mottainai-server v0.1.0 // indirect
	github.com/Sabayon/pkgs-checker v0.8.4
	github.com/codegangsta/inject v0.0.0-20150114235600-33e0aa1cb7c0 // indirect
	github.com/cyphar/filepath-securejoin v0.2.3 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/genuinetools/img v0.5.11 // indirect
	github.com/go-logr/logr v1.1.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/hashicorp/go-version v1.3.0 // indirect
	github.com/icza/dyno v0.0.0-20210726202311-f1bafe5d9996 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/jinzhu/copier v0.3.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/lxc/lxd v0.0.0-20210916213034-14c28f0636d4 // indirect
	github.com/magefile/mage v1.11.0 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-sqlite3 v1.14.8 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.4.2 // indirect
	github.com/mudler/luet v0.0.0-20210604142351-a7b4ae67c9b8
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.16.0
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/cobra v1.2.1
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	go.uber.org/zap v1.19.1 // indirect
	golang.org/x/crypto v0.0.0-20210915214749-c084706c2272 // indirect
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f // indirect
	golang.org/x/term v0.0.0-20210916214954-140adaaadfaf // indirect
	gopkg.in/ini.v1 v1.63.2 // indirect
	gopkg.in/yaml.v2 v2.4.0
	helm.sh/helm/v3 v3.7.0 // indirect
	k8s.io/klog/v2 v2.20.0 // indirect
	k8s.io/utils v0.0.0-20210820185131-d34e5cb4466e // indirect
)

replace github.com/containerd/containerd => github.com/containerd/containerd v1.3.1-0.20200227195959-4d242818bf55

replace github.com/docker/docker => github.com/Luet-lab/moby v17.12.0-ce-rc1.0.20200605210607-749178b8f80d+incompatible
