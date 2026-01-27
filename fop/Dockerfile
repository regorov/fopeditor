FROM golang:1.21-alpine AS builder
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/fopserver ./cmd/server

FROM eclipse-temurin:17-jre

ARG FOP_VERSION=2.9
ARG FOP_BASE_URL=https://archive.apache.org/dist/xmlgraphics/fop/binaries
ENV FOP_HOME=/opt/fop

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        curl \
        ca-certificates \
        fontconfig \
        fonts-dejavu-core \
        fonts-liberation \
        fonts-noto-core \
        fonts-noto-cjk \
    && rm -rf /var/lib/apt/lists/* \
    && curl -fsSL ${FOP_BASE_URL}/fop-${FOP_VERSION}-bin.tar.gz -o /tmp/fop.tar.gz \
    && mkdir -p /opt \
    && tar -xzf /tmp/fop.tar.gz -C /opt \
    && rm /tmp/fop.tar.gz \
    && mv /opt/fop-${FOP_VERSION} ${FOP_HOME} \
    && mkdir -p /usr/share/fonts/custom \
    && cat <<'EOF' >/usr/local/bin/fop
#!/bin/sh
export FOP_HOME=/opt/fop/fop
exec /opt/fop/fop/fop -c /etc/fop/fonts.xml "$@"
EOF
RUN chmod +x /usr/local/bin/fop
COPY times-new-roman-cyrillic.ttf /usr/share/fonts/custom/
COPY fonts.xml /etc/fop/fonts.xml
RUN fc-cache -fv
COPY --from=builder /out/fopserver /usr/local/bin/fopserver

EXPOSE 8090
ENTRYPOINT ["/usr/local/bin/fopserver"]
