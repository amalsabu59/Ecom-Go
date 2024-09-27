# Use the latest Red Hat Universal Base Image
FROM registry.access.redhat.com/ubi8/ubi:latest

# Install dependencies
RUN yum -y install wget tar git && \
    yum clean all

# Install Go 1.23.1
RUN wget https://golang.org/dl/go1.23.1.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.23.1.linux-amd64.tar.gz && \
    rm go1.23.1.linux-amd64.tar.gz

# Set up Go environment variables
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Install the correct version of Go Air
RUN  go install github.com/cosmtrek/air@v1.50.0

# Copy the entire project code
COPY . .

# Expose the port your app will run on
EXPOSE 8081

# Use entrypoint for easier debugging
ENTRYPOINT ["air", "-c", "cmd/.air.toml"]
