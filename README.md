# protodoc
Protobuf to documentation generator. See [Cloudprober](https://github.com/cloudprober/cloudprober)'s [config documentation](https://cloudprober.org/docs/config/overview/) for an example.

### Example:

To generate documentation for Cloudprober config:
```
go run ./cmd/protodoc/. --proto_root_dir=<path_to_cloudprober_code> --package_prefix=github.com/cloudprober/cloudprober
```
