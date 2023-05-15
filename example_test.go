package smspartner_test

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hoflish/smspartner-go/v1"
)

func Example_client_CheckCredits() {
	client, err := smspartner.NewClient(&http.Client{})
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.CheckCredits()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Credits - response: %#v\n", resp)
}

func Example_client_SendSMS() {
	client, err := smspartner.NewClient(&http.Client{})
	if err != nil {
		log.Fatal(err)
	}

	date := smspartner.NewDate(2018, 8, 16, 17, 45)
	minute, err := date.MinuteToSendSMS()
	if err != nil {
		log.Fatal(err)
	}
	sms := &smspartner.SMS{
		PhoneNumbers: "+212620123456, +212621123456",
		Message:      "This is your message",
		Gamme:        smspartner.LowCost,
		ScheduledDeliveryDate: date.ScheduledDeliveryDate(),
		Time:   date.Time.Hour(),
		Minute: minute,
	}

	resp, err := client.SendSMS(sms)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("SMS sending - response: %#v\n", resp)
}

func Example_client_SendBulkSMS() {
	client, err := smspartner.NewClient(&http.Client{})
	if err != nil {
		log.Fatal(err)
	}

	date := smspartner.NewDate(2018, 8, 16, 17, 45)
	minute, err := date.MinuteToSendSMS()
	if err != nil {
		log.Fatal(err)
	}

	bulksms := &smspartner.BulkSMS{
		SMSList: []*smspartner.SMSPayload{
			{
				PhoneNumber: "+212620xxxxxx",
				Message:     "This is your message",
			},
			{
				PhoneNumber: "+212620xxxxxx",
				Message:     "This is your message",
			},
		},
		Gamme: smspartner.Premium,
		ScheduledDeliveryDate: date.ScheduledDeliveryDate(),
		Time:   date.Time.Hour(),
		Minute: minute,
	}
	resp, err := client.SendBulkSMS(bulksms)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Bulk SMS sending - response: %#v\n", resp)
}

func Example_client_SendVirtualNumber() {
	client, err := smspartner.NewClient(&http.Client{})
	if err != nil {
		log.Fatal(err)
	}

	vn := &smspartner.VNumber{
		To:      "+212620123456",
		From:    "+212620123456",
		Message: "This is your message",
	}

	resp, err := client.SendVirtualNumber(vn)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("VN sending - response: %#v\n", resp)
}

func Example_client_VerifyNumber() {
	client, err := smspartner.NewClient(&http.Client{})
	if err != nil {
		log.Fatal(err)
	}

	payload := &smspartner.NumberVerificationRequest{
		PhoneNumbers: "+212620123456,+212621123456",
		NotifyURL:    "http://example.com/api/hlr/notify",
	}

	resp, err := client.VerifyNumber(payload)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Number verification - response: %#v\n", resp)
}

func Example_client_VerifyNumberFormat() {
	client, err := smspartner.NewClient(&http.Client{})
	if err != nil {
		log.Fatal(err)
	}

	phoneNumbers := []string{"+212620123456", "+212621123456"}
	resp, err := client.VerifyNumberFormat(phoneNumbers...)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Format Number verification - response: %#v\n", resp)
}

func Example_client_CancelSMS() {
	client, err := smspartner.NewClient(&http.Client{})
	if err != nil {
		log.Fatal(err)
	}

	messageID := 2274024
	resp, err := client.CancelSMS(messageID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("SMS cancelling - response: %#v\n", resp)
}

func Example_client_GetSMSStatus() {
	client, err := smspartner.NewClient(&http.Client{})
	if err != nil {
		log.Fatal(err)
	}

	messageID := 2274024
	phoneNumber := "+212620123456"
	resp, err := client.GetSMSStatus(messageID, phoneNumber)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("SMS status - response: %#v\n", resp)
}

func Example_client_GetMultiSMSStatus() {
	client, err := smspartner.NewClient(&http.Client{})
	if err != nil {
		log.Fatal(err)
	}
	ss := &smspartner.MultiSMSStatusReq{
		SMSStatusList: []*smspartner.MultiSMSStatusPayload{
			{
				PhoneNumber: "+212620123456",
				MessageID:   2270142,
			},
			{
				PhoneNumber: "+212621123456",
				MessageID:   2270142,
			},
		},
	}
	resp, err := client.GetMultiSMSStatus(ss)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Multi SMS status - response: %#v\n", resp)
}
