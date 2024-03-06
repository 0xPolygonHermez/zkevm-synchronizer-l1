
DOCKER_COMPOSE := docker-compose -f docker-compose.yml
DOCKER_COMPOSE_STATE_DB := zkevm-state-db
DOCKER_COMPOSE_POOL_DB := zkevm-pool-db
DOCKER_COMPOSE_RPC_DB := zkevm-rpc-db
DOCKER_COMPOSE_BRIDGE_DB := zkevm-bridge-db
DOCKER_COMPOSE_ZKEVM_NODE := zkevm-synchronizer-l1
DOCKER_COMPOSE_L1_NETWORK := zkevm-mock-l1-network
DOCKER_COMPOSE_ZKPROVER := zkevm-prover
DOCKER_COMPOSE_BRIDGE := zkevm-bridge-service
DOCKER_COMPOSE_SYNC_DB := zkevm-sync-db

RUN_STATE_DB := $(DOCKER_COMPOSE) up -d $(DOCKER_COMPOSE_STATE_DB)
RUN_POOL_DB := $(DOCKER_COMPOSE) up -d $(DOCKER_COMPOSE_POOL_DB)
RUN_BRIDGE_DB := $(DOCKER_COMPOSE) up -d $(DOCKER_COMPOSE_BRIDGE_DB)
RUN_DBS := ${RUN_BRIDGE_DB} && ${RUN_STATE_DB} && ${RUN_POOL_DB}
RUN_NODE := $(DOCKER_COMPOSE) up -d $(DOCKER_COMPOSE_ZKEVM_NODE)
RUN_NODE_V1TOV2 := $(DOCKER_COMPOSE) up -d $(DOCKER_COMPOSE_ZKEVM_NODE_V1TOV2)
RUN_AGGREGATOR_V1TOV2 := $(DOCKER_COMPOSE) up -d $(DOCKER_COMPOSE_ZKEVM_AGGREGATOR_V1TOV2)
RUN_L1_NETWORK := $(DOCKER_COMPOSE) up -d $(DOCKER_COMPOSE_L1_NETWORK)
RUN_L1_NETWORK_V1TOV2 := $(DOCKER_COMPOSE) up -d $(DOCKER_COMPOSE_L1_NETWORK_V1TOV2)
RUN_ZKPROVER := $(DOCKER_COMPOSE) up -d $(DOCKER_COMPOSE_ZKPROVER)
RUN_ZKPROVER_V1TOV2 := $(DOCKER_COMPOSE) up -d $(DOCKER_COMPOSE_ZKPROVER_V1TOV2)
RUN_BRIDGE := $(DOCKER_COMPOSE) up -d $(DOCKER_COMPOSE_BRIDGE)
RUN_BRIDGE_V1TOV2 := $(DOCKER_COMPOSE) up -d $(DOCKER_COMPOSE_BRIDGE_V1TOV2)

