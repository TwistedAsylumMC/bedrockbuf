protobufs:
	protoc -I ./ -I ../ --go_out=./ --go_opt=paths=source_relative ./*.proto

force:
	rm -f *.pb.go
	make protobufs