package eventmesh

import (
	"context"
	"fmt"
	"go.uber.org/zap"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"

	eventingclient "knative.dev/eventing/pkg/client/clientset/versioned"
)

// Endpoint is the HTTP handler that's used to serve the event mesh data.
type Endpoint struct {
	inClusterConfig *rest.Config
	logger          *zap.SugaredLogger
}

// ensure that Endpoint implements the StrictServerInterface
var _ StrictServerInterface = &Endpoint{}

func NewEndpoint(inClusterConfig *rest.Config, logger *zap.SugaredLogger) *Endpoint {
	return &Endpoint{
		inClusterConfig: inClusterConfig,
		logger:          logger,
	}
}

func (e Endpoint) GetEventMesh(ctx context.Context, _ GetEventMeshRequestObject) (GetEventMeshResponseObject, error) {
	logger := e.logger

	authToken, ok := GetAuthToken(ctx)
	if !ok {
		return GetEventMesh401JSONResponse{
			Error: "Authorization header is missing",
		}, nil
	}

	// config := rest.CopyConfig(e.inClusterConfig)
	// config.BearerToken = authToken

	config := rest.AnonymousClientConfig(e.inClusterConfig)
	config.BearerToken = authToken

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		logger.Errorw("Error creating dynamic client", "error", err)
		return nil, fmt.Errorf("error creating dynamic client: %w", err)
	}

	eventingClient, err := eventingclient.NewForConfig(config)
	if err != nil {
		logger.Errorw("Error creating eventing client", "error", err)
		return nil, fmt.Errorf("error creating eventing client: %w", err)
	}

	eventMesh, err := BuildEventMesh(ctx, eventingClient, dynamicClient, logger)
	if err != nil {
		logger.Errorw("Error building event mesh", "error", err)
		return nil, fmt.Errorf("error building event mesh: %w", err)
	}

	return GetEventMesh200JSONResponse{
		Brokers:    eventMesh.Brokers,
		EventTypes: eventMesh.EventTypes,
	}, nil
}
