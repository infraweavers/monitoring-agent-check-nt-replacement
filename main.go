package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"monitoring-agent-client-check-nt-replacement/internal/httpclient"
	"monitoring-agent-client-check-nt-replacement/internal/nagios"
	"net/http"
	"os"
	"time"
)

type stringFlag struct {
	set   bool
	value string
}

func (sf *stringFlag) Set(x string) error {
	sf.value = x
	sf.set = true
	return nil
}

func (sf *stringFlag) String() string {
	return sf.value
}

type CounterResult struct {
	Results []CounterResultItem
}
type CounterResultItem struct {
	CounterName  string
	InstanceName string
	Value        string
}

func die(stdout io.Writer, message string) int {
	fmt.Fprint(stdout, message)
	return nagios.StateUNKNOWNExitCode
}

func enableTimeout(timeout string) time.Duration {
	timeoutDuration, timeoutParseError := time.ParseDuration(timeout)
	if timeoutParseError != nil {
		panic(fmt.Errorf("error parsing timeout value %s", timeoutParseError.Error()))
	}

	time.AfterFunc(timeoutDuration, func() {
		panic(fmt.Sprintf("Client timeout reached: %s\n", timeoutDuration))
	})
	return timeoutDuration
}

func main() {
	httpClient := httpclient.NewHTTPClient()
	invokeClient(os.Stdout, httpClient)
}

func invokeClient(stdout io.Writer, httpClient httpclient.Interface) {

	var plugin = nagios.Plugin{
		ExitStatusCode: nagios.StateUNKNOWNExitCode,
	}

	defer plugin.ReturnCheckResults()

	hostname := flag.String("host", "", "hostname or ip")
	port := flag.Int("port", 9000, "port number")
	username := flag.String("username", os.Getenv("MONITORING_AGENT_USERNAME"), "username")
	password := flag.String("password", os.Getenv("MONITORING_AGENT_PASSWORD"), "password")
	counterName := flag.String("counter", "", "counter path (i.e. \\PhysicalDisk(_Total)\\Avg. Disk Queue Length)")

	var warningThreshold stringFlag
	var criticalThreshold stringFlag

	flag.Var(&warningThreshold, "warning", "warning threshold")
	flag.Var(&criticalThreshold, "critical", "critical threshold")

	if warningThreshold.set {
		plugin.WarningThreshold = warningThreshold.value
	}
	if criticalThreshold.set {
		plugin.CriticalThreshold = criticalThreshold.value
	}

	//counterlabel := flag.String("label", "", "output label")
	counterUnit := flag.String("unit", "%", "unit of measurement")

	cacertificateFilePath := flag.String("cacert", os.Getenv("MONITORING_AGENT_CA_CERTIFICATE_PATH"), "CA certificate")
	certificateFilePath := flag.String("certificate", os.Getenv("MONITORING_AGENT_CLIENT_CERTIFICATE_PATH"), "certificate file")
	privateKeyFilePath := flag.String("key", os.Getenv("MONITORING_AGENT_CLIENT_KEY_PATH"), "key file")
	timeoutString := flag.String("timeout", "10s", "timeout (e.g. 10s)")
	makeInsecure := flag.Bool("insecure", false, "ignore TLS Certificate checks")

	flag.Parse()

	if *hostname == "" {
		die(stdout, "hostname is not set")
		return
	}
	if *password == "" {
		die(stdout, "password is not set")
		return
	}

	timeout := enableTimeout(*timeoutString)

	restRequest := map[string]interface{}{
		"CounterPath": counterName,
	}

	url := fmt.Sprintf("https://%s:%d/v1/os_specific", *hostname, *port)

	httpClient.SetTimeout(timeout)

	transport := new(http.Transport)
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: *makeInsecure,
	}

	if *certificateFilePath != "" && *privateKeyFilePath != "" {
		certificateToLoad, err := tls.LoadX509KeyPair(*certificateFilePath, *privateKeyFilePath)
		if err != nil {
			die(stdout, fmt.Sprintf("error loading certificate pair %s", err.Error()))
			return
		}
		transport.TLSClientConfig.Certificates = []tls.Certificate{certificateToLoad}
	}

	if *cacertificateFilePath != "" {
		caCertificate, err := ioutil.ReadFile(*cacertificateFilePath)
		if err != nil {
			die(stdout, fmt.Sprintf("error loading ca certificate %s", err.Error()))
			return
		}
		CACertificatePool := x509.NewCertPool()
		CACertificatePool.AppendCertsFromPEM(caCertificate)
		transport.TLSClientConfig.RootCAs = CACertificatePool
	}

	httpClient.SetTransport(transport)

	byteArray, _ := json.Marshal(restRequest)
	byteArrayBuffer := bytes.NewBuffer(byteArray)

	req, err := http.NewRequest(http.MethodPost, url, byteArrayBuffer)
	if err != nil {
		panic(fmt.Errorf("got http request error %s", err.Error()))
	}
	req.SetBasicAuth(*username, *password)

	response, err := httpClient.Do(req)

	if err != nil {
		die(stdout, fmt.Sprintf("got httpClient error %s", err.Error()))
		return
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		errorBodyContent, _ := ioutil.ReadAll(response.Body)
		die(stdout, fmt.Sprintf("Response code: %d\n%s", response.StatusCode, errorBodyContent))
		return
	}

	var decodedResponse CounterResult

	decoder := json.NewDecoder(response.Body)
	decoder.DisallowUnknownFields()
	decoder.Decode(&decodedResponse)

	plugin.ExitStatusCode = nagios.StateOKExitCode

	for _, outputValue := range decodedResponse.Results {

		perfdata := nagios.PerformanceData{
			Label:             outputValue.InstanceName,
			Value:             outputValue.Value,
			UnitOfMeasurement: *counterUnit,
		}
		if warningThreshold.set {
			perfdata.Warn = warningThreshold.value
		}
		if criticalThreshold.set {
			perfdata.Crit = criticalThreshold.value
		}
		plugin.AddPerfData(false, perfdata)
		plugin.EvaluateThreshold(perfdata)
	}

	plugin.ServiceOutput = nagios.StateOKLabel

	if plugin.ExitStatusCode == nagios.StateWARNINGExitCode {
		plugin.ServiceOutput = nagios.StateWARNINGLabel
	}

	if plugin.ExitStatusCode == nagios.StateCRITICALExitCode {
		plugin.ServiceOutput = nagios.StateCRITICALLabel
	}

	if plugin.ExitStatusCode == nagios.StateUNKNOWNExitCode {
		plugin.ServiceOutput = nagios.StateUNKNOWNLabel
	}
}
