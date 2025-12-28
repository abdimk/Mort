build:
	@go build -o ./bin/Mort

run: build
	 @./bin/Mort

runfollower: build
	@./bin/Mort --listenaddr :4000 --leaderaddr :3000

# clean:
# 	rm -rf bin

# start:clean build run