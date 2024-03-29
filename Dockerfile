# Use an Alpine base image
FROM alpine:latest

# Install Git
RUN apk --no-cache add git

# Clone the repository
RUN git clone https://github.com/gunnarmorling/1brc.git

# Change to the cloned directory
WORKDIR /1brc/src/main/python

# Install Python 3
RUN apk add --no-cache python3

# Run the Python script to generate the data
RUN python3 create_measurements.py 1000000000