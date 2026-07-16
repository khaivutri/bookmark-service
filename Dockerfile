FROM golang:1.26-alpine AS base

RUN mkdir -p /opt/app

WORKDIR /opt/app

COPY . .

RUN go mod download 

FROM base AS build

RUN apk add build-base

RUN GOOS=linux go build -tags musl -ldflags "-w -s" -o bookmark-service cmd/api/main.go

FROM base AS test-exec

ARG _outputdir="/tmp/coverage"
ARG COVERAGE_EXCLUDE
ARG COVERAGE_THRESHOLD

RUN mkdir -p "${_outputdir}" &&  \
    go test ./... -coverprofile=coverage.tmp -coverpkg=./... -covermode=atomic -p 1 && \
    grep -vE "${COVERAGE_EXCLUDE}" coverage.tmp > ${_outputdir}/coverage.out && \
    go tool cover -html=${_outputdir}/coverage.out -o ${_outputdir}/coverage.html && \
    total=$(go tool cover -func=${_outputdir}/coverage.out | grep total: | awk '{print $3}' | sed 's/%//') && \
    echo "Current total coverage: $total%" && \
    if [ $(echo "$total < ${COVERAGE_THRESHOLD}" | bc -l) -eq 1 ]; then \
        echo "❌ Coverage ($total%) is below threshold (${COVERAGE_THRESHOLD}%)"; \
        exit 1; \
    else \
        echo "✅ Coverage ($total%) meets threshold (${COVERAGE_THRESHOLD}%)"; \
    fi

FROM scratch AS test 

ARG _outputdir="/tmp/coverage"

COPY --from=test-exec ${_outputdir}/coverage.out /

COPY --from=test-exec ${_outputdir}/coverage.html /


FROM alpine AS final

ARG app_name=app
ENV TZ=Asia/Ho_Chi_Minh

WORKDIR /app

COPY --from=build /opt/app/bookmark-service /app/bookmark-service
COPY --from=build /opt/app/docs /app/docs

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

CMD ["/app/bookmark-service"]