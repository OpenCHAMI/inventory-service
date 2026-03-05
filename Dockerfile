# Copyright © 2026 OpenCHAMI a Series of LF Projects, LLC
# SPDX-FileCopyrightText: Copyright © 2026 OpenCHAMI a Series of LF Projects, LLC
#
# SPDX-License-Identifier: MIT

FROM debian:bookworm-slim

ENV SMD_PORT=8080

RUN apt-get update && apt-get install -y \
    ca-certificates \
    git \
    bash \
    && rm -rf /var/lib/apt/lists/*

RUN groupadd -g 1000 smd && \
    useradd -r -u 1000 -g smd smd && \
    mkdir -p /data && \
    chown 1000:1000 /data

WORKDIR /home/smd

COPY bin/smd2-server /usr/local/bin/smd2-server

RUN chown -R smd:smd /home/smd

USER smd

ENTRYPOINT ["/usr/local/bin/smd2-server"]
CMD ["serve", "--port", "8080", "--database-url", "file:/data/smd2.db?_fk=1"]
