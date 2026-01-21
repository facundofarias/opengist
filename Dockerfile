FROM alpine:3.19

RUN echo "Step 1: installing tools..." \
  && apk add --no-cache curl

RUN echo "Step 2: intentional failure for testing" \
  && exit 1
