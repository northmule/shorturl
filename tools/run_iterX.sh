clear
echo $SHORTURL_PROJECT_PATH
cd $SHORTURL_PROJECT_PATH/cmd/shortener
go build -buildvcs=false -o shortener
chmod +x shortener
cd $SHORTURL_PROJECT_PATH

# Запуск теста
shortenertestbeta -test.v -test.run=^TestIteration15$ -binary-path=$SHORTURL_PROJECT_PATH/cmd/shortener/shortener -database-dsn='postgres://postgres:123@localhost:5456/shorturl?sslmode=disable'
