FROM scratch

ENV PORT 32001
EXPOSE $PORT

COPY hello-world /
