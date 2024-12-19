package fileupload

import (
	"fmt"
	"os"
	"path"

	"github.com/majid-cj/go-chat-server/util"
)

// FormatFile ...
func FormatFile(filepath string) string {
	ext := path.Ext(filepath)
	return fmt.Sprintf("%s-%s%s", util.GetTimeNow().Format(os.Getenv("TIME_FORMATE")), util.ULID(), ext)
}
