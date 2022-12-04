package bumper

//go:generate go run github.com/abice/go-enum -f $GOFILE --marshal --names

// Field ENUM(auto, major, minor, patch, prerelease)
type Field string
