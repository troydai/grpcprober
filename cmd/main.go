package main

import (
	"context"
	"crypto/x509"
	"log"
	"os"
	"os/signal"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	pb "github.com/troydai/grpcprober/gen/api/protos"
)

const _beaconDemoCert = `-----BEGIN CERTIFICATE-----
MIICzjCCAjCgAwIBAgIUelbDmVczSZ6zTDiNBTM9Qv4QbCIwCgYIKoZIzj0EAwIw
bjELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAldBMRAwDgYDVQQHDAdSZWRtb25kMQ8w
DQYDVQQKDAZUREZ1bmQxETAPBgNVBAsMCEtleVNtaXRoMRwwGgYDVQQDDBNrZXlz
bWl0aC50cm95ZGFpLmNjMB4XDTI0MDgyNDIxMTk1MFoXDTI0MDkyMzIxMTk1MFow
bjELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAldBMRAwDgYDVQQHDAdSZWRtb25kMQ8w
DQYDVQQKDAZUREZ1bmQxETAPBgNVBAsMCEtleVNtaXRoMRwwGgYDVQQDDBNrZXlz
bWl0aC50cm95ZGFpLmNjMIGbMBAGByqGSM49AgEGBSuBBAAjA4GGAAQBjmDTB1fu
lnOXQYr3ZbJcDMxQK427c+sodhsIJbgY7h5pfIVFhDUpPczy88cimz8sELLKmyOh
AfXft8wkaO6tPn0A3S6lnyhak36NFsxsbZkp9QHEFLCrtCgPXqiXMBJ5icRYutqF
RIBCTGKM+NH9Nn/ekhP6817Wfa2iZkO1oanI2MmjaTBnMB0GA1UdDgQWBBQ8bgwv
MBkiJIHlofjSpJbf7qL9HzAfBgNVHSMEGDAWgBQ8bgwvMBkiJIHlofjSpJbf7qL9
HzAPBgNVHRMBAf8EBTADAQH/MBQGA1UdEQQNMAuCCWxvY2FsaG9zdDAKBggqhkjO
PQQDAgOBiwAwgYcCQgCJrmIRkXMIy3qjQD8e74JiIfInCIySFkQfzAxFnFzqVap3
Gq6q+yxqGWPkX5aooCglZBMq3t8zjwO2KXIImR/wRwJBB5z+EHFKqHadDpoiRoIG
OFCuIvDRtq6U6j20s/e0rno4lkiuc7MblNRWkKeIuEtu1nYfyjaBsszI6FfgKrlm
+Zk=
-----END CERTIFICATE-----`

func main() {
	conn, err := grpc.NewClient(
		"localhost:8080",
		grpc.WithTransportCredentials(getTransportCredential()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a new gRPC client for the beacon service.
	client := pb.NewBeaconServiceClient(conn)

	exit := make(chan struct{})

	go func() {
		defer close(exit)
		term := make(chan os.Signal, 1)
		signal.Notify(term, os.Interrupt)

		<-term
	}()

	t := time.NewTicker(10 * time.Minute)
	for {
		// Make a gRPC request to the beacon service.
		response, err := client.Signal(context.Background(), &pb.SignalRequest{})
		if err != nil {
			log.Fatalf("Failed to get beacon: %v", err)
		}

		// Process the response from the beacon service.
		log.Printf("Received beacon: %v", response.GetDetails())

		select {
		case <-exit:
			return
		case <-t.C:
		}
	}
}

func getTransportCredential() credentials.TransportCredentials {
	if os.Getenv("PROBER_TLS") == "true" {
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM([]byte(_beaconDemoCert))
		return credentials.NewClientTLSFromCert(certPool, "")
	}

	return insecure.NewCredentials()
}
