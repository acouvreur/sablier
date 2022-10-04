FROM golang:1.18 AS build

ENV PORT 10000

WORKDIR /go/src/sablier

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . /go/src/sablier

ARG TARGETOS
ARG TARGETARCH
RUN make ${TARGETOS}/${TARGETARCH}

FROM alpine

RUN addgroup -S sablier && adduser -S sablier -G sablier
USER sablier:sablier

COPY --from=build --chown=sablier:sablier /go/src/sablier/sablier* /go/bin/sablier

EXPOSE 10000

ENTRYPOINT [ "/go/bin/sablier" ]
CMD [ "start", "--provider.name=docker"]