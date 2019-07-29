FROM golang:1.12 as sdk
WORKDIR /app/sdk
COPY sdk/go.mod go.mod
COPY sdk/go.sum go.sum

FROM golang:1.12 as middleware
WORKDIR /app/middleware
COPY middleware/go.mod go.mod
COPY middleware/go.sum go.sum

FROM golang:1.12 as authentication
WORKDIR /app/authentication
COPY authentication/go.mod go.mod
COPY authentication/go.sum go.sum

FROM golang:1.12 as api-gateway
WORKDIR /app/api-gateway
COPY api-gateway/go.mod go.mod
COPY api-gateway/go.sum go.sum

FROM golang:1.12 as static
WORKDIR /app/static
COPY static/go.mod go.mod
COPY static/go.sum go.sum

FROM golang:1.12 as ethereum
COPY --from=sdk /app/sdk /app/sdk
COPY --from=middleware /app/middleware /app/middleware
WORKDIR /app/ethereum
COPY ethereum/go.mod go.mod
COPY ethereum/go.sum go.sum

FROM golang:1.12 as app
WORKDIR /app

COPY --from=api-gateway /app/api-gateway /app/api-gateway
RUN cd api-gateway && go mod download

COPY --from=authentication /app/authentication /app/authentication
RUN cd authentication && go mod download

COPY --from=middleware /app/middleware /app/middleware
RUN cd middleware && go mod download

COPY --from=sdk /app/sdk /app/sdk
RUN cd sdk && go mod download

COPY --from=static /app/static /app/static
RUN cd static && go mod download

COPY --from=ethereum /app/ethereum /app/ethereum
RUN cd ethereum && go mod download
RUN cd api-gateway && go mod download
RUN cd authentication && go mod download
RUN cd sdk && go mod download
RUN cd static && go mod download

COPY . .
