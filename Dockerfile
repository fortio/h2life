FROM scratch
COPY h2life /usr/bin/h2life
ENTRYPOINT ["/usr/bin/h2life"]
