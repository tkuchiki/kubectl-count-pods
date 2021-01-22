module github.com/tkuchiki/kubectl-count-pods

go 1.15

replace (
	k8s.io/api => k8s.io/api v0.0.0-20201209045733-fcac651617f2
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20201209085528-15c5dba13c59
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20201217091618-950728861c24
	k8s.io/client-go => k8s.io/client-go v0.0.0-20201217085940-0964d4be7536
)

require (
	github.com/gizak/termui v3.1.0+incompatible // indirect
	github.com/k0kubun/pp v3.0.1+incompatible
	github.com/olekukonko/tablewriter v0.0.4
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	k8s.io/apimachinery v0.0.0-20201209085528-15c5dba13c59
	k8s.io/cli-runtime v0.0.0-20201217091618-950728861c24
	k8s.io/client-go v0.0.0-20201217085940-0964d4be7536
	k8s.io/klog v1.0.0 // indirect

)
