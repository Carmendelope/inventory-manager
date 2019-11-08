/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package bus

import (
	"context"
	"github.com/nalej/inventory-manager/internal/pkg/server/agent"
	"github.com/nalej/inventory-manager/internal/pkg/server/edgecontroller"
	"github.com/nalej/nalej-bus/pkg/queue/inventory/events"
	"github.com/rs/zerolog/log"
	"time"
)

const InventoryEventsTimeout = time.Second * 30

// TODO Refactor this package to be outside of server, and move service.go to other package.

type InventoryEventsHandler struct {
	ecHandler    *edgecontroller.Handler
	agentHandler *agent.Handler
	consumer     *events.InventoryEventsConsumer
}

func NewInventoryEventsHandler(ecHandler *edgecontroller.Handler, agentHandler *agent.Handler, consumer *events.InventoryEventsConsumer) *InventoryEventsHandler {
	return &InventoryEventsHandler{
		ecHandler:    ecHandler,
		agentHandler: agentHandler,
		consumer:     consumer,
	}
}

func (ieh *InventoryEventsHandler) Run() {
	go ieh.consumeEICStart()
	go ieh.consumeEdgeControllerId()
	go ieh.consumeAgentAlive()
	go ieh.consumeAgentUninstalled()
	go ieh.waitRequests()
}

// Endless loop waiting for requests
func (ieh *InventoryEventsHandler) waitRequests() {
	log.Debug().Msg("wait for requests to be received by the inventory events queue")
	for {
		ctx, cancel := context.WithTimeout(context.Background(), InventoryEventsTimeout)
		// in every iteration this loop consumes data and sends it to the corresponding channels
		currentTime := time.Now()
		err := ieh.consumer.Consume(ctx)
		cancel()
		select {
		case <-ctx.Done():
			// the timeout was reached
			log.Debug().Msgf("no message received since %s", currentTime.Format(time.RFC3339))
		default:
			// we received something or an error
			if err != nil {
				log.Error().Err(err).Msg("error consuming data from application ops")
			}
		}
	}
}

func (ieh *InventoryEventsHandler) consumeEICStart() {
	log.Debug().Msg("consuming EICStart")
	for {
		received := <-ieh.consumer.Config.ChEICStart
		log.Debug().Interface("message", received).Msg("EICSTart received")
		ieh.ecHandler.EICStart(nil, received)
	}
}

func (ieh *InventoryEventsHandler) consumeEdgeControllerId() {
	log.Debug().Msg("consuming EdgeControllerId")
	for {
		received := <-ieh.consumer.Config.ChEdgeControllerId
		log.Debug().Interface("message", received).Msg("EdgeControllerId received")
		ieh.ecHandler.EICAlive(nil, received)
	}
}

func (ieh *InventoryEventsHandler) consumeAgentAlive() {
	log.Debug().Msg("consuming AgentAlive")
	for {
		received := <-ieh.consumer.Config.ChAgentsAlive
		log.Debug().Interface("message", received).Msg("AgentAlive received")
		ieh.agentHandler.LogAgentAlive(nil, received)
	}
}
func (ieh *InventoryEventsHandler) consumeAgentUninstalled() {
	log.Debug().Msg("consuming AgentUninstalled")
	for {
		received := <-ieh.consumer.Config.ChUninstalledAssetId
		log.Debug().Interface("message", received).Msg("AgentUninstalled received")
		ieh.agentHandler.UninstalledAgent(nil, received)
	}
}