# Start by building the application.
FROM golang:1.18 as build

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/app

# Now copy it into our base image.
FROM python:3-slim

ARG USERNAME=nonroot
ARG USER_UID=1000
ARG USER_GID=$USER_UID

# Create the user
RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME

# Install weasyprint
RUN apt-get -y update \
    && apt-get -y dist-upgrade \
    && apt-get install -y \
        libpango-1.0-0 libpangoft2-1.0-0 \
    && apt-get -y clean \
    && pip install weasyprint

USER $USERNAME

COPY --from=build /go/bin/app /

CMD ["/app", "serve"]
