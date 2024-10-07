# Browser-less bundle
FROM alpine:3.18 AS release

RUN adduser -D -u 12345 -g 12345 k6

COPY --chmod=755 {{.Executable}} /usr/bin/k6

USER k6
WORKDIR /home/k6

ENTRYPOINT ["k6"]

# Browser-enabled bundle
FROM release AS with-browser

USER root

COPY --chmod=755 --from=release /usr/bin/k6 /usr/bin/k6
RUN apk --no-cache add chromium-swiftshader

USER k6

ENV CHROME_BIN=/usr/bin/chromium-browser
ENV CHROME_PATH=/usr/lib/chromium/

ENV K6_BROWSER_HEADLESS=true
# no-sandbox chrome arg is required to run chrome browser in
# alpine and avoids the usage of SYS_ADMIN Docker capability
ENV K6_BROWSER_ARGS=no-sandbox

ENTRYPOINT ["k6"]
