package provision

import (
	"context"

	"github.com/sapcc/ironic_temper/pkg/clients"
	"github.com/sapcc/ironic_temper/pkg/model"
	log "github.com/sirupsen/logrus"
)

type SchedulerError struct {
	Err  string
	Node *model.IronicNode
}

func (n *SchedulerError) Error() string {
	return n.Err
}

type ErrorHandler struct {
	Errors  chan error
	ctx     context.Context
	clients *clients.Client
}

func NewErrorHandler(ctx context.Context, c *clients.Client) (e ErrorHandler) {
	errors := make(chan error, 0)
	e.Errors = errors
	e.ctx = ctx
	e.clients = c
	go e.initHandler()
	return e
}

func (e ErrorHandler) initHandler() {
	go func() {
		select {
		case err := <-e.Errors:
			if serr, ok := err.(*SchedulerError); ok {
				log.Infof("error tempering node %s. err: %s", serr.Node.UUID, serr.Err)
				if serr.Node.InstanceUUID != "" {
					e.clients.DeleteNodeTestDeployment(serr.Node)
				}
			} else {
				log.Error(err.Error())
			}
		case <-e.ctx.Done():
			return
		}
	}()
}
