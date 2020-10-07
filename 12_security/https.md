# HTTPS

С HTTPS немного посложнее: мы должны будем клиенту явно прописать сертификат, чтобы он "не ругался". Т.к. по умолчанию используются системные значения, а у нас - самоподписанный сертификат.

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
		log.Print(err)
		os.Exit(1)
	}
}

func execute(addr string, certificatePath string, privateKeyPath string) (err error) {
	return http.ListenAndServeTLS(addr, certificatePath, privateKeyPath, &handler{});
}
```

## Client

Клиенту мы заменяем корневые сертификаты на самоподписанный:

```go
const defaultPort = "9999"
const defaultHost = "netology.local"
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
	certificate, err := ioutil.ReadFile(certificatePath)
	if err != nil {
		return err
	}

	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(certificate) {
		return errors.New("certificate not added")
	}

	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: pool,
		},
	}}

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://%s", addr), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("status not 200")
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Print(cerr)
		}
	}()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Printf("%s", data)

	return nil
}
```
