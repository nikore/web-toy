FROM golang:1.11.5

RUN cd /tmp \
    && curl -LO http://http.us.debian.org/debian/pool/main/r/recode/librecode0_3.6-21_amd64.deb \
    && curl -LO http://http.us.debian.org/debian/pool/main/f/fortune-mod/fortune-mod_1.99.1-7_amd64.deb \
    && curl -LO http://http.us.debian.org/debian/pool/main/f/fortune-mod/fortunes-min_1.99.1-7_all.deb \
    && dpkg -i librecode0_3.6-21_amd64.deb fortune-mod_1.99.1-7_amd64.deb fortunes-min_1.99.1-7_all.deb

COPY run.sh /run.sh

RUN mkdir -p /opt/web-toy

WORKDIR /opt/web-toy

ENV PATH /opt/web-toy/bin:$PATH

COPY . /opt/web-toy/

EXPOSE 8080

RUN make clean && make build

CMD ["/opt/web-toy/bin/web-toy"]
