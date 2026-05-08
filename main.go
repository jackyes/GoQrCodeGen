package main

import (
	"log"
	"net/http"
)

const (
	QRSmall      = 128
	QRMedium     = 256
	QRLarge      = 512
	QRExtraLarge = 1024

	LogoPercent = 0.25
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", serveHTML)
	http.HandleFunc("/generate", generateQRCodeHandler)
	http.HandleFunc("/generate_instagram", generateInstagramQRCodeHandler)
	http.HandleFunc("/generate_facebook", generateFacebookQRCodeHandler)
	http.HandleFunc("/generate_tiktok", generateTikTokQRCodeHandler)
	http.HandleFunc("/generate_linkedin", generateLinkedInQRCodeHandler)
	http.HandleFunc("/generate_youtube", generateYouTubeQRCodeHandler)
	http.HandleFunc("/generate_vcard", generateVCardQRCodeHandler)
	http.HandleFunc("/generate_wifi", generateWiFiQRCodeHandler)
	http.HandleFunc("/generate_map", generateMapQRCodeHandler)
	http.HandleFunc("/generate_event", generateEventQRCodeHandler)
	http.HandleFunc("/generate_paypal", generatePayPalQRCodeHandler)
	http.HandleFunc("/generate_whatsapp", generateWhatsAppQRCodeHandler)
	http.HandleFunc("/generate_x", generateXQRCodeHandler)
	http.HandleFunc("/generate_email", generateEmailQRCodeHandler)
	http.HandleFunc("/generate_sms", generateSMSQRCodeHandler)
	http.HandleFunc("/generate_phone", generatePhoneQRCodeHandler)
	http.HandleFunc("/generate_spotify", generateSpotifyQRCodeHandler)
	http.HandleFunc("/generate_telegram", generateTelegramQRCodeHandler)
	http.HandleFunc("/generate_zoom", generateZoomQRCodeHandler)

	log.Println("Server running on port 5555")
	log.Fatal(http.ListenAndServe(":5555", nil))
}

func serveHTML(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}
