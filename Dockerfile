FROM golang:1.18 AS build

WORKDIR /myoperator

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . /myoperator
RUN make install-tools
RUN make build

FROM scratch AS final
COPY --from=build /myoperator/bin/otemplate /otemplate
ENTRYPOINT ["/otemplate"]
