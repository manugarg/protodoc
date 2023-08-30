# This Dockerfile expects "protodoc" multi-platform binaries to exist in the
# same directory.
#
# Docker image built using this can executed in the following manner:
#   docker run --net host -v $PWD/protodoc.cfg:/etc/protodoc.cfg \
#                         protodoc/protodoc
FROM alpine
COPY protodoc-linux-* ./

ARG TARGETPLATFORM
RUN if [ "$TARGETPLATFORM" = "linux/amd64" ]; then \
  mv protodoc-linux-amd64 protodoc && rm protodoc-linux-*; fi
RUN if [ "$TARGETPLATFORM" = "linux/arm64" ]; then \
  mv protodoc-linux-arm64 protodoc && rm protodoc-linux-*; fi
RUN if [ "$TARGETPLATFORM" = "linux/arm/v7" ]; then \
  mv protodoc-linux-armv7 protodoc && rm protodoc-linux-*; fi

# Metadata params
ARG BUILD_DATE
ARG VERSION
ARG VCS_REF
# Metadata
LABEL org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.name="Protodoc" \
      org.label-schema.vcs-url="https://github.com/manugarg/protodoc" \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.version=$VERSION \
      com.microscaling.license="Apache-2.0"

ENTRYPOINT ["/protodoc"]
