
openssl genrsa -out test_ca.key 2048
openssl req -x509 -new -nodes -key test_ca.key -sha256 -days 2048 -out test_ca.crt -subj "/CN=My CA"


openssl genrsa -out test_server.key 2048
openssl req -new -key test_server.key -out test_server.csr -subj "/CN=localhost"


openssl x509 -req -in test_server.csr -CA test_ca.crt -CAkey test_ca.key -CAcreateserial -out test_server.crt -days 500 -sha256  -extfile <(printf "subjectAltName=DNS:localhost,DNS:example.com")


openssl genrsa -out test_client.key 2048
openssl req -new -key test_client.key -out test_client.csr -subj "/CN=Client"


openssl x509 -req -in test_client.csr -CA test_ca.crt -CAkey test_ca.key -CAcreateserial -out test_client.crt -days 500 -sha256  -extfile <(printf "subjectAltName=DNS:client")

rm test_ca.key test_ca.srl test_client.csr test_server.csr

openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:512
openssl rsa -pubout -in private_key.pem -out public_key.pem