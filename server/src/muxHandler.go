package src

import (
	"fmt"
	"getMeMod/server/store"
	"getMeMod/server/utils"
	"net/http"
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
