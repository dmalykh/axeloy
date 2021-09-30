build-ways:
	DIR="$$(pwd)/src/ways"; \
	for name in `ls -1 $$DIR`; do (cd "$$DIR/$$name" && bash -c "go build -buildmode=plugin -o $$name.so $$name.go"); done
mocks:
	 (cd src && mockery --all --keeptree)
 migration:
	 migrate create -ext sql -dir $$(pwd)/src/repository/db/_migrations -seq $$@
