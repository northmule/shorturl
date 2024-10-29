echo $SHORTURL_PROJECT_PATH
cd $SHORTURL_PROJECT_PATH/cmd/shortener
go build -buildvcs=false -o shortener
chmod +x shortener

cd $SHORTURL_PROJECT_PATH/cmd/staticlint
go build -buildvcs=false -o staticlint
chmod +x staticlint