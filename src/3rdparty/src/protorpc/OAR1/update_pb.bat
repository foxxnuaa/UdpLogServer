
@rem gen cxx code
--protoc.exe --cxx_out=. msg.proto
--protoc.exe --cxx_out=. server_push.proto


@rem gen go code
protoc.exe --go_out=. msg.proto
protoc.exe --go_out=. server_push.proto
protoc.exe --go_out=. addressbook.proto

