# Go Exposed Config Scanner

This project is a multi-threaded tool designed to scan URLs for exposed configurations based on predefined templates.

## Features

- Load and scan templates from a configuration directory.
- Multi-threaded scanning with customizable thread count.
- Timeout settings for HTTP requests.
- Supports scanning specific templates or all templates.
- Outputs results to files.

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/faizdotid/go-exposed-config-scanner.git
   cd go-exposed-config-scanner
   ```

2. Install dependencies and build the project:

   ```bash
   go mod tidy
   cd cmd/go-scanner
   go build
   ```

## Usage

### Command-Line Arguments

- `-id`: Specify the template ID(s) to scan, comma-separated for multiple templates.
- `-all`: Scan all available templates.
- `-filelist`: List of URLs to scan.
- `-threads`: Number of threads to use for scanning, default 10 (***Threads will apllied for each templates***).
- `-timeout`: Timeout for HTTP requests in seconds (***This will be applied to all templates***).
- `-verbose`: Print error verbose.
- `-match`: Print only match URLs.
- `-show`: Display available templates and their details.

### Example Commands

1. **Scan using a specific template:**

   ```bash
   ./go-exposed-config-scanner -filelist urls.txt -id template1
   ```

2. **Scan using all templates:**

   ```bash
   ./go-exposed-config-scanner -filelist urls.txt -all
   ```

3. **Show available templates:**

   ```bash
   ./go-exposed-config-scanner -show
   ```

### Output

Results are stored in the `results` directory, with each template's output saved in a separate file.
