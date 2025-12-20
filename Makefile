build:
	@go build -o ./bin/Mort

run: build
	 ./bin/Mort


# clean:
# 	rm -rf bin

# start:clean build run