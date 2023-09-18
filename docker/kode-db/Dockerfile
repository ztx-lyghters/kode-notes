FROM postgres:15.4-alpine3.18
COPY ./init.sql /docker-entrypoint-initdb.d/
RUN mkdir -p /usr/local/bin
COPY ./first-run /usr/local/bin/
RUN chmod +x /usr/local/bin/first-run
ENTRYPOINT ["bash"]
