.PHONY: all agent backend flink deploy clean

GO_FLAGS := -ldflags="-s -w" -tags=containers
GO_PGO := -pgo=auto

all: agent backend

# Build Agent with PGO optimization
agent:
	cd agent && go build $(GO_FLAGS) $(GO_PGO) -o bin/monitor-agent ./cmd/agent/

# Build Backend
backend:
	cd backend && go build $(GO_FLAGS) -o bin/monitor-backend ./cmd/server/

# Build Flink job uber-jar
flink:
	cd flink-job && mvn clean package -DskipTests

# Run database migrations
migrate:
	mysql -h 127.0.0.1 -P 4000 -u root < schema/tidb.sql

# Docker images
docker-agent:
	docker build -t pingan/monitor-agent:latest -f agent/Dockerfile agent/

docker-backend:
	docker build -t pingan/monitor-backend:latest -f backend/Dockerfile backend/

docker-flink:
	docker build -t pingan/flink-monitor-job:latest flink-job/

# Local dev environment
dev-up:
	docker compose -f deploy/docker-compose.yml up -d

dev-down:
	docker compose -f deploy/docker-compose.yml down

# Clean artifacts
clean:
	rm -rf agent/bin backend/bin flink-job/target
