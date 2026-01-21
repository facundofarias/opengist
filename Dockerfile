RUN make

# Intentional failure for testing
RUN echo "Failing build on purpose" && exit 1
