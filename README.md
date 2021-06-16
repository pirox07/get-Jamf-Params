# get_Jamf-Params

## Required

- [Go install](https://golang.org/doc/install)
- Jamf Pro Account
    - Privileges ... `Policies:READ`, `Scripts: READ`

## Usage

```
export JAMF_BASE_URL=https://<your tenant>.jamfcloud.com
export JAMF_USER=<Jamf Pro Username>
export JAMF_USER_PASSWORD=<Jamf Pro User Password>

go get
go run main.go
```