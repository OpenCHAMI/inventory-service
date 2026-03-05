# Copyright © 2025 OpenCHAMI a Series of LF Projects, LLC
#
# SPDX-License-Identifier: MIT

FROM debian:bookworm-slim AS builder


RUN mkdir -p /data && chown 1000:1000 /data

# FROM rockylinux:9.3.20231119 
# FROM docker.io/library/alpine:3.15
# FROM docker.io/library/alpine:3.23



FROM gcr.io/distroless/static-debian12:nonroot

ENV SMD_PORT=8080

WORKDIR /app

COPY --from=builder /data /data

COPY bin/smd2-server /usr/local/bin/smd2-server

USER nonroot:nonroot

ENTRYPOINT ["/usr/local/bin/smd2-server"]
CMD ["serve", "--database-url", "file:/data/smd2.db?_fk=1"]
