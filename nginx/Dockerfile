FROM nginx:latest

RUN mkdir -p /data/nginx/cache/webp

COPY nginx.conf /etc/nginx/
COPY conf.d/* /etc/nginx/conf.d/

RUN mkdir -p /var/www/html
RUN cp /usr/share/nginx/html/index.html /var/www/html

# forward request and error logs to docker log collector
RUN ln -sf /dev/stdout /var/log/nginx/access.log \
	&& ln -sf /dev/stderr /var/log/nginx/error.log

EXPOSE 80

STOPSIGNAL SIGTERM

CMD ["nginx", "-g", "daemon off;"]
