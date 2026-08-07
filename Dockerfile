# Copyright © 2026 OpenCHAMI a Series of LF Projects, LLC
# SPDX-FileCopyrightText: Copyright © 2026 OpenCHAMI a Series of LF Projects, LLC
#
# SPDX-License-Identifier: MIT

FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y \
    ca-certificates \
    git \
    bash && \
    apt-get install -y --no-install-recommends curl && \
    rm -rf /var/lib/apt/lists/*

RUN groupadd -g 1000 inventory && \
    useradd -r -u 1000 -g inventory inventory && \
    mkdir -p /data && \
    chown 1000:1000 /data

WORKDIR /home/inventory

COPY inventory-service /usr/local/bin/inventory-service
COPY schemas /schemas

RUN chown -R inventory:inventory /home/inventory
RUN chown -R inventory:inventory /schemas

USER inventory

ENTRYPOINT ["/usr/local/bin/inventory-service"]
CMD ["serve", "--port", "8080", "--database-url", "file:/data/inventory.db?_fk=1", "--auth-enabled=false"]
