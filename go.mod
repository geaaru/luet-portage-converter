module github.com/Luet-lab/luet-portage-converter

go 1.15

require (
	github.com/MottainaiCI/lxd-compose v0.7.2
	github.com/Sabayon/pkgs-checker v0.8.1
	github.com/codegangsta/inject v0.0.0-20150114235600-33e0aa1cb7c0 // indirect
	github.com/genuinetools/img v0.5.11 // indirect
	github.com/magefile/mage v1.11.0 // indirect
	github.com/magiconair/properties v1.8.4 // indirect
	github.com/mattn/go-sqlite3 v1.14.6 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/mudler/luet v0.0.0-20210604142351-a7b4ae67c9b8
	github.com/onsi/ginkgo v1.14.2
	github.com/onsi/gomega v1.10.3
	github.com/pelletier/go-toml v1.8.1 // indirect
	github.com/sirupsen/logrus v1.8.0 // indirect
	github.com/spf13/afero v1.5.1 // indirect
	github.com/spf13/cobra v1.1.3
	golang.org/x/sys v0.0.0-20210218155724-8ebf48af031b // indirect
	golang.org/x/text v0.3.5 // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/containerd/containerd => github.com/containerd/containerd v1.3.1-0.20200227195959-4d242818bf55

replace github.com/docker/docker => github.com/Luet-lab/moby v17.12.0-ce-rc1.0.20200605210607-749178b8f80d+incompatible
