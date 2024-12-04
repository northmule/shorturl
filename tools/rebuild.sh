echo $SHORTURL_PROJECT_PATH
cd $SHORTURL_PROJECT_PATH/cmd/shortener
go build -buildvcs=false -o shortener
chmod +x shortener

cd $SHORTURL_PROJECT_PATH/cmd/staticlint
go build -buildvcs=false -o staticlint
chmod +x staticlint

cd $SHORTURL_PROJECT_PATH/cmd/shortener_grpc
go build -buildvcs=false -o shortener_grpc
chmod +x shortener_grpc