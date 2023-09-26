package middlerwares

import (
	"math/rand"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/oklog/ulid"
)

func RequestIDMiddleware(c *gin.Context) {
	reqID := genUniqueID()
	c.Set("req_id", reqID)
	c.Header("X-Request-ID", reqID)
}

// Use pool to avoid concurrent access for rand.Source
var entropyPool = sync.Pool{
	New: func() interface{} {
		return rand.New(rand.NewSource(time.Now().UnixNano()))
	},
}

// Generate Unique ID
// Currently using ULID, this maybe conflict with other process with very low possibility
func genUniqueID() string {
	entropy := entropyPool.Get().(*rand.Rand)
	defer entropyPool.Put(entropy)
	id := ulid.MustNew(ulid.Now(), entropy)
	return id.String()
}
