FROM alpine:3.10
ARG BINARY=./fitness-backend
EXPOSE 3000

COPY ${BINARY} /opt/fitness-backend
ENTRYPOINT ["/opt/fitness-backend"]