FROM alpine:3.10
ARG BINARY=./fitness-backend
EXPOSE 3000

COPY ${BINARY} /opt/fitness-backend
RUN chmod +x /opt/fitness-backend
RUN apk add --no-cache tzdata
ENTRYPOINT ["/opt/fitness-backend"]