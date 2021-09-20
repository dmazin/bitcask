FROM golang:1.16 as build

# Create appuser.
# See https://stackoverflow.com/a/55757473/12429735
ENV USER=appuser
ENV UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

RUN apt-get update && apt-get install -y ca-certificates
RUN ls
# RUN go get github.com/dmazin/naivedb

# Build
WORKDIR /go/src/github.com/rakyll/hey
# RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux make build

###############################################################################
# final stage
FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
USER appuser:appuser

ARG APPLICATION="naivedb"
ARG DESCRIPTION="Re-implementation of Bitcask"
ARG PACKAGE="dmazin/naivedb"

COPY --from=build /go/bin/${APPLICATION} /naivedb-server
ENTRYPOINT ["/naivedb-server"]