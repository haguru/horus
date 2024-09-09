.PHONY: docker up down lint

docker:
	cd ./crumbdb_service; make docker; cd -
	cd ./user_acct_service; make docker; cd -
	cd ./follower_service; make docker; cd -

down: 
	docker compose down

fmt:
	cd ./crumbdb_service; make fmt; cd -
	cd ./user_acct_service; make fmt; cd -
	cd ./follower_service; make fmt; cd -

fixfmt:
	cd ./crumbdb_service; gofmt -l -w .; cd -
	cd ./user_acct_service; gofmt -l -w .; cd -
	cd ./follower_service; gofmt -l -w .; cd -
 
lint:
	cd ./crumbdb_service; make lint; cd -
	cd ./user_acct_service; make lint; cd -
	cd ./follower_service; make lint; cd -

up: 
	docker compose up -d



