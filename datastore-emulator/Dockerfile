FROM debian:stable-slim

LABEL maintainer="Marcos Jauregui <it.marcosj@gmail.com>"
LABEL title="Datastore emulator"
LABEL description="GCP Datastore emulator based on Debian stable slim"

# SETUP ----------------------------------------------------------

ARG DS_USER="dsuser"
ARG DS_HOME="/home/${DS_USER}"

# DEPENDENCIES
RUN set -eux \
    && DEBIAN_FRONTEND=noninteractive apt --yes update \
    && DEBIAN_FRONTEND=noninteractive apt-get install --yes --no-install-recommends \
        curl \
        default-jre \
    # Clean installation
    && DEBIAN_FRONTEND=noninteractive apt-get autoremove --yes \
    && DEBIAN_FRONTEND=noninteractive apt-get clean \
    && rm -frv /var/lib/apt/lists/* \
        /tmp/* \
        /var/tmp/* \ 
        ~/.cache \
    # USER
    && useradd --comment "datastore emulator user" \
        --create-home \
        --home-dir ${DS_HOME} \        
        ${DS_USER}

# EMULATOR ----------------------------------------------------------

USER ${DS_USER}

ARG GCLOUD_VERSION="458.0.1"
ARG GCLOUD_ARCH="x86_64"
ARG GCLOUD_FILE="google-cloud-cli-${GCLOUD_VERSION}-linux-${GCLOUD_ARCH}.tar.gz"

ENV CLOUDSDK_CORE_PROJECT="test-project"
ENV DS_HOST="0.0.0.0"
ENV GCLOUD_DIR="${DS_HOME}/google-cloud-sdk"
ENV GCLOUD_STORE_ON_DISK=false

RUN set -eux \
    && curl -o ${DS_HOME}/${GCLOUD_FILE} https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/${GCLOUD_FILE} \
    && tar -xvf ${DS_HOME}/${GCLOUD_FILE} -C ${DS_HOME}/ \
    && rm -v ${DS_HOME}/${GCLOUD_FILE} \
    && ${GCLOUD_DIR}/install.sh \
        --quiet \
        --usage-reporting false \
        --path-update true \
        --command-completion true \ 
    && export PATH=$PATH:${GCLOUD_DIR}/bin \
    && gcloud components install \
        beta \
        cloud-datastore-emulator \
        --quiet

# ----------------------------------------------------------

COPY ./entrypoint.sh /entrypoint.sh

EXPOSE 8081

ENTRYPOINT [ "./entrypoint.sh" ]
