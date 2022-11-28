package utils

import (
	"google.golang.org/grpc/metadata"
)

func firstMetadataWithName(md metadata.MD, name string) string {
	values := md.Get(name)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
