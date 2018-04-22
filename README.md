# ocp-router-overrider

# mkdir $GOPATH/src/ocp-router-overrider

# vendoring
https://github.com/Masterminds/glide

```sh
go get github.com/Masterminds/glide
go install github.com/Masterminds/glide
```

```sh



 glide install --strip-vendor

cat <<EOF > glide.yaml
package: gettingstarted 
import:
- package: github.com/openshift/client-go
  version: release-3.9
EOF

```