ARG BUILDPLATFORM
FROM --platform=${BUILDPLATFORM:-linux/amd64} busybox AS builder

WORKDIR /go/src

ARG COMMIT_SHA=""

ARG PROJECT_NAME=""

ARG TARGETPLATFORM
RUN echo ${TARGETPLATFORM} > a.txt && touch b.txt

FROM busybox AS builder2

WORKDIR /go/src

RUN touch b.txt

FROM busybox

WORKDIR /todo

COPY --from=builder2 /go/src/b.txt ./

COPY --from=builder /go/src/a.txt ./

