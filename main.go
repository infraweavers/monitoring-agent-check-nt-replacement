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
	"net/http"
	"os"
	"strconv"
	"time"
)

type CounterResult struct {
	Results []CounterResultItem
}
type CounterResultItem struct {
	CounterName  string
	InstanceName string
	Value        string
}

const okExitCode = 0
const warningExitCode = 1
const criticalExitCode = 2
const unknownExitCode = 3

var exitCodeToString = map[int]string{
	okExitCode:       "OK",
	warningExitCode:  "WARNING",
	criticalExitCode: "CRITICAL",
	unknownExitCode:  "UNKNOWN",
}

func die(stdout io.Writer, message string) int {
	fmt.Fprint(stdout, message)
	return unknownExitCode
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
	os.Exit(invokeClient(os.Stdout, httpClient))
}

func invokeClient(stdout io.Writer, httpClient httpclient.Interface) int {
	_ = flag.String("template", "", "pnp4nagios template")

	hostname := flag.String("host", "", "hostname or ip")
	port := flag.Int("port", 9000, "port number")
	username := flag.String("username", os.Getenv("MONITORING_AGENT_USERNAME"), "username")
	password := flag.String("password", os.Getenv("MONITORING_AGENT_PASSWORD"), "password")
	counterName := flag.String("counter", "", "counter path (i.e. \\PhysicalDisk(_Total)\\Avg. Disk Queue Length)")

	warningThreshold := flag.Float64("warning", -1, "warning threshold")
	criticalThreshold := flag.Float64("critical", -1, "critical threshold")

	counterlabel := flag.String("label", "", "output label")
	counterUnit := flag.String("unit", "%", "unit of measurement")

	cacertificateFilePath := flag.String("cacert", os.Getenv("MONITORING_AGENT_CA_CERTIFICATE_PATH"), "CA certificate")
	certificateFilePath := flag.String("certificate", os.Getenv("MONITORING_AGENT_CLIENT_CERTIFICATE_PATH"), "certificate file")
	privateKeyFilePath := flag.String("key", os.Getenv("MONITORING_AGENT_CLIENT_KEY_PATH"), "key file")
	timeoutString := flag.String("timeout", "10s", "timeout (e.g. 10s)")
	makeInsecure := flag.Bool("insecure", false, "ignore TLS Certificate checks")

	flag.Parse()

	if *hostname == "" {
		return die(stdout, "hostname is not set")
	}
	if *password == "" {
		return die(stdout, "password is not set")
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
			return die(stdout, fmt.Sprintf("error loading certificate pair %s", err.Error()))
		}
		transport.TLSClientConfig.Certificates = []tls.Certificate{certificateToLoad}
	}

	if *cacertificateFilePath != "" {
		caCertificate, err := ioutil.ReadFile(*cacertificateFilePath)
		if err != nil {
			return die(stdout, fmt.Sprintf("error loading ca certificate %s", err.Error()))
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
		return die(stdout, fmt.Sprintf("got httpClient error %s", err.Error()))
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		errorBodyContent, _ := ioutil.ReadAll(response.Body)
		return die(stdout, fmt.Sprintf("Response code: %d\n%s", response.StatusCode, errorBodyContent))
	}

	var decodedResponse CounterResult

	decoder := json.NewDecoder(response.Body)
	decoder.DisallowUnknownFields()
	decoder.Decode(&decodedResponse)

	outputValue := decodedResponse.Results[0].Value
	outputFloatValue, err := strconv.ParseFloat(outputValue, 64)

	outputCode := unknownExitCode

	if err != nil {

	} else {

		if *criticalThreshold >= *warningThreshold {
			/*
				on the number line this looks like:
				|-----------------------
				0             w    c
			*/
			if outputFloatValue > *criticalThreshold {
				outputCode = criticalExitCode
			} else if outputFloatValue > *warningThreshold && outputFloatValue <= *criticalThreshold {
				outputCode = warningExitCode
			} else {
				outputCode = okExitCode
			}

		} else if *criticalThreshold < *warningThreshold {
			/*
				on the number line this looks like:
				|-----------------------
				0    c    w
			*/
			if outputFloatValue < *criticalThreshold {
				outputCode = criticalExitCode
			} else if outputFloatValue < *warningThreshold && outputFloatValue >= *criticalThreshold {
				outputCode = warningExitCode
			} else {
				outputCode = okExitCode
			}

		}

		fmt.Printf("%s = %s %s | '%s'=%s%s;%f;%f;\n", *counterlabel, outputValue, *counterUnit, *counterlabel, outputValue, *counterUnit, *warningThreshold, *criticalThreshold)
	}

	return outputCode
}
