#Using an official Golang runtime image that is used across most containers in DSCC
FROM cds-harbor.rtplab.nimblestorage.com/docker_proxy/library/golang:1.19-alpine3.16 as base

#Setting the working Directory
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Build the Golang application
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Define environment variables to pass to the Golang application
ENV PORT=8080
ENV GIN_MODE=release

# Run the Golang application
CMD ["./main"]
