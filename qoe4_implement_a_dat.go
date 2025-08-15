package main

import (
	"log"
	"net/http"
	"encoding/json"
	"fmt"
	"time"
	"strings"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

// Define structure for notification data
type Notification struct {
	Title       string `json:"title"`
	Message     string `json:"message"`
	Sender      string `json:"sender"`
	Receiver    string `json:"receiver"`
	Timestamp   int64  `json:"timestamp"`
}

// Define structure for mobile device data
type MobileDevice struct {
	DeviceID  string `json:"device_id"`
	Platform  string `json:"platform"`
}

func main() {
	// Initialize Firebase Firestore
	ctx := context.Background()
	sa := option-WithCredentialsFile("path/to/serviceAccountKey.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error initializing Firestore: %v", err)
	}
	defer client.Close()

	// Define HTTP handler for receiving notification data
	http.HandleFunc("/notifications", func(w http.ResponseWriter, r *http.Request) {
		var notification Notification
		err := json.NewDecoder(r.Body).Decode(&notification)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Store notification data in Firestore
		_, _, err = client.Collection("notifications").Add(ctx, notification)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Get list of mobile devices to notify
		devices, err := getMobileDevices(client, ctx, notification.Receiver)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send notification to each mobile device
		for _, device := range devices {
			sendNotification(device, notification)
		}
	})

	http.ListenAndServe(":8080", nil)
}

func getMobileDevices(client *firestore.Client, ctx context.Context, receiver string) ([]MobileDevice, error) {
	var devices []MobileDevice
	iter := client.Collection("mobile_devices").Where("receiver", "==", receiver).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var device MobileDevice
		err = doc.DataTo(&device)
		if err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}
	return devices, nil
}

func sendNotification(device MobileDevice, notification Notification) {
	// Implement notification sending logic based on device platform
	if device.Platform == "iOS" {
		// Send notification using APNs
		fmt.Println("Sending iOS notification...")
	} else if device.Platform == "Android" {
		// Send notification using FCM
		fmt.Println("Sending Android notification...")
	}
}