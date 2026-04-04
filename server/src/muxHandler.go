package src

import (
	"fmt"
	"net/http"

	"github.com/AatirNadim/getMe/server/store"
	"github.com/AatirNadim/getMe/server/utils"
)

func muxHandler(storeInstance *store.Store) *http.ServeMux {
	mux := http.NewServeMux()

	c := &Controllers{}

	mux.HandleFunc(fmt.Sprintf("GET %s", utils.GetRoute), c.GetController(storeInstance))

	mux.HandleFunc(fmt.Sprintf("POST %s", utils.PutRoute), c.PutController(storeInstance))

	mux.HandleFunc(fmt.Sprintf("DELETE %s", utils.DeleteRoute), c.DeleteController(storeInstance))

	mux.HandleFunc(fmt.Sprintf("DELETE %s", utils.ClearStoreRoute), c.ClearStoreController(storeInstance))

	mux.HandleFunc(fmt.Sprintf("POST %s", utils.BatchPutRoute), c.BatchPutController(storeInstance))

	mux.HandleFunc(fmt.Sprintf("POST %s", utils.BatchGetRoute), c.BatchGetController(storeInstance))

	mux.HandleFunc(fmt.Sprintf("POST %s", utils.BatchDeleteRoute), c.BatchDeleteController(storeInstance))
	return mux
}
