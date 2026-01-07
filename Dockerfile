FROM golang:1.24-alpine

ARG TOKEN
ENV GOPRIVATE="github.com/Out-Of-India-Theory"

RUN apk add --no-cache \
    git \
    chromium \
    curl \
    fontconfig \
    ttf-freefont \
    unzip \
    && ln -sf /usr/bin/chromium-browser /usr/bin/google-chrome

# Create font directory
RUN mkdir -p /usr/share/fonts/noto

# Download only the specific fonts you need
WORKDIR /usr/share/fonts/noto

# Download required fonts
RUN curl -L -o NotoSansTamil-Regular.ttf https://github.com/googlefonts/noto-fonts/raw/main/hinted/ttf/NotoSansTamil/NotoSansTamil-Regular.ttf && \
    curl -L -o NotoSansTelugu-Regular.ttf https://github.com/googlefonts/noto-fonts/raw/main/hinted/ttf/NotoSansTelugu/NotoSansTelugu-Regular.ttf && \
    curl -L -o NotoSansDevanagari-Regular.ttf https://github.com/googlefonts/noto-fonts/raw/main/hinted/ttf/NotoSansDevanagari/NotoSansDevanagari-Regular.ttf && \
    curl -L -o NotoSansKannada-Regular.ttf https://github.com/googlefonts/noto-fonts/raw/main/hinted/ttf/NotoSansKannada/NotoSansKannada-Regular.ttf && \
    curl -L -o NotoSansGujarati-Regular.ttf https://github.com/googlefonts/noto-fonts/raw/main/hinted/ttf/NotoSansGujarati/NotoSansGujarati-Regular.ttf

# Rebuild font cache
RUN fc-cache -f -v

RUN go env -w GOPRIVATE="github.com/Out-Of-India-Theory" \
    && git config --global url."https://oit-devops:${TOKEN}@github.com".insteadOf "https://github.com"

COPY . /go/src/github.com/Out-Of-India-Theory/helper-service

WORKDIR /go/src/github.com/Out-Of-India-Theory/helper-service

RUN echo $GOPRIVATE

RUN go mod tidy \
    && go mod download

RUN GOOS=linux GOARCH=amd64 go build -o main .

EXPOSE 8080

CMD ["./main"]
