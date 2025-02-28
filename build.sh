#!/bin/sh

go build -o ggcr main.go
sudo mv ggcr /bin/ggcr
rm *.mp3
