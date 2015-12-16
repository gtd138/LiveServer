protoc --go_out=. *.proto

copy *.go ..\..\..\src\msg_proto
del *.go

pause...