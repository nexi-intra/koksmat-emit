package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nexi-intra/koksmat-emit/internal/emitter"
	//"github.com/koksmat-com/koksmat/model"
	//"github.com/magicbutton/magic-mix/model"
)

const webhooksTag = "Webhooks"

type WebhookEventStruct struct {
	SubscriptionID                 string    `json:"subscriptionId"`
	SubscriptionExpirationDateTime time.Time `json:"subscriptionExpirationDateTime"`
	ChangeType                     string    `json:"changeType"`
	Resource                       string    `json:"resource"`
	ResourceData                   struct {
		OdataType string `json:"@odata.type"`
		OdataID   string `json:"@odata.id"`
		OdataEtag string `json:"@odata.etag"`
		ID        string `json:"id"`
	} `json:"resourceData"`
	ClientState string `json:"clientState"`
	TenantID    string `json:"tenantId"`
}
type Callback struct {
	Value []WebhookEventStruct `json:"value"`
}

// webhook_MicrosoftGraph handles incoming HTTP requests for Microsoft Graph webhooks.
// It performs validation of the subscription by checking for a "validationToken" query parameter.
// If the token is present, it confirms the subscription by echoing the token back to the client.
// If the token is not present, it decodes the request body into a Callback struct and processes the values.
// It responds with a 200 status code and a "received" message upon successful processing.
//
// Parameters:
//   - w: http.ResponseWriter to write the HTTP response.
//   - r: *http.Request containing the HTTP request.
//
// Responses:
//   - 200 OK: If the validation token is confirmed or the request body is successfully processed.
//   - 400 Bad Request: If there is an error decoding the request body.
func webhook_MicrosoftGraph(app *emitter.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token := r.URL.Query().Get("validationToken")
		if token != "" {
			app.Obs.Info("Confirming subscription")

			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.WriteHeader(200)
			fmt.Fprint(w, token)
			return
		}

		p := &Callback{}
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println(err)
			return

		}

		for _, v := range p.Value {
			fmt.Println(v)

		}

		w.WriteHeader(200)
		fmt.Fprint(w, "received")

	}
}

// func getWebHooks() usecase.Interactor {
// 	type GetRequest struct {
// 		//	Paging `bson:",inline"`
// 	}

// 	type GetResponse struct {
// 		Webhooks []*officegraph.MicrosoftGraphSubscription `json:"webhooks"`
// 		// NumberOfRecords int64                                     `json:"numberofrecords"`
// 		// Pages           int64                                     `json:"pages"`
// 		// CurrentPage     int64                                     `json:"currentpage"`
// 		// PageSize        int64                                     `json:"pagesize"`
// 	}
// 	u := usecase.NewInteractor(func(ctx context.Context, input GetRequest, output *GetResponse) error {

// 		data, err := officegraph.SubscriptionList()
// 		output.Webhooks = data
// 		// output.NumberOfRecords = int64(len(data))
// 		// output.Pages = 1
// 		// output.CurrentPage = 1
// 		// output.PageSize = 100

// 		return err

// 	})

// 	u.SetTitle("Get webhooks ")

// 	u.SetExpectedErrors(status.InvalidArgument)
// 	u.SetTags(
// 		webhooksTag,
// 	)
// 	return u
// }
