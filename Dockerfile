FROM golang:1.21 as build
WORKDIR /app
# Copy dependencies list
ADD go.mod go.sum ./
RUN go mod download
# Build with optional lambda.norpc tag
ADD cmd/lambda/proxy/main.go  main.go
RUN go build -tags lambda.norpc -o main main.go
# Copy artifacts to a clean image
FROM public.ecr.aws/lambda/provided:al2
COPY --from=build /app/main ./main
ENTRYPOINT [ "./main" ]