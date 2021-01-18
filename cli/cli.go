package cli

import (
	"github.com/hackaio/pk"
	"github.com/spf13/cobra"
)

type wrapper struct {
	service pk.Service
}

func (w wrapper) start(cmd *cobra.Command, args []string) {
	panic("implement me")
}

func (w wrapper) add(cmd *cobra.Command, args []string) {
	panic("implement me")
}

func (w wrapper) get(cmd *cobra.Command, args []string) {
	panic("implement me")
}

func (w wrapper) list(cmd *cobra.Command, args []string) {
	panic("implement me")
}

func (w wrapper) delete(cmd *cobra.Command, args []string) {
	panic("implement me")
}

func (w wrapper) update(cmd *cobra.Command, args []string) {
	panic("implement me")
}

func NewWrapper(service pk.Service) CLI {
	return wrapper{service: service}
}

type CLI interface {
	start(cmd *cobra.Command, args []string)
	add(cmd *cobra.Command, args []string)
	get(cmd *cobra.Command, args []string)
	list(cmd *cobra.Command, args []string)
	delete(cmd *cobra.Command, args []string)
	update(cmd *cobra.Command, args []string)
}
