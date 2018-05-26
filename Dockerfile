FROM golang:alpine as builder
COPY . /go/src/github.com/DennisDenuto/boshver-resource
ENV CGO_ENABLED 0
RUN go build -o /assets/in github.com/DennisDenuto/boshver-resource/in
RUN go build -o /assets/out github.com/DennisDenuto/boshver-resource/out
RUN go build -o /assets/check github.com/DennisDenuto/boshver-resource/check
WORKDIR /go/src/github.com/DennisDenuto/boshver-resource
RUN set -e; for pkg in $(go list ./...); do \
		go test -o "/tests/$(basename $pkg).test" -c $pkg; \
	done

FROM alpine:edge AS resource
RUN apk add --no-cache bash tzdata ca-certificates git jq openssh
RUN git config --global user.email "git@localhost"
RUN git config --global user.name "git"
COPY --from=builder assets/ /opt/resource/
RUN chmod +x /opt/resource/*

FROM resource AS tests
ARG SEMVER_TESTING_ACCESS_KEY_ID
ARG SEMVER_TESTING_SECRET_ACCESS_KEY
ARG SEMVER_TESTING_BUCKET
ARG SEMVER_TESTING_REGION
ARG SEMVER_TESTING_V2_SIGNING
COPY --from=builder /tests /go-tests
WORKDIR /go-tests
RUN set -e; for test in /go-tests/*.test; do \
		$test; \
	done
COPY test/ /opt/resource-tests
RUN /opt/resource-tests/all.sh


FROM resource
