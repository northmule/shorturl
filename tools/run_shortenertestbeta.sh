# shellcheck disable=SC2164
cd ../cmd/shortener
go build -buildvcs=false -o shortener
chmod +x shortener
cd ../
clear
shortenertestbeta -test.v -test.run=^TestIteration -binary-path=/home/djo/GolandProjects/shorturl/cmd/shortener/shortener
