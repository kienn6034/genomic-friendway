clean-subnet:
	./scripts/clean-subnet.sh

run-subnet:
	./scripts/run-subnet.sh

deploy-contract:
	./scripts/deploy-contract.sh

generate-artifacts:
	./scripts/generate-artifacts.sh


test-contracts:
	cd genomicdao && npx hardhat test

test-backend:
	cd genomic-service && go test ./...

test: test-contracts test-backend