# PromGen: Prometheus Listener Generator

## Introduction
PromGen is a command-line tool written in Go, designed to generate fully functional Prometheus listener code for various monitoring scenarios. It supports generating tests for CPU load, certificate expiry, and host availability.

## Requirements
- Go (1.15 or later)

## Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/promgen.git
Navigate to the PromGen directory:
cd promgen
Compile the program:

go build promgen.go


Usage
To use PromGen, run it with the required flags. Below are examples for different types of tests:

CPU Load Test
Generate a Prometheus listener for CPU load:

./promgen -test-type "CPU Load" -interval 1m -threshold 0.75 -port 9090 cputest.go

Host Up/Down Test
Generate a Prometheus listener to check if a host is up or down:

./promgen -test-type "host-up" -interval 1m -port 9090 -host "example.com" hosttest.go

Certificate Expiry Test
Generate a Prometheus listener for SSL certificate expiry:

./promgen -test-type "cert-expiry" -interval 1m -port 9090 -domain "example.com" certtest.go

After generating the test file (e.g., cputest.go, hosttest.go, certtest.go):

Navigate to the directory containing the generated file.

Initialize a new Go module:

go mod init your/module

Replace your/module with a unique module name for your project.

Add the Prometheus client dependency:

go get github.com/prometheus/client_golang@latest

Build and run the generated file:

go build <generated-file>.go
./<generated-file>

Customization
PromGen allows customization for various test parameters:

-test-type: The type of test (e.g., "CPU Load", "cert-expiry", "host-up").
-interval: The interval between checks (e.g., 30s, 1m, 5m).
-threshold: The threshold value for CPU load triggering an alert or action.
-port: The port on which the Prometheus listener will run.
-domain: The domain name for the SSL certificate expiry check.
-host: The host address for the up/down check.
Output
The output is a .go file containing the complete code for a Prometheus listener, ready to be compiled and run.

Additional test type examples:

Load Average Test
For generating a Prometheus listener to monitor system load average:

./promgen -test-type "load-average" -interval 1m -threshold 1.5 -port 9090 loadtest.go

-test-type "load-average": Specifies that the test is for monitoring load average.
-interval 1m: Sets the check interval to 1 minute.
-threshold 1.5: Sets the threshold for the load average, above which an alert or action should be triggered.
-port 9090: Specifies the port for the Prometheus listener.
loadtest.go: The output file where the generated code will be saved.

Swap Usage Test
For generating a Prometheus listener to monitor swap usage:
./promgen -test-type "swap-usage" -interval 5m -threshold 2048 -port 9090 swapttest.go


