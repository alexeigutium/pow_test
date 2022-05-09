FROM scratch

COPY data/quotes.txt /data/quotes.txt
COPY build/server /bin/pow_server

CMD ["/bin"]
