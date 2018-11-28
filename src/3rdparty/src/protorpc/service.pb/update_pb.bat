
@rem gen cxx code
protoc.exe --cxx_out=. arith.proto


@rem gen go code
protoc.exe --go_out=. arith.proto

