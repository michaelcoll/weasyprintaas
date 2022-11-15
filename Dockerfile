# Start by building the application.
FROM golang:1.19 as build

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/app

# Now copy it into our base image.
FROM alpine:3

ARG USERNAME=nonroot
ARG USER_UID=1000
ARG USER_GID=$USER_UID

# Create the user
RUN addgroup --gid $USER_GID -S $USERNAME \
    && adduser -u $USER_UID -S $USERNAME -G $USERNAME

# Install weasyprint
RUN apk --update --upgrade --no-cache add py3-pip py3-pillow py3-cffi py3-brotli gcc musl-dev python3 pango fontconfig font-noto \
    && pip install weasyprint \
    && apk del py3-pip py3-brotli gcc musl-dev \
    && apk add py3-six

USER $USERNAME

COPY --from=build /go/bin/app /

CMD ["/app", "serve"]
