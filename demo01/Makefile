default:run

env:
	export https_proxy=http://127.0.0.1:8118;export http_proxy=http://127.0.0.1:8118

	wget https://raw.githubusercontent.com/ossrs/srs/develop/trunk/doc/source.200kbps.768x320.flv

run:
	CGO_LDFLAGS=-L/usr/local/ffmpeg/lib CGO_CFLAGS=-I/usr/local/ffmpeg/include go run demo01.go