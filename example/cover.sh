curl http://localhost:8000/ping

# goc
goc build  --buildflags="-tags='goc'" . -o=main_goc ; ./main_goc

wget http://localhost:7777/cover.out
go tool cover -func=cover.out
go tool cover -html cover.out -o coverage.html

# go cov

go build  -cover -o app  .; mkdir cov ; GOCOVERDIR=cov ./app

go tool covdata textfmt -i=cov -o cover_report.txt
go tool cover -func=cover_report.txt
go tool cover -html=cover_report.txt -o cover_report.html

