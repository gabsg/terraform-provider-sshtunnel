FROM golang:1.17 as builder
ARG BINARY
ARG ARCH
ARG VERSION
WORKDIR /app
COPY ./ /app/
# disable cgo to us libc-musl
RUN go mod tidy && \
    go mod vendor && \
    CGO_ENABLED=0 go build -o ./bin/${BINARY}

FROM hashicorp/terraform:1.0.11
ARG ARCH
ARG BINARY
ARG HOSTNAME
ARG NAMESPACE
ARG NAME
ARG VERSION
# install provider binary from builder
RUN mkdir -p  $HOME/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${ARCH}
COPY --from=builder  /app/bin/${BINARY}  /${BINARY} 
RUN cp /${BINARY} $HOME/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${ARCH}

