mongo:
	docker run -d --network host -v ${HOME}/temp_mongo:/data/db mvertes/alpine-mongo:4.0.6-1

api:
	go run -mod=vendor -race github.com/holive/feedado/app/cmd/api

worker:
	go run -mod=vendor -race github.com/holive/feedado/app/cmd/worker

run: mongo api worker
