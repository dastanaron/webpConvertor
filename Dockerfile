FROM alpine:3.13.5

ARG USER=webp
ARG UID=1000
ARG GID=1000

USER root

RUN apk add --update --no-cache libc6-compat libpng-dev libjpeg-turbo-dev giflib-dev tiff-dev autoconf automake make gcc g++ wget

RUN mkdir -p /var/www/html/webpconvertor/data

#Installing webp lib
RUN wget https://storage.googleapis.com/downloads.webmproject.org/releases/webp/libwebp-0.6.0.tar.gz && \
tar -xvzf libwebp-0.6.0.tar.gz && \
mv libwebp-0.6.0 libwebp && \
rm libwebp-0.6.0.tar.gz && \
cd /libwebp && \
./configure && \
make && \
make install && \
rm -rf libwebp

COPY ./data/config.docker.yaml /var/www/html/webpconvertor/data/config.yaml

RUN set -eux; \
	addgroup -g ${GID} -S ${USER}; \
	adduser -u ${UID} -D -S -G ${USER} ${USER}

COPY ./build/webpConvertor /var/www/html/webpconvertor/run

RUN chown -R ${UID}:${GID} /var/www/html/webpconvertor \
    && chmod -R 775 /var/www/html/webpconvertor \
    && chown -R ${UID}:${GID} /run

RUN chmod +x /var/www/html/webpconvertor/run


USER ${UID}:${GID}

WORKDIR /var/www/html/webpconvertor

CMD ["/var/www/html/webpconvertor/run", "-c", "/var/www/html/webpconvertor/data/config.yaml", ">", "/dev/stderr"]
