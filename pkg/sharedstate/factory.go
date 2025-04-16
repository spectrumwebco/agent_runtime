package sharedstate

import (
	"github.com/spectrumwebco/agent_runtime/internal/eventstream"
	"github.com/spectrumwebco/agent_runtime/pkg/ai/vercel"
	"github.com/spectrumwebco/agent_runtime/pkg/sharedstate/adapters"
)

type Factory struct {
	sharedStateManager *SharedStateManager
	eventStream        *eventstream.Stream
}

func NewFactory(sharedStateManager *SharedStateManager, eventStream *eventstream.Stream) *Factory {
	return &Factory{
		sharedStateManager: sharedStateManager,
		eventStream:        eventStream,
	}
}

func (f *Factory) CreateRSCAdapter(rscClient *vercel.RSCAdapter) *adapters.RSCAdapter {
	return adapters.NewRSCAdapter(rscClient)
}

func (f *Factory) CreateRSCIntegration(
	rscAdapter *adapters.RSCAdapter,
	rscClient *vercel.RSCAdapter,
	rscIntegration *vercel.RSCIntegration,
) *RSCIntegration {
	return NewRSCIntegration(
		f.sharedStateManager,
		rscAdapter,
		f.eventStream,
		rscClient,
		rscIntegration,
	)
}
