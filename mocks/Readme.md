To generate mocks, run the following command:

```cassandraql
go get github.com/vektra/mockery/.../
cd sdk
mockery -all -recursive -keeptree
```

If will generate `mocks` in `sdk` folder. After that you can copy it to `mocks` folder.