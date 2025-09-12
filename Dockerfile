
---

---

### ✅ 2️⃣: Dockerfile (Optional if you want containerized setup)

```dockerfile
# Base image
FROM golang:1.25

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN go build -o myapp

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./myapp"]
