package bus

import (
	"context"
	"github.com/nalej/inventory-manager/internal/pkg/server/agent"
	"github.com/nalej/nalej-bus/pkg/queue/inventory/ops"
	"github.com/rs/zerolog/log"
	"time"
)

const InventoryOpsTimeout = time.Second * 30


type InventoryOpsHandler struct {
	agentHandler *agent.Handler
	consumer     *ops.InventoryOpsConsumer
}

func NewInventoryOpsHandler(agentHandler *agent.Handler, consumer *ops.InventoryOpsConsumer) *InventoryOpsHandler {
	return &InventoryOpsHandler{
		agentHandler: agentHandler,
		consumer:     consumer,
	}
}

func (ioh *InventoryOpsHandler) Run() {
	go ioh.consumeAgentOpResponse()
	go ioh.waitRequests()
}

// Endless loop waiting for requests
func (ioh *InventoryOpsHandler) waitRequests() {
	log.Debug().Msg("wait for requests to be received by the inventory ops queue")
	for {
		ctx, cancel := context.WithTimeout(context.Background(), InventoryOpsTimeout)
		// in every iteration this loop consumes data and sends it to the corresponding channels
		currentTime := time.Now()
		err := ioh.consumer.Consume(ctx)
		cancel()
		select {
		case <-ctx.Done():
			// the timeout was reached
			log.Debug().Msgf("no message received since %s", currentTime.Format(time.RFC3339))
		default:
			// we received something or an error
			if err != nil {
				log.Error().Err(err).Msg("error consuming data from inventory events")
			}
		}
	}
}

func (ioh *InventoryOpsHandler) consumeAgentOpResponse() {
	log.Debug().Msg("AgentOpResponse")
	for {
		received := <- ioh.consumer.Config.ChAgentOpResponse
		log.Debug().Msg("agentOpResponse received")
		ioh.agentHandler.CallbackAgentOperation(nil, received)
	}
}

