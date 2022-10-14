FROM golang:1.19-alpine

# Create a directory for the app
RUN mkdir /app
 
# Copy all files from the current directory to the app directory
COPY . /app
 
# Set working directory
WORKDIR /app

# Set environnement variables
ARG APP_MODE
ARG PORT
ARG ALLOWED_ORIGIN
ARG DSN
ARG JWT_SECRET

ENV APP_MODE=${APP_MODE} \
    PORT=${PORT} \
    ALLOWED_ORIGIN=${ALLOWED_ORIGIN} \
    DSN=${DSN} \
    JWT_SECRET=${JWT_SECRET}
    
# go build will build an executable file named server in the current directory
RUN set GOOS=linux && set GOARCH=amd64 && set CGO_ENABLED=0 && go build -o ./bin/tenidditerApi ./cmd/api

# Run the server executable
CMD [ "/app/bin/tenidditerApi" ]