package src

import (
	"fmt"
	"net/http"

	"github.com/AatirNadim/getMe/server/store"
	// "github.com/AatirNadim/getMe/server/utils"
	"github.com/AatirNadim/getMe/commons"
)

func muxHandler(storeInstance *store.Store) *http.ServeMux {
	mux := http.NewServeMux()

	c := &Controllers{
		StoreInstance: storeInstance,
	}

	mux.HandleFunc(fmt.Sprintf("GET %s", commons.GetRoute), c.GetController())

	mux.HandleFunc(fmt.Sprintf("POST %s", commons.PutRoute), c.PutController())

	mux.HandleFunc(fmt.Sprintf("DELETE %s", commons.DeleteRoute), c.DeleteController())

	mux.HandleFunc(fmt.Sprintf("DELETE %s", commons.ClearStoreRoute), c.ClearStoreController())

	mux.HandleFunc(fmt.Sprintf("POST %s", commons.BatchPutRoute), c.BatchPutController())

	mux.HandleFunc(fmt.Sprintf("POST %s", commons.BatchGetRoute), c.BatchGetController())

	mux.HandleFunc(fmt.Sprintf("DELETE %s", commons.BatchDeleteRoute), c.BatchDeleteController())
	return mux
}
