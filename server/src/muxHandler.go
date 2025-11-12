package src

import (
	"fmt"
	"net/http"

	"github.com/AatirNadim/getMe/server/store"
	"github.com/AatirNadim/getMe/server/utils"
)

func muxHandler(storeInstance *store.Store) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc(fmt.Sprintf("GET %s", utils.GetRoute), GetController(storeInstance))

	mux.HandleFunc(fmt.Sprintf("POST %s", utils.PutRoute), PutController(storeInstance))

	mux.HandleFunc(fmt.Sprintf("DELETE %s", utils.DeleteRoute), DeleteController(storeInstance))

	mux.HandleFunc(fmt.Sprintf("DELETE %s", utils.ClearStoreRoute), ClearStoreController(storeInstance))

	mux.HandleFunc(fmt.Sprintf("POST %s", utils.BatchPutRoute), BatchPutCController(storeInstance))
	return mux
}
