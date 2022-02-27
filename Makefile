.PHONY: buildways
buildways:
	DIR="$$(pwd)/src/ways"; \
	for name in `ls -1 $$DIR`; do (cd "$$DIR/$$name" && bash -c "go build -buildmode=plugin -o $$name.so $$name.go"); done


.PHONY: build
build:
	(make buildways && cd src && go mod tidy && go build main.go)

mocks:
	 (cd src && mockery --all --keeptree)

.PHONY: reformgen
reformgen:
	DIR="$$(pwd)/src/repository/db"; \
	for name in `ls -1 $$DIR`; do \
	  	if [ -d "$$DIR/$$name" ]; then \
			(cd "$$DIR/$$name" && bash -c "go generate model/*.go");\
		fi \
  	done
	#(cd src/repository/db && go generate */model/*.go )


 db.migration.new:
	 @echo "Hi $@" \
	# bash -c "migrate create -ext sql -dir $$(pwd)/src/repository/db/_migrations -seq $$@"
