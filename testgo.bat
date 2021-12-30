cd D:\workspace\GO\Toy_Prj\basic_struct\cmd\web
go test -coverprofile=coverage.out && go tool cover -html=coverage.out
timeout 2 > NUL
cd ../../