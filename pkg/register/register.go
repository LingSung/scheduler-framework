package register

import (
	"github.com/LingSung/scheduler-framework/pkg/xtutx"
	"github.com/spf13/cobra"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"
)

func Register() *cobra.Command {
	return app.NewSchedulerCommand(
		app.WithPlugin(xtutx.Name, xtutx.New),
	)
}