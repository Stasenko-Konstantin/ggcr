#!/bin/sh

go build main.go -o ggcr
sudo mv ggcr /bin/ggcr
rm *.mp3
