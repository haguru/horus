.PHONY: docker up down

docker:
	cd ./crumbdb_service; make docker; cd -
	cd ./user_acct_service; make docker; cd -
	cd ./follower_service; make docker; cd -

up: 
	docker compose up -d

down: 
	docker compose down