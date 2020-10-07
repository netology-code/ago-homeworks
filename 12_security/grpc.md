# gRPC

Для gRPC всё достаточно просто: вы прописываете пути к файлу сертификата и файлу ключа (в примере - каталог `tls`).

## Server

```go
const defaultPort = "9999"
const defaultHost = "0.0.0.0"
const defaultCertificatePath = "./tls/certificate.pem"
const defaultPrivateKeyPath = "./tls/key.pem"

func main() {
	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("APP_HOST")
	if !ok {
		host = defaultHost
	}

	certificatePath, ok := os.LookupEnv("APP_CERT_PATH")
	if !ok {
		certificatePath = defaultCertificatePath
	}

	privateKeyPath, ok := os.LookupEnv("APP_PRIVATE_KEY_PATH")
	if !ok {
		privateKeyPath = defaultPrivateKeyPath
	}

	if err := execute(net.JoinHostPort(host, port), certificatePath, privateKeyPath); err != nil {
		os.Exit(1)
	}
}

func execute(addr string, certificatePath string, privateKeyPath string) (err error) {
	creds, err := credentials.NewServerTLSFromFile(certificatePath, privateKeyPath)
	if err != nil {
		log.Fatalf("Failed to setup TLS: %v", err)
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	server := app.NewServer()
	eventV1Pb.RegisterEventServiceServer(grpcServer, server)

	return grpcServer.Serve(listener)
}
```

## Client

Клиенту нужен только сертификат:

```go
const defaultPort = "9999"
const defaultHost = "0.0.0.0"
const defaultCertificatePath = "./tls/certificate.pem"

func main() {
	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("APP_HOST")
	if !ok {
		host = defaultHost
	}

	certificatePath, ok := os.LookupEnv("APP_CERTIFICATE_PATH")
	if !ok {
		certificatePath = defaultCertificatePath
	}

	if err := execute(net.JoinHostPort(host, port), certificatePath); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func execute(addr string, certificatePath string) (err error) {
	creds, err := credentials.NewClientTLSFromFile(certificatePath, "")
	if err != nil {
		return err
	}

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Print(err)
		}
	}()

	client := eventV1Pb.NewEventServiceClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	response, err := client.Unary(ctx, &eventV1Pb.EventRequest{Id: 1, Payload: "Request"})
	if err != nil {
		return err
	}

	log.Print(response)
	return nil
}
```
