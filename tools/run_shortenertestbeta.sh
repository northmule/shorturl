clear
echo $SHORTURL_PROJECT_PATH
cd $SHORTURL_PROJECT_PATH/cmd/shortener
go build -buildvcs=false -o shortener
chmod +x shortener
cd $SHORTURL_PROJECT_PATH

# Запуск тестов

# Спринт 1
shortenertestbeta -test.v -test.run=^TestIteration1$ -binary-path=$SHORTURL_PROJECT_PATH/cmd/shortener/shortener
shortenertestbeta -test.v -test.run=^TestIteration2$ -source-path=$SHORTURL_PROJECT_PATH
shortenertestbeta -test.v -test.run=^TestIteration3$ -source-path=$SHORTURL_PROJECT_PATH
shortenertestbeta -test.v -test.run=^TestIteration4$ -binary-path=$SHORTURL_PROJECT_PATH/cmd/shortener/shortener -server-port=9880
shortenertestbeta -test.v -test.run=^TestIteration5$ -binary-path=$SHORTURL_PROJECT_PATH/cmd/shortener/shortener -server-port=9880

# Спринт 2
shortenertestbeta -test.v -test.run=^TestIteration6$ -source-path=$SHORTURL_PROJECT_PATH
shortenertestbeta -test.v -test.run=^TestIteration7$ -source-path=$SHORTURL_PROJECT_PATH -binary-path=$SHORTURL_PROJECT_PATH/cmd/shortener/shortener
shortenertestbeta -test.v -test.run=^TestIteration8$ -binary-path=$SHORTURL_PROJECT_PATH/cmd/shortener/shortener
shortenertestbeta -test.v -test.run=^TestIteration9$ -binary-path=$SHORTURL_PROJECT_PATH/cmd/shortener/shortener -source-path=$SHORTURL_PROJECT_PATH -file-storage-path="/tmp/iter9_short_autotest.json"
