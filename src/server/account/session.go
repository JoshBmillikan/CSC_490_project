package account

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"math/rand"
	"net/http"
	"sort"
	"sync"
)

type IdToken struct {
	username string
	uuid     uuid.UUID
}

var mutex sync.RWMutex
var tokens []IdToken

func SessionAuth(ctx *gin.Context) {
	cookie, err := ctx.Cookie("SHADER_VIEW_SESSION")
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
	var token IdToken
	err = json.Unmarshal([]byte(cookie), &token)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
	mutex.RLock()
	view := tokens
	found := false
	for true {
		if len(view) == 0 {
			if view[0].uuid == token.uuid {
				found = true
			}
			break
		}

		mid := len(view) / 2
		if view[mid].uuid == token.uuid {
			found = true
			break
		} else if view[mid].uuid.String() < token.uuid.String() {
			view = view[0:mid]
		} else {
			view = view[mid : len(view)-1]
		}
	}
	mutex.RUnlock()
	if !found {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	} else {
		ctx.Set("session", token)
	}
}

func newSession(username string) (IdToken, error) {
	rnd := rand.New(rand.NewSource(1))
	id, err := uuid.NewRandomFromReader(rnd)
	if err != nil {
		return IdToken{}, err
	}

	token := IdToken{
		username: username,
		uuid:     id,
	}
	go addToken(token)
	return token, nil
}

func addToken(token IdToken) {
	mutex.Lock()
	tokens = append(tokens, token)
	sort.Slice(tokens, func(i, j int) bool {
		return tokens[i].uuid.String() < tokens[j].uuid.String()
	})
	mutex.Unlock()
}