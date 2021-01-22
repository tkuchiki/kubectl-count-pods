package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd/api"
)

var (
	countPodsExample = `
	# Count the number of pods per status
	kubectl count pods -n NAMESPACE
`

	errNoContext = fmt.Errorf("no context is currently set, use %q to select a new one", "kubectl config use-context <context>")
	version      = "v0.0.1"
)

type Options struct {
	configFlags *genericclioptions.ConfigFlags

	resultingContext     *api.Context
	resultingContextName string

	userSpecifiedCluster   string
	userSpecifiedContext   string
	userSpecifiedAuthInfo  string
	userSpecifiedNamespace string

	rawConfig      api.Config
	listNamespaces bool
	args           []string

	genericclioptions.IOStreams
}

// NewOptions provides an instance of Options with default values
func NewOptions(streams genericclioptions.IOStreams) *Options {
	return &Options{
		configFlags: genericclioptions.NewConfigFlags(true),

		IOStreams: streams,
	}
}

// NewCmdNamespace provides a cobra command wrapping Options
func NewCmd(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewOptions(streams)

	cmd := &cobra.Command{
		Use:          "count pods [flags]",
		Example:      countPodsExample,
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(c, args); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}
			if err := o.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	o.configFlags.AddFlags(cmd.Flags())

	return cmd
}

func (o *Options) Complete(cmd *cobra.Command, args []string) error {
	o.args = args

	var err error
	o.rawConfig, err = o.configFlags.ToRawKubeConfigLoader().RawConfig()
	if err != nil {
		return err
	}

	var namespace string
	namespace, err = cmd.Flags().GetString("namespace")
	if err != nil {
		return err
	}

	if namespace == "" {
		namespace, _, err = o.configFlags.ToRawKubeConfigLoader().Namespace()
		if err != nil {
			return err
		}
	}
	o.userSpecifiedNamespace = namespace

	o.userSpecifiedContext, err = cmd.Flags().GetString("context")
	if err != nil {
		return err
	}

	o.userSpecifiedCluster, err = cmd.Flags().GetString("cluster")
	if err != nil {
		return err
	}

	o.userSpecifiedAuthInfo, err = cmd.Flags().GetString("user")
	if err != nil {
		return err
	}

	currentContext, exists := o.rawConfig.Contexts[o.rawConfig.CurrentContext]
	if !exists {
		return errNoContext
	}

	o.resultingContext = api.NewContext()
	o.resultingContext.Cluster = currentContext.Cluster
	o.resultingContext.AuthInfo = currentContext.AuthInfo

	// if a target context is explicitly provided by the user,
	// use that as our reference for the final, resulting context
	if len(o.userSpecifiedContext) > 0 {
		o.resultingContextName = o.userSpecifiedContext
		if userCtx, exists := o.rawConfig.Contexts[o.userSpecifiedContext]; exists {
			o.resultingContext = userCtx.DeepCopy()
		}
	}

	// override context info with user provided values
	o.resultingContext.Namespace = o.userSpecifiedNamespace

	if len(o.userSpecifiedCluster) > 0 {
		o.resultingContext.Cluster = o.userSpecifiedCluster
	}
	if len(o.userSpecifiedAuthInfo) > 0 {
		o.resultingContext.AuthInfo = o.userSpecifiedAuthInfo
	}

	// generate a unique context name based on its new values if
	// user did not explicitly request a context by name
	if len(o.userSpecifiedContext) == 0 {
		o.resultingContextName = generateContextName(o.resultingContext)
	}

	return nil
}

func generateContextName(fromContext *api.Context) string {
	name := fromContext.Namespace
	if len(fromContext.Cluster) > 0 {
		name = fmt.Sprintf("%s/%s", name, fromContext.Cluster)
	}
	if len(fromContext.AuthInfo) > 0 {
		cleanAuthInfo := strings.Split(fromContext.AuthInfo, "/")[0]
		name = fmt.Sprintf("%s/%s", name, cleanAuthInfo)
	}

	return name
}

func (o *Options) Validate() error {
	if len(o.rawConfig.CurrentContext) == 0 {
		return errNoContext
	}
	if len(o.args) > 1 {
		return fmt.Errorf("either one or no arguments are allowed")
	}

	return nil
}

func (o *Options) Run() error {
	r := resource.
		NewBuilder(o.configFlags).
		Unstructured().
		NamespaceParam(o.userSpecifiedNamespace).
		DefaultNamespace().
		ResourceTypeOrNameArgs(true, "pods").
		Latest().
		Flatten().
		Do()
	if err := r.Err(); err != nil {
		return err
	}

	infos, err := r.Infos()
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Status", "Count"})
	statuses := make(map[string]int)
	objs := make([]*unstructured.Unstructured, len(infos))
	total := 0
	for ix := range infos {
		objs[ix] = infos[ix].Object.(*unstructured.Unstructured)
		o := objs[ix].Object
		status := o["status"].(map[string]interface{})
		phase := status["phase"].(string)
		if phase == "Succeeded" {
			phase = "Completed"
		}
		statuses[status["phase"].(string)]++
		total++
	}

	for status, count := range statuses {
		table.Append([]string{status, fmt.Sprint(count)})
	}
	table.SetFooter([]string{"Total", fmt.Sprint(total)})
	table.Render()

	return nil
}
