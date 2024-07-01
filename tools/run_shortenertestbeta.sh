clear
echo $SHORTURL_PROJECT_PATH
cd $SHORTURL_PROJECT_PATH/cmd/shortener
go build -buildvcs=false -o shortener
chmod +x shortener
cd $SHORTURL_PROJECT_PATH
shortenertestbeta -test.v -test.run=^TestIteration -binary-path=$SHORTURL_PROJECT_PATH/cmd/shortener/shortener
