go test -coverprofile=cov.out
go tool cover -html=cov.out

go install rsc.io/uncover@latest
uncover cov.out
