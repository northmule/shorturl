# shellcheck disable=SC2164
cd ../cmd/shortener
chmod +x shortener
go build -buildvcs=false -o shortener
cd ../
clear
shortenertestbeta -test.v -test.run=^TestIteration1$ -binary-path=cmd/shortener/shortener
