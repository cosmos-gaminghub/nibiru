# Simple usage with a mounted data directory:
# > docker build -t nbr .
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.nbrd:/root/.nbrd -v ~/.nbrcli:/root/.nbrcli nbr nbrd init
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.nbrd:/root/.nbrd -v ~/.nbrcli:/root/.nbrcli nbr nbrd start
FROM golang:alpine AS build-env

# Set up dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python vim

# Set working directory for the build
WORKDIR /go/src/github.com/cosmos-gaminghub/nibiru

# Add source files
COPY . .

# Install minimum necessary dependencies, build Cosmos SDK, remove packages
RUN apk add --no-cache $PACKAGES && \
    export GO111MODULE=on && \
    make install

# # Final image
# FROM alpine:edge
#
# # Install ca-certificates
# RUN apk add --update ca-certificates
# WORKDIR /root
#
# # Copy over binaries from the build-env
# COPY --from=build-env /go/bin/nbrd /usr/bin/nbrd
# COPY --from=build-env /go/bin/nbrcli /usr/bin/nbrcli
#
# # Run nbrd by default, omit entrypoint to ease using container with nbrcli
# # CMD ["nbrd start"]
