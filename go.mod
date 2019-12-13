module github.com/ameydev/groot

go 1.13

replace github.com/infracloudio/ksearch => github.com/arush-sal/ksearch v0.0.0-20191204064629-07f4be858800

require (
	github.com/ameydev/ksearch v0.0.0-20191205152724-081db6afc2b0
	github.com/infracloudio/ksearch v0.0.0-00010101000000-000000000000
	github.com/mitchellh/go-homedir v1.1.0
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.5.0
	k8s.io/api v0.0.0-20191121015604-11707872ac1c
	k8s.io/apimachinery v0.0.0-20191203211716-adc6f4cd9e7d
	k8s.io/client-go v0.0.0-20190620085101-78d2af792bab
)
