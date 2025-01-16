package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
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

func webhook_MicrosoftGraph(w http.ResponseWriter, r *http.Request) {

	token := r.URL.Query().Get("validationToken")
	if token != "" {
		fmt.Println("Confirming subscription")
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