STOP_NODE_DB := $(DOCKER_COMPOSE) stop $(DOCKER_COMPOSE_NODE_DB) && $(DOCKER_COMPOSE) rm -f $(DOCKER_COMPOSE_NODE_DB)
STOP_BRIDGE_DB := $(DOCKER_COMPOSE) stop $(DOCKER_COMPOSE_BRIDGE_DB) && $(DOCKER_COMPOSE) rm -f $(DOCKER_COMPOSE_BRIDGE_DB)
STOP_DBS := ${STOP_NODE_DB} && ${STOP_BRIDGE_DB}
STOP_NODE := $(DOCKER_COMPOSE) stop $(DOCKER_COMPOSE_ZKEVM_NODE) && $(DOCKER_COMPOSE) rm -f $(DOCKER_COMPOSE_ZKEVM_NODE)
STOP_NODE_V1TOV2 := $(DOCKER_COMPOSE) stop $(DOCKER_COMPOSE_ZKEVM_NODE_V1TOV2) && $(DOCKER_COMPOSE) rm -f $(DOCKER_COMPOSE_ZKEVM_NODE_V1TOV2)
STOP_AGGREGATOR_V1TOV2 := $(DOCKER_COMPOSE) stop $(DOCKER_COMPOSE_ZKEVM_AGGREGATOR_V1TOV2) && $(DOCKER_COMPOSE) rm -f $(DOCKER_COMPOSE_ZKEVM_AGGREGATOR_V1TOV2)
STOP_NETWORK := $(DOCKER_COMPOSE) stop $(DOCKER_COMPOSE_L1_NETWORK) && $(DOCKER_COMPOSE) rm -f $(DOCKER_COMPOSE_L1_NETWORK)
STOP_NETWORK_V1TOV2 := $(DOCKER_COMPOSE) stop $(DOCKER_COMPOSE_L1_NETWORK_V1TOV2) && $(DOCKER_COMPOSE) rm -f $(DOCKER_COMPOSE_L1_NETWORK_V1TOV2)
STOP_ZKPROVER := $(DOCKER_COMPOSE) stop $(DOCKER_COMPOSE_ZKPROVER) && $(DOCKER_COMPOSE) rm -f $(DOCKER_COMPOSE_ZKPROVER)
STOP_ZKPROVER_V1TOV2 := $(DOCKER_COMPOSE) stop $(DOCKER_COMPOSE_ZKPROVER_V1TOV2) && $(DOCKER_COMPOSE) rm -f $(DOCKER_COMPOSE_ZKPROVER_V1TOV2)
STOP_BRIDGE := $(DOCKER_COMPOSE) stop $(DOCKER_COMPOSE_BRIDGE) && $(DOCKER_COMPOSE) rm -f $(DOCKER_COMPOSE_BRIDGE)
STOP_BRIDGE_V1TOV2 := $(DOCKER_COMPOSE) stop $(DOCKER_COMPOSE_BRIDGE_V1TOV2) && $(DOCKER_COMPOSE) rm -f $(DOCKER_COMPOSE_BRIDGE_V1TOV2)
STOP := $(DOCKER_COMPOSE) down --remove-orphans



GO_BASE := $(shell pwd)
GO_BIN := $(GO_BASE)/dist
GO_ENV_VARS := GO_BIN=$(GO_BIN)
GO_BINARY := zkevm-bridge
GO_CMD := $(GO_BASE)/cmd


.PHONY: run-db-sync
run-db-sync: ## Runs the synchronizer database
	$(DOCKER_COMPOSE) up -d $(DOCKER_COMPOSE_SYNC_DB)

.PHONY: stop-db-sync
stop-db-sync: ## Stops the node database
	$(DOCKER_COMPOSE) stop $(DOCKER_COMPOSE_SYNC_DB) && $(DOCKER_COMPOSE) rm -f $(DOCKER_COMPOSE_SYNC_DB)

.PHONY: run-dbs
run-dbs: run-db-sync ## Runs the node database
	

.PHONY: stop-dbs
stop-dbs: stop-db-sync ## Stops the node database
	

.PHONY: run-node
run-node: ## Runs the node
	$(RUN_NODE)

.PHONY: stop-node
stop-node: ## Stops the node
	$(STOP_NODE)

.PHONY: run-network
run-network: ## Runs the l1 network
	$(RUN_L1_NETWORK)

.PHONY: stop-network
stop-network: ## Stops the l1 network
	$(STOP_NETWORK)



.PHONY: run-prover
run-prover: ## Runs the zk prover
	$(RUN_ZKPROVER)

.PHONY: stop-prover
stop-prover: ## Stops the zk prover
	$(STOP_ZKPROVER)



.PHONY: stop
stop: ## Stops all services
	$(STOP)

.PHONY: restart
restart: stop run ## Executes `make stop` and `make run` commands

.PHONY: run
run: stop run-dbs ## runs all services
	$(RUN_DBS)
	$(RUN_L1_NETWORK)
	sleep 5
	$(RUN_ZKPROVER)
	sleep 3
	$(RUN_NODE)
	sleep 7
	$(RUN_BRIDGE)




## Help display.
## Pulls comments from beside commands and prints a nicely formatted
## display with the commands and their usage information.
.DEFAULT_GOAL := help

.PHONY: help
help: ## Prints this help
		@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

