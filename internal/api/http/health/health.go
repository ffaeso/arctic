package health

import (
	"net/http"

	"github.com/fraeso/arctic/internal/api/response"
)

type PingResponse struct {
	Status string `json:"status"`
}

func Ping(w http.ResponseWriter, r *http.Request) {
	res := &PingResponse{
		Status: "up",
	}
	response.WriteJSON(w, http.StatusOK, res)
}
