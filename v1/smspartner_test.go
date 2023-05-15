package smspartner_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/hoflish/smspartner-go/v1"
)

func TestCredits(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		b, err := fixture("credits.json")
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprint(w, string(b))
	})

	cli, teardown := testingHTTPClient(t, h)
	defer teardown()

	res, err := cli.CheckCredits()
	if err != nil {
		t.Fatal(err)
	}

	gotUsername := res.User.Username
	wantUsernme := "example@gmail.com"
	if gotUsername != wantUsernme {
		t.Errorf("got: '%s', want: '%s'", gotUsername, wantUsernme)
	}
}

func TestSendSMS(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		b, err := fixture("send_sms.json")
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprint(w, string(b))
	})

	cli, teardown := testingHTTPClient(t, h)
	defer teardown()

	d := smspartner.NewDate(2018, 8, 16, 17, 45)
	minute, err := d.MinuteToSendSMS()
	if err != nil {
		t.Error(err)
	}
	sms := &smspartner.SMS{
		PhoneNumbers: "0620123456",
		Message:      "Your message goes here",
		Gamme:        smspartner.LowCost,
		ScheduledDeliveryDate: d.ScheduledDeliveryDate(),
		Time:   d.Time.Hour(),
		Minute: minute,
	}

	res, err := cli.SendSMS(sms)
	if err != nil {
		t.Fatal(err)
	}

	wantNbOfSMS := 1
	wantCurrency := "EUR"

	if res.NumberOfSMS != wantNbOfSMS {
		t.Errorf("got: %d, want: %d", res.NumberOfSMS, 1)
	}

	if res.Currency != wantCurrency {
		t.Errorf("got: %s, want: %s", res.Currency, wantCurrency)
	}
}

func TestSendSMSWithError(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		b, err := fixture("send_sms_error.json")
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprint(w, string(b))
	})

	cli, teardown := testingHTTPClient(t, h)
	defer teardown()

	sms := &smspartner.SMS{}
	res, err := cli.SendSMS(sms)

	if res != nil {
		t.Errorf("response should be nil, but got: %#v", res)
	}

	wantErr := "Le message est requis (and 5 other errors)"
	if err.Error() != wantErr {
		t.Errorf("got: %s want: %s", err, wantErr)
	}
}

func TestSendBulkSMS(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		b, err := fixture("send_bulksms.json")
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprint(w, string(b))
	})

	cli, teardown := testingHTTPClient(t, h)
	defer teardown()

	d := smspartner.NewDate(2018, 8, 16, 18, 30)
	minute, err := d.MinuteToSendSMS()
	if err != nil {
		t.Error(err)
	}

	bulksms := &smspartner.BulkSMS{
		SMSList: []*smspartner.SMSPayload{
			{
				PhoneNumber: "0620123456",
				Message:     "Your message goes here",
			},
			{
				PhoneNumber: "0621123456",
				Message:     "Your message goes here",
			},
		},
		ScheduledDeliveryDate: d.ScheduledDeliveryDate(),
		Time:   d.Time.Hour(),
		Minute: minute,
	}

	res, err := cli.SendBulkSMS(bulksms)
	if err != nil {
		t.Fatal(err)
	}

	var totalCost float64
	for _, smsrl := range res.SMSResponseList {
		totalCost += smsrl.Cost
	}

	if res.Cost != totalCost {
		t.Errorf("got: %f, want: %f", res.Cost, totalCost)
	}
}

func TestSendVirtualNumber(t *testing.T) {
	t.Skip("Need to communicate with remote API team")
}

func TestVerifyNumber(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		b, err := fixture("verify_number.json")
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprint(w, string(b))
	})

	cli, teardown := testingHTTPClient(t, h)
	defer teardown()

	reqPayload := &smspartner.NumberVerificationRequest{
		PhoneNumbers: "+212620123456,+212621123456",
	}

	res, err := cli.VerifyNumber(reqPayload)
	if err != nil {
		t.Fatal(err)
	}

	gotCost := res.Cost
	wantCost := 0.005

	if gotCost != wantCost {
		t.Errorf("got: %f, want: %f", gotCost, wantCost)
	}
}

func TestVerifyNumberFormat(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		b, err := fixture("lookup.json")
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprint(w, string(b))
	})

	cli, teardown := testingHTTPClient(t, h)
	defer teardown()

	phoneNumbers := []string{"+212620123456", "+212666123456"}
	res, err := cli.VerifyNumberFormat(phoneNumbers...)
	if err != nil {
		t.Fatal(err)
	}

	gotLookupLength := len(res.Lookup) // number of phoneNumbers to verify
	if gotLookupLength != 2 {
		t.Errorf("got: %d, want: %d", gotLookupLength, 2)
	}

	for _, tt := range res.Lookup {
		if tt.Success == true {
			gotNetwork := tt.Network
			wantNetwork := "Méditel"

			if gotNetwork != wantNetwork {
				t.Errorf("got: %s, want: %s", gotNetwork, wantNetwork)
			}
		}
	}
}

func TestCancelSMS(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		b, err := fixture("cancel_sms.json")
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprint(w, string(b))
	})

	cli, teardown := testingHTTPClient(t, h)
	defer teardown()

	msgID := 2271595
	res, err := cli.CancelSMS(msgID)
	if err != nil {
		t.Fatal(err)
	}

	gotMessage := res["message"]
	wantMessage := "L'envoi du SMS a été annulé."

	if gotMessage != wantMessage {
		t.Errorf("got: %s, want: %s", gotMessage, wantMessage)
	}

}

func TestGetSMSStatus(t *testing.T) {
	t.Error("Not implemented yet")
}

func TestGetMultiSMSStatus(t *testing.T) {
	t.Error("Not implemented yet")
}

func TestGetBulkSMSStatus(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		b, err := fixture("bulk_status.json")
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprint(w, string(b))
	})

	cli, teardown := testingHTTPClient(t, h)
	defer teardown()

	msgID := 2270142
	res, err := cli.GetBulkSMSStatus(msgID)
	if err != nil {
		t.Fatal(err)
	}

	gotMessageId := res.MessageID
	wantMessageId := strconv.Itoa(msgID)

	if gotMessageId != wantMessageId {
		t.Errorf("got: %s, want: %s", gotMessageId, wantMessageId)
	}
}

func TestListStops(t *testing.T) {
	t.Error("Not implemented yet")
}

func TestAddToStops(t *testing.T) {
	t.Error("Not implemented yet")
}

func TestDeleteFromStops(t *testing.T) {
	t.Error("Not implemented yet")
}

func testingHTTPClient(t *testing.T, handler http.Handler) (*smspartner.Client, func()) {
	server := httptest.NewServer(handler)

	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, server.Listener.Addr().String())
			},
		},
	}

	testApiKey := smspartner.APIKey("TEST_API_KEY")

	spClient, err := smspartner.NewClient(cli, testApiKey)
	if err != nil {
		t.Fatalf("error creating client: %v", err)
	}

	return spClient, server.Close
}

func fixture(path string) ([]byte, error) {
	f, err := os.Open("testdata/" + path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return b, nil
}
