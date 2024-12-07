build:
	mkdir -p ./bin
	go build -o ./bin/signservice cmd/main.go


run: build
	chmod +x ./bin/signservice
	lsof -ti :8080 | xargs kill -9
	./bin/signservice


run-mock:
	sleep 1 # Wait for the service to start
	curl --location 'http://localhost:8080/api/v0/signature-device' \
		--header 'Content-Type: application/json' \
		--data '{"id":"1","algorithm":"ECC"}' &

	curl --location 'http://localhost:8080/api/v0/signature-device' \
		--header 'Content-Type: application/json' \
		--data '{"id":"2","algorithm":"RSA"}' &

	curl --location 'http://localhost:8080/api/v0/signing-creation' \
		--header 'Content-Type: application/json' \
		--data '{"device_id":"1","data":"test 1 with alg: 1 "}' &

	curl --location 'http://localhost:8080/api/v0/signing-creation' \
		--header 'Content-Type: application/json' \
		--data '{"device_id":"1","data":"test 2 with alg: 1 "}' &

	curl --location 'http://localhost:8080/api/v0/signing-creation' \
		--header 'Content-Type: application/json' \
		--data '{"device_id":"1","data":"test 3 with alg: 1 "}' &

	curl --location 'http://localhost:8080/api/v0/signing-creation' \
		--header 'Content-Type: application/json' \
		--data '{"device_id":"1","data":"test 4 with alg: 1 "}' &

	curl --location 'http://localhost:8080/api/v0/signing-creation' \
		--header 'Content-Type: application/json' \
		--data '{"device_id":"2","data":"test 1 with alg: 2"}' &

	curl --location 'http://localhost:8080/api/v0/signing-creation' \
		--header 'Content-Type: application/json' \
		--data '{"device_id":"2","data":"test 2 with alg: 2"}' &

	curl --location 'http://localhost:8080/api/v0/signing-creation' \
		--header 'Content-Type: application/json' \
		--data '{"device_id":"2","data":"test 3 with alg: 2"}' &

	curl --location 'http://localhost:8080/api/v0/signing-creation' \
		--header 'Content-Type: application/json' \
		--data '{"device_id":"2","data":"test 4 with alg: 2 "}' &
	wait # Wait for all background tasks to complete
