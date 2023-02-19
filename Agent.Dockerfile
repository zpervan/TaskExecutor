# Image
FROM golang:1.19.1

# Define the work directory
WORKDIR ./Agent

# Copy the project folders into the container's working directory
COPY . .

# Build it
RUN go build taskexecutor/agent

# Give executable privileges
RUN chmod +x agent

# Run the task execution agent
CMD [ "./agent" ]