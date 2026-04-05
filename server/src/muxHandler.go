package src

import (
	"fmt"
	"net/http"

	"github.com/AatirNadim/getMe/server/store"
	"github.com/AatirNadim/getMe/server/utils"
)

func muxHandler(storeInstance *store.Store) *http.ServeMux {
	mux := http.NewServeMux()

	c := &Controllers{
		StoreInstance: storeInstance,
	}

	mux.HandleFunc(fmt.Sprintf("GET %s", utils.GetRoute), c.GetController())

	mux.HandleFunc(fmt.Sprintf("POST %s", utils.PutRoute), c.PutController())

	mux.HandleFunc(fmt.Sprintf("DELETE %s", utils.DeleteRoute), c.DeleteController())

	mux.HandleFunc(fmt.Sprintf("DELETE %s", utils.ClearStoreRoute), c.ClearStoreController())

	mux.HandleFunc(fmt.Sprintf("POST %s", utils.BatchPutRoute), c.BatchPutController())

	mux.HandleFunc(fmt.Sprintf("POST %s", utils.BatchGetRoute), c.BatchGetController())

	mux.HandleFunc(fmt.Sprintf("POST %s", utils.BatchDeleteRoute), c.BatchDeleteController())
	return mux
}
