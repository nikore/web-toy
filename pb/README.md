The helloworld.pb.go file is compiled from the protobuff schema. I am leaving the generated file checked 
in so the build dosen't have to setup the compiler and check it in. If you change the schema at all you
need to recompile the file to do that:

1. Download the compiler software here https://github.com/google/protobuf/releases newest version should
be fine and make sure you pick the one that works with the os. 
2. Run the following command: `protoc -I pb/ pb/helloworld.proto --go_out=plugins=grpc:pb` from project root dir

once you run the command the `helloworld.pb.go` will have been regenerated and any new messages schemas
should be available for use. You can adjust the outputs if you wish to generate other clients.