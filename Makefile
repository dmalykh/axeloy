ways.build:
	DIR="$$(pwd)/src/ways"; \
	for name in `ls -1 $$DIR`; do (cd "$$DIR/$$name" && bash -c "go build -buildmode=plugin -o $$name.so $$name.go"); done
mocks:
	 (cd src && mockery --all --keeptree)

db.reform.gen:
	DIR="$$(pwd)/src/repository/db"; \
	for name in `ls -1 $$DIR`; do \
		(cd "$$DIR/$$name" && bash -c "go generate model/*.go");\
  	done
	#(cd src/repository/db && go generate */model/*.go )

 db.migration.new:
	 @echo "Hi $@" \
	# bash -c "migrate create -ext sql -dir $$(pwd)/src/repository/db/_migrations -seq $$@"
