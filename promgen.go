package main

import (
    "flag"
    "fmt"
    "os"
    "text/template"
    "time"
)

// PrometheusListener represents the structure of the generated listener code
type PrometheusListener struct {
    TestType    string
    Interval    time.Duration
    Threshold   float64
    Port        int
    Domain      string // For cert expiry check
    HostAddress string // For host up/down check
}

// listenerTemplate is the Go template for generating the Prometheus listener code
const listenerTemplate = `
package main

import (
    "crypto/tls"
    "net"
    "net/http"
    "time"
    "log"
    "fmt"

    "github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusListener struct {
    TestType    string
    Interval    time.Duration
    Threshold   float64
    Port        int
    Domain      string
    HostAddress string
}

{{ if eq .TestType "cert-expiry" }}
func CheckCertificateExpiry(domain string) (time.Duration, error) {
    conn, err := tls.Dial("tcp", domain+":443", nil)
    if err != nil {
        return 0, err
    }
    defer conn.Close()

    cert := conn.ConnectionState().PeerCertificates[0]
    return cert.NotAfter.Sub(time.Now()), nil
}
{{ end }}

{{ if eq .TestType "host-up" }}
func IsHostUp(host string) bool {
    _, err := net.DialTimeout("tcp", host, time.Second*10)
    return err == nil
}
{{ end }}

func NewPrometheusListener(testType string, interval time.Duration, threshold float64, port int, domain string, hostAddress string) *PrometheusListener {
    return &PrometheusListener{
        TestType:    testType,
        Interval:    interval,
        Threshold:   threshold,
        Port:        port,
        Domain:      domain,
        HostAddress: hostAddress,
    }
}

func (p *PrometheusListener) Start() {
    http.Handle("/metrics", promhttp.Handler())
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", p.Port), nil))
}

func main() {
    listener := NewPrometheusListener("{{ .TestType }}", {{ .Interval }}, {{ .Threshold }}, {{ .Port }}, "{{ .Domain }}", "{{ .HostAddress }}")
    listener.Start()
}
`

func main() {
    testType := flag.String("test-type", "", "Type of test to monitor (cert-expiry, host-up)")
    interval := flag.Duration("interval", 1*time.Minute, "Check interval")
    threshold := flag.Float64("threshold", 0.0, "Threshold for alerts")
    port := flag.Int("port", 9090, "Port for Prometheus listener")
    domain := flag.String("domain", "", "Domain name for certificate expiry check")
    hostAddress := flag.String("host", "", "Host address for up/down check")

    flag.Parse()

    args := flag.Args()
    if len(args) == 0 {
        fmt.Println("Error: No output file specified")
        os.Exit(1)
    }
    outputFileName := args[len(args)-1]

    listener := PrometheusListener{
        TestType:    *testType,
        Interval:    *interval,
        Threshold:   *threshold,
        Port:        *port,
        Domain:      *domain,
        HostAddress: *hostAddress,
    }

    tmpl, err := template.New("listener").Parse(listenerTemplate)
    if err != nil {
        fmt.Println("Error parsing template:", err)
        os.Exit(1)
    }

    outputFile, err := os.Create(outputFileName)
    if err != nil {
        fmt.Println("Error creating output file:", err)
        os.Exit(1)
    }
    defer outputFile.Close()

    if err := tmpl.Execute(outputFile, listener); err != nil {
        fmt.Println("Error executing template:", err)
        os.Exit(1)
    }

    fmt.Printf("Prometheus listener code generated in %s\n", outputFileName)
}
