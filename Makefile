.PHONY: all

all:
	go build && ./baseweb

count:
	fd | grep .go$ | xargs wc -l
