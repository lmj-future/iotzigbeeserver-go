FROM --platform=$TARGETPLATFORM alpine
ARG TARGETPLATFORM
COPY $TARGETPLATFORM/bin/ /workspace/server/
COPY $TARGETPLATFORM/config/ /workspace/server/config/
WORKDIR /workspace/server
ENTRYPOINT ["/workspace/server/iot-zigbee"]