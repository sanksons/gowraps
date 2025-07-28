# GoWraps

[![Go Reference](https://pkg.go.dev/badge/github.com/sanksons/gowraps.svg)](https://pkg.go.dev/github.com/sanksons/gowraps)
[![Go Report Card](https://goreportcard.com/badge/github.com/sanksons/gowraps)](https://goreportcard.com/report/github.com/sanksons/gowraps)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A collection of utility packages for Go that provide common functionality with clean, simple APIs. GoWraps aims to reduce boilerplate code and provide reliable, tested utilities for everyday Go development.

## 📦 Packages

### 🔐 Cipher
Cryptographic utilities for encryption, decryption, and encoding.

- **AES Encryption/Decryption** with CFB mode
- **Base64 Encoding/Decoding** utilities
- Secure key handling with 32-byte keys

```go
import "github.com/sanksons/gowraps/cipher"

// Base64 operations
encoded := cipher.Base64Encode("Hello, World!")
decoded, err := cipher.Base64Decode(encoded)

// AES encryption/decryption
var key [32]byte
copy(key[:], "your-32-byte-encryption-key-here")
encrypted, err := cipher.Encrypt(key, []byte("secret data"))
decrypted, err := cipher.Decrypt(key, encrypted)
```

### ⚡ Concurrency
Utilities for parallel execution with panic safety and result ordering.

- **Parallel function execution** with guaranteed result ordering
- **Panic recovery** for safe concurrent operations
- Simple API for complex parallel workflows

```go
import "github.com/sanksons/gowraps/concurrency"

functions := []func() interface{}{
    func() interface{} { return "Result 1" },
    func() interface{} { return "Result 2" },
    func() interface{} { return "Result 3" },
}

results := concurrency.Parallelize(functions)
// Results maintain the same order as input functions
```

### 🔄 Convert
Type conversion utilities with comprehensive type support.

- **Safe type conversions** between common Go types
- Support for `int`, `string`, `float64`, `float32`, etc.
- Error handling for invalid conversions

```go
import "github.com/sanksons/gowraps/convert"

// Convert various types to int
intVal, err := convert.ToInt("123")        // string to int
intVal, err := convert.ToInt(123.45)       // float to int
intVal, err := convert.ToInt(int64(123))   // int64 to int

// Convert various types to string
strVal, err := convert.ToString(123)       // int to string
strVal, err := convert.ToString(123.45)    // float to string
```

### 📁 Filesystem
File system operations with enhanced error handling.

- **File existence checking**
- **File reading/writing** with proper error categorization
- Custom error types for common file operations

```go
import "github.com/sanksons/gowraps/filesystem"

// Check if file exists
exists, err := filesystem.CheckIfFileExists("/path/to/file.txt")

// Read file content
content, err := filesystem.GetFile("/path/to/file.txt")
if err == filesystem.ErrFileNotFound {
    // Handle file not found
} else if err == filesystem.ErrPermissionDenied {
    // Handle permission denied
}
```

### 🗺️ HMap
Thread-safe concurrent hash map implementation.

- **Concurrent access** with RWMutex protection
- **Generic key-value storage** with `interface{}` types
- Rich API for map operations

```go
import "github.com/sanksons/gowraps/hmap"

// Create a new concurrent map
m := hmap.New()

// Thread-safe operations
m.Put("key1", "value1")
m.PutIfAbsent("key2", "value2")

value, exists := m.Get("key1")
m.Remove("key1")

// Iterate safely
m.Each(func(key, value interface{}) {
    fmt.Printf("%v: %v\n", key, value)
})
```

### 🖼️ Imaging
Image processing utilities for common formats.

- **MIME type detection** from byte data
- **Format conversion** between PNG, JPEG, GIF
- **Image encoding/decoding** with validation

```go
import "github.com/sanksons/gowraps/imaging"

// Detect image format
mime, err := imaging.GetMime(imageBytes)

// Convert between formats
ext, err := imaging.GetExtension4mMime("image/jpeg")
mime, err := imaging.GetMime4mExt("png")

// Image processing
img, err := imaging.PngToImage(pngBytes)
jpegBytes, err := imaging.ImageToBytesJpeg(img)
```

### 🗄️ MySQLDB
MySQL database wrapper with struct mapping and connection pooling.

- **Automatic struct mapping** for query results
- **Connection pooling** with configurable limits
- **Prepared statement support** for safety and performance
- Support for both single-row and multi-row operations

```go
import "github.com/sanksons/gowraps/mysqldb"

type User struct {
    Name       string
    Data       string
    Occupation *string
}

config := mysqldb.MySqlConfig{
    User:               "root",
    Passwd:             "password",
    Addr:               "localhost:3306",
    DBName:             "mydb",
    MaxOpenConnections: 10,
    MaxIdleConnections: 2,
}

pool, err := mysqldb.Initiate(config)
defer pool.Close()

// Query single row
var user User
err = pool.QuerySingle("SELECT name, data, occupation FROM users WHERE id = ?", &user, 1)

// Query multiple rows
var users []User
err = pool.Query("SELECT name, data, occupation FROM users", &users)
```

### 🔤 Regexp
Regular expression utilities for common text processing.

- **Text sanitization** for alphanumeric-only content
- **Pattern replacement** with custom replacements

```go
import "github.com/sanksons/gowraps/regexp"

// Remove non-alphanumeric characters
cleaned, err := regexp.AlphaNumericOnly("Hello, World! 123", "_")
// Result: "Hello__World__123"
```

### ⏰ Timer
Time manipulation and formatting utilities.

- **HTTP time formatting** (RFC1123)
- **Timezone-aware** current time retrieval
- Support for various timezone formats

```go
import "github.com/sanksons/gowraps/timer"

// Format time for HTTP headers
httpTime := timer.GetHttpTime(time.Now())

// Get current time in specific timezone
utcTime, err := timer.GetCurrentTime("UTC")
nyTime, err := timer.GetCurrentTime("America/New_York")
```

### 🛠️ Util
Database and query utilities.

- **Multi-insert query generation** for bulk operations
- **SQL helper functions** for common database tasks

```go
import "github.com/sanksons/gowraps/util"

fields := []string{"name", "email", "age"}
values := [][]interface{}{
    {"John", "john@example.com", 30},
    {"Jane", "jane@example.com", 25},
}

query := util.GetMultiInsertQuery("users", fields, values)
// Result: INSERT INTO users (name,email,age) VALUES (?,?,?),(?,?,?)
```

## 🚀 Installation

```bash
go get github.com/sanksons/gowraps
```

Or install specific packages:

```bash
go get github.com/sanksons/gowraps/cipher
go get github.com/sanksons/gowraps/concurrency
go get github.com/sanksons/gowraps/mysqldb
# ... etc
```

## 📖 Usage

Each package can be imported and used independently:

```go
package main

import (
    "fmt"
    "github.com/sanksons/gowraps/cipher"
    "github.com/sanksons/gowraps/concurrency"
    "github.com/sanksons/gowraps/timer"
)

func main() {
    // Use cipher for encoding
    encoded := cipher.Base64Encode("Hello, GoWraps!")
    
    // Use concurrency for parallel execution
    functions := []func() interface{}{
        func() interface{} { return "Task 1 Complete" },
        func() interface{} { return "Task 2 Complete" },
    }
    results := concurrency.Parallelize(functions)
    
    // Use timer for current time
    currentTime, _ := timer.GetCurrentTime("UTC")
    
    fmt.Printf("Encoded: %s\n", encoded)
    fmt.Printf("Results: %v\n", results)
    fmt.Printf("Current time: %v\n", currentTime)
}
```

## 🧪 Testing

The project includes comprehensive test suites for all packages:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./cipher
go test ./mysqldb
```

## 🏗️ Requirements

- **Go 1.24.4+**
- External dependencies:
  - `github.com/go-errors/errors v1.5.1`
  - `github.com/go-sql-driver/mysql v1.9.3`
  - `github.com/sanksons/go-reflexer v1.0.0`

## 📋 Examples

Check out the `examples/` directories in individual packages for more detailed usage examples:

- [MySQL Single Row Example](mysqldb/examples/singlerow/example.go)
- [MySQL Multi Row Example](mysqldb/examples/multirow/example.go)

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Add tests for your changes
4. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
5. Push to the branch (`git push origin feature/AmazingFeature`)
6. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📞 Support

If you have any questions or need help using GoWraps, please:

1. Check the [documentation](https://pkg.go.dev/github.com/sanksons/gowraps)
2. Look at the [examples](mysqldb/examples/)
3. Open an [issue](https://github.com/sanksons/gowraps/issues) on GitHub

## 🎯 Roadmap

- [ ] Add more image format support (WebP, TIFF)
- [ ] Extend cipher package with more encryption algorithms
- [ ] Add Redis wrapper similar to MySQL wrapper
- [ ] Performance optimizations for concurrent operations
- [ ] Additional utility functions based on community feedback

---

**GoWraps** - Making Go development more productive, one utility at a time! 🚀
