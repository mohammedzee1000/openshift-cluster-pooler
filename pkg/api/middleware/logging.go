package middleware

import (
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/api/servercontext"
	"log"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		ctx, err := servercontext.NewAPIServerContext()
		if err != nil {
			log.Fatal("unable in instantiate context : ", err.Error())
		}
		ctx.Log.Info("api server", r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}