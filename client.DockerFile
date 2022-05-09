FROM scratch

COPY build/client /bin/pow_client

CMD ["/bin"]
