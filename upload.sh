curl -X POST -F "files=@readme.md"  -F "file_path=zero"  'http://localhost:8000/index/upload'
curl -X POST -F "files=@uploadPage.go"  -F "file_path=zero"  'http://localhost:8000/index/upload'