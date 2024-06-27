package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
	"github.com/skip2/go-qrcode"
)

const (
	// QR code size constants
	QRSmall      = 128  // Small QR code size in pixels
	QRMedium     = 256  // Medium QR code size in pixels
	QRLarge      = 512  // Large QR code size in pixels
	QRExtraLarge = 1024 // Extra large QR code size in pixels

	// Logo size configuration
	LogoPercent = 0.25 // Percentage of the QR code occupied by the logo

	// Logo image paths
	InstagramLogoPath = "static/instagram_logo.png"
	FacebookLogoPath  = "static/facebook_logo.png"
	TikTokLogoPath    = "static/tiktok_logo.png"
	LinkedInLogoPath  = "static/linkedin_logo.png"
	YouTubeLogoPath   = "static/youtube_logo.png"
	WiFiLogoPath      = "static/wifi_logo.png"
	MapLogoPath       = "static/map_logo.png"
	EventLogoPath     = "static/event_logo.png"
	PayPalLogoPath    = "static/paypal_logo.png"
	WhatsAppLogoPath  = "static/whatsapp_logo.png"
	XLogoPath         = "static/x_logo.png"
	EmailLogoPath     = "static/email_logo.png"
	SMSLogoPath       = "static/sms_logo.png"
	PhoneLogoPath     = "static/phone_logo.png"
	SpotifyLogoPath   = "static/spotify_logo.png"
	TelegramLogoPath  = "static/telegram_logo.png"
	ZoomLogoPath      = "static/zoom_logo.png"
)

func main() {
	// Serve static files from the "static" directory
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Define handler functions for different QR code generation requests
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

	// Log server startup message
	log.Println("Server running on port 5555")

	// Start the server and handle fatal errors
	log.Fatal(http.ListenAndServe(":5555", nil))
}

func serveHTML(w http.ResponseWriter, r *http.Request) {
	// Serve the index.html file
	http.ServeFile(w, r, "static/index.html")
}

func isValidQRCodeSize(size int) bool {
	// Check if the provided size is one of the predefined valid sizes
	return size == QRSmall || size == QRMedium || size == QRLarge || size == QRExtraLarge
}

func generateMapQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST, otherwise return an error
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("generateMapQRCodeHandler: Method not allowed")
		return
	}

	// Extract latitude and longitude from the request form
	latitude := r.FormValue("latitude")
	longitude := r.FormValue("longitude")
	// Check if both latitude and longitude are present in the request
	if latitude == "" || longitude == "" {
		http.Error(w, "Missing latitude or longitude", http.StatusBadRequest)
		log.Printf("generateMapQRCodeHandler: Missing latitude or longitude")
		return
	}

	// Validate the format and range of latitude
	lat, err := strconv.ParseFloat(latitude, 64)
	if err != nil || lat < -90 || lat > 90 {
		http.Error(w, "Invalid latitude", http.StatusBadRequest)
		log.Printf("generateMapQRCodeHandler: Invalid latitude")
		return
	}

	// Validate the format and range of longitude
	lon, err := strconv.ParseFloat(longitude, 64)
	if err != nil || lon < -180 || lon > 180 {
		http.Error(w, "Invalid longitude", http.StatusBadRequest)
		log.Printf("generateMapQRCodeHandler: Invalid longitude")
		return
	}

	// Extract the requested QR code size from the form
	sizeStr := r.FormValue("size")

	// Check if the size parameter is present
	if sizeStr == "" {
		http.Error(w, "Missing size", http.StatusBadRequest)
		log.Printf("generateMapQRCodeHandler: Missing size")
		return
	}

	// Convert the size string to an integer and validate it against allowed sizes
	size, err := strconv.Atoi(sizeStr)
	if err != nil || !isValidQRCodeSize(size) {
		http.Error(w, "Invalid size", http.StatusBadRequest)
		log.Printf("generateMapQRCodeHandler: Invalid size - %v", err)
		return
	}

	// Construct the geo URI for Google Maps using the validated latitude and longitude
	geoURL := fmt.Sprintf("geo:%f,%f", lat, lon)

	// Generate the QR code for the geo URI with the requested size
	qrCode, err := generateQRCode(geoURL, size)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		log.Printf("generateMapQRCodeHandler: Failed to generate QR code - %v", err)
		return
	}

	// Open the map logo image file
	mapLogoFile, err := http.Dir(".").Open(MapLogoPath)
	if err != nil {
		http.Error(w, "Failed to open map logo", http.StatusInternalServerError)
		log.Printf("generateMapQRCodeHandler: Failed to open map logo - %v", err)
		return
	}
	defer mapLogoFile.Close() // Close the file after processing

	// Decode the map logo image
	mapLogo, err := decodeImage(mapLogoFile)
	if err != nil {
		http.Error(w, "Failed to decode map logo", http.StatusInternalServerError)
		log.Printf("generateMapQRCodeHandler: Failed to decode map logo - %v", err)
		return
	}

	// Overlay the map logo onto the QR code with a specific logo size percentage
	qrCode, err = overlayImageOnQRCode(qrCode, mapLogo, LogoPercent)
	if err != nil {
		http.Error(w, "Failed to overlay map logo on QR code", http.StatusInternalServerError)
		log.Printf("generateMapQRCodeHandler: Failed to overlay map logo on QR code - %v", err)
		return
	}

	// Set the content type header to indicate PNG image data
	w.Header().Set("Content-Type", "image/png")

	// Encode the QR code image as PNG format and write it to the HTTP response writer
	err = png.Encode(w, qrCode)
	if err != nil {
		log.Printf("generateMapQRCodeHandler: Failed to encode QR code as PNG - %v", err)
	}
}

func generateWiFiQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST, otherwise return an error
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("generateWiFiQRCodeHandler: Method not allowed")
		return
	}

	// Extract SSID, password, security type, and size from the request form
	ssid := r.FormValue("ssid")
	password := r.FormValue("password")
	security := r.FormValue("security")
	sizeStr := r.FormValue("size")

	// Validate the presence of SSID
	if ssid == "" {
		http.Error(w, "Missing SSID", http.StatusBadRequest)
		log.Printf("generateWiFiQRCodeHandler: Missing SSID")
		return
	}

	// Define a map of valid security types for Wi-Fi networks
	validSecurities := map[string]bool{"WPA": true, "WPA2": true, "WPA3": true, "WEP": true, "nopass": true}

	// Validate the provided security typ
	if !validSecurities[security] {
		http.Error(w, "Invalid security type", http.StatusBadRequest)
		log.Printf("generateWiFiQRCodeHandler: Invalid security type")
		return
	}

	// Validate password requirements for WPA/WPA2/WPA3 security
	if security == "WPA" || security == "WPA2" || security == "WPA3" {
		if password == "" {
			http.Error(w, "Password is required for WPA/WPA2/WPA3 security", http.StatusBadRequest)
			log.Printf("generateWiFiQRCodeHandler: Password is required for WPA/WPA2/WPA3 security")
			return
		}
		if len(password) < 8 || len(password) > 63 {
			http.Error(w, "Password for WPA/WPA2/WPA3 must be between 8 and 63 characters", http.StatusBadRequest)
			log.Printf("generateWiFiQRCodeHandler: Password for WPA/WPA2/WPA3 must be between 8 and 63 characters")
			return
		}
	}

	// Validate password requirements for WEP security
	if security == "WEP" {
		if len(password) != 5 && len(password) != 13 {
			http.Error(w, "Password for WEP must be exactly 5 or 13 characters", http.StatusBadRequest)
			log.Printf("generateWiFiQRCodeHandler: Password for WEP must be exactly 5 or 13 characters")
			return
		}
	}

	// Validate the presence of size parameter
	if sizeStr == "" {
		http.Error(w, "Missing size", http.StatusBadRequest)
		log.Printf("generateWiFiQRCodeHandler: Missing size")
		return
	}

	// Convert size string to integer and validate it against allowed sizes
	size, err := strconv.Atoi(sizeStr)
	if err != nil || !isValidQRCodeSize(size) {
		http.Error(w, "Invalid size", http.StatusBadRequest)
		log.Printf("generateWiFiQRCodeHandler: Invalid size - %v", err)
		return
	}

	// Construct the Wi-Fi network information string using the validated parameters
	wifiString := fmt.Sprintf("WIFI:T:%s;S:%s;P:%s;;", security, ssid, password)
	qrCode, err := generateQRCode(wifiString, size)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		log.Printf("generateWiFiQRCodeHandler: Failed to generate QR code - %v", err)
		return
	}

	// Generate the QR code for the Wi-Fi network information string with the requested size
	wifiLogoFile, err := http.Dir(".").Open(WiFiLogoPath)
	if err != nil {
		http.Error(w, "Failed to open Wi-Fi logo", http.StatusInternalServerError)
		log.Printf("generateWiFiQRCodeHandler: Failed to open Wi-Fi logo - %v", err)
		return
	}
	defer wifiLogoFile.Close() // Close the file after processing

	// Decode the Wi-Fi logo image
	wifiLogo, err := decodeImage(wifiLogoFile)
	if err != nil {
		http.Error(w, "Failed to decode Wi-Fi logo", http.StatusInternalServerError)
		log.Printf("generateWiFiQRCodeHandler: Failed to decode Wi-Fi logo - %v", err)
		return
	}

	// Overlay the Wi-Fi logo onto the QR code with a specific logo size percentage
	qrCode, err = overlayImageOnQRCode(qrCode, wifiLogo, LogoPercent)
	if err != nil {
		http.Error(w, "Failed to overlay Wi-Fi logo on QR code", http.StatusInternalServerError)
		log.Printf("generateWiFiQRCodeHandler: Failed to overlay Wi-Fi logo on QR code - %v", err)
		return
	}

	// Set the content type header to indicate PNG image data
	w.Header().Set("Content-Type", "image/png")

	// Encode the QR code image as PNG format and write it to the HTTP response writer
	err = png.Encode(w, qrCode)
	if err != nil {
		log.Printf("generateWiFiQRCodeHandler: Failed to encode QR code as PNG - %v", err)
	}
}

func generateLinkedInQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST, otherwise return an error
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("generateLinkedInQRCodeHandler: Method not allowed")
		return
	}

	// Extract username from the request form
	username := r.FormValue("username")

	// Validate the presence of username
	if username == "" {
		http.Error(w, "Missing username", http.StatusBadRequest)
		log.Printf("generateLinkedInQRCodeHandler: Missing username")
		return
	}

	// Extract size string from the request form
	sizeStr := r.FormValue("size")
	// Validate the presence of size parameter
	if sizeStr == "" {
		http.Error(w, "Missing size", http.StatusBadRequest)
		log.Printf("generateLinkedInQRCodeHandler: Missing size")
		return
	}

	// Convert size string to integer and validate it against allowed sizes
	size, err := strconv.Atoi(sizeStr)
	if err != nil || !isValidQRCodeSize(size) {
		http.Error(w, "Invalid size", http.StatusBadRequest)
		log.Printf("generateLinkedInQRCodeHandler: Invalid size - %v", err)
		return
	}

	// Construct the LinkedIn profile URL using the extracted username
	url := "https://www.linkedin.com/in/" + username

	// Generate the QR code for the LinkedIn profile URL with the requested size
	qrCode, err := generateQRCode(url, size)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		log.Printf("generateLinkedInQRCodeHandler: Failed to generate QR code - %v", err)
		return
	}

	// Open the LinkedIn logo image file
	linkedinLogoFile, err := http.Dir(".").Open(LinkedInLogoPath)
	if err != nil {
		http.Error(w, "Failed to open LinkedIn logo", http.StatusInternalServerError)
		log.Printf("generateLinkedInQRCodeHandler: Failed to open LinkedIn logo - %v", err)
		return
	}
	defer linkedinLogoFile.Close() // Close the file after processing

	// Decode the LinkedIn logo image
	linkedinLogo, err := decodeImage(linkedinLogoFile)
	if err != nil {
		http.Error(w, "Failed to decode LinkedIn logo", http.StatusInternalServerError)
		log.Printf("generateLinkedInQRCodeHandler: Failed to decode LinkedIn logo - %v", err)
		return
	}

	// Overlay the LinkedIn logo onto the QR code with a specific logo size percentage
	qrCode, err = overlayImageOnQRCode(qrCode, linkedinLogo, LogoPercent)
	if err != nil {
		http.Error(w, "Failed to overlay LinkedIn logo on QR code", http.StatusInternalServerError)
		log.Printf("generateLinkedInQRCodeHandler: Failed to overlay LinkedIn logo on QR code - %v", err)
		return
	}

	// Set the content type header to indicate PNG image data
	w.Header().Set("Content-Type", "image/png")

	// Encode the QR code image as PNG format and write it to the HTTP response writer
	err = png.Encode(w, qrCode)
	if err != nil {
		log.Printf("generateLinkedInQRCodeHandler: Failed to encode QR code as PNG - %v", err)
	}
}

func generateYouTubeQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST, otherwise return an error
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("generateYouTubeQRCodeHandler: Method not allowed")
		return
	}
	// Extract YouTube channel name from the request form
	channel := r.FormValue("channel")
	// Validate the presence of channel name
	if channel == "" {
		http.Error(w, "Missing channel", http.StatusBadRequest)
		log.Printf("generateYouTubeQRCodeHandler: Missing channel")
		return
	}

	// Extract size string from the request form
	sizeStr := r.FormValue("size")
	// Validate the presence of size parameter
	if sizeStr == "" {
		http.Error(w, "Missing size", http.StatusBadRequest)
		log.Printf("generateYouTubeQRCodeHandler: Missing size")
		return
	}

	// Convert size string to integer and validate it against allowed sizes
	size, err := strconv.Atoi(sizeStr)
	if err != nil || !isValidQRCodeSize(size) {
		http.Error(w, "Invalid size", http.StatusBadRequest)
		log.Printf("generateYouTubeQRCodeHandler: Invalid size - %v", err)
		return
	}

	// Construct the YouTube channel URL using the extracted channel name
	url := "https://www.youtube.com/channel/" + channel

	// Generate the QR code for the YouTube channel URL with the requested size
	qrCode, err := generateQRCode(url, size)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		log.Printf("generateYouTubeQRCodeHandler: Failed to generate QR code - %v", err)
		return
	}
	// Open the YouTube logo image file
	youtubeLogoFile, err := http.Dir(".").Open(YouTubeLogoPath)
	if err != nil {
		http.Error(w, "Failed to open YouTube logo", http.StatusInternalServerError)
		log.Printf("generateYouTubeQRCodeHandler: Failed to open YouTube logo - %v", err)
		return
	}
	defer youtubeLogoFile.Close() // Close the file after processing

	// Decode the YouTube logo image
	youtubeLogo, err := decodeImage(youtubeLogoFile)
	if err != nil {
		http.Error(w, "Failed to decode YouTube logo", http.StatusInternalServerError)
		log.Printf("generateYouTubeQRCodeHandler: Failed to decode YouTube logo - %v", err)
		return
	}

	// Overlay the YouTube logo onto the QR code with a specific logo size percentage
	qrCode, err = overlayImageOnQRCode(qrCode, youtubeLogo, LogoPercent)
	if err != nil {
		http.Error(w, "Failed to overlay YouTube logo on QR code", http.StatusInternalServerError)
		log.Printf("generateYouTubeQRCodeHandler: Failed to overlay YouTube logo on QR code - %v", err)
		return
	}
	// Set the content type header to indicate PNG image data
	w.Header().Set("Content-Type", "image/png")
	// Encode the QR code image as PNG format and write it to the HTTP response writer
	err = png.Encode(w, qrCode)
	if err != nil {
		log.Printf("generateYouTubeQRCodeHandler: Failed to encode QR code as PNG - %v", err)
	}
}

func generateQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST, otherwise return an error
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("generateQRCodeHandler: Method not allowed")
		return
	}

	// Extract URL from the request form
	url := r.FormValue("url")

	// Validate the presence of URL
	if url == "" {
		http.Error(w, "Missing URL", http.StatusBadRequest)
		log.Printf("generateQRCodeHandler: Missing URL")
		return
	}

	// Check if an image file was uploaded
	file, _, err := r.FormFile("image")

	// Handle errors except for missing file (handled separately)
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "Error reading image", http.StatusInternalServerError)
		log.Printf("generateQRCodeHandler: Error reading image - %v", err)
		return
	}
	// Extract size string from the request form
	sizeStr := r.FormValue("size")
	// Validate the presence of size parameter
	if sizeStr == "" {
		http.Error(w, "Missing size", http.StatusBadRequest)
		log.Printf("generateQRCodeHandler: Missing size")
		return
	}

	// Convert size string to integer and validate it against allowed sizes
	size, err := strconv.Atoi(sizeStr)
	if err != nil || !isValidQRCodeSize(size) {
		http.Error(w, "Invalid size", http.StatusBadRequest)
		log.Printf("generateQRCodeHandler: Invalid size - %v", err)
		return
	}

	// Extract logo width percentage string (optional)
	logoWidthPercentStr := r.FormValue("logoWidthPercent")
	// Extract logo opacity string (optional, defaults to 1 if missing)
	logoOpacityStr := r.FormValue("logoOpacity")

	logoWidthPercent, err := strconv.ParseFloat(logoWidthPercentStr, 64)
	if err != nil {
		http.Error(w, "Invalid logo width percent", http.StatusBadRequest)
		log.Printf("generateQRCodeHandler: Invalid logo width percent - %v", err)
		return
	}

	// Parse logo opacity as float64 (handle potential parsing error with default value)
	logoOpacity, err := strconv.ParseFloat(logoOpacityStr, 64)
	if err != nil {
		logoOpacity = 1 // Use default opacity of 1 if parsing fails
	}

	// Generate the QR code for the provided URL with the requested size
	qrCode, err := generateQRCode(url, size)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		log.Printf("generateQRCodeHandler: Failed to generate QR code - %v", err)
		return
	}
	// If an image file was uploaded, process it
	if file != nil {
		// Decode the uploaded image
		overlayImage, err := decodeImage(file)
		if err != nil {
			http.Error(w, "Failed to decode image", http.StatusInternalServerError)
			log.Printf("generateQRCodeHandler: Failed to decode image - %v", err)
			return
		}

		// Overlay the uploaded image onto the QR code with specified width percentage and opacity
		qrCode, err = overlayImageOnQRCodeWithOpacity(qrCode, overlayImage, logoWidthPercent, logoOpacity)
		if err != nil {
			http.Error(w, "Failed to overlay image on QR code", http.StatusInternalServerError)
			log.Printf("generateQRCodeHandler: Failed to overlay image on QR code - %v", err)
			return
		}
	}

	// Set the content type header to indicate PNG image data
	w.Header().Set("Content-Type", "image/png")

	// Encode the QR code image as PNG format and write it to the HTTP response writer
	err = png.Encode(w, qrCode)
	if err != nil {
		log.Printf("generateQRCodeHandler: Failed to encode QR code as PNG - %v", err)
	}
}

func generateFacebookQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST, otherwise return an error
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("generateFacebookQRCodeHandler: Method not allowed")
		return
	}
	// Extract Facebook username from the request form
	username := r.FormValue("username")

	// Validate the presence of username
	if username == "" {
		http.Error(w, "Missing username", http.StatusBadRequest)
		log.Printf("generateFacebookQRCodeHandler: Missing username")
		return
	}

	// Extract size string from the request form
	sizeStr := r.FormValue("size")

	// Validate the presence of size parameter
	if sizeStr == "" {
		http.Error(w, "Missing size", http.StatusBadRequest)
		log.Printf("generateFacebookQRCodeHandler: Missing size")
		return
	}

	// Convert size string to integer and validate it against allowed sizes
	size, err := strconv.Atoi(sizeStr)
	if err != nil || !isValidQRCodeSize(size) {
		http.Error(w, "Invalid size", http.StatusBadRequest)
		log.Printf("generateFacebookQRCodeHandler: Invalid size - %v", err)
		return
	}

	// Construct the Facebook profile URL using the extracted username
	url := "https://www.facebook.com/" + username

	// Generate the QR code for the Facebook profile URL with the requested size
	qrCode, err := generateQRCode(url, size)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		log.Printf("generateFacebookQRCodeHandler: Failed to generate QR code - %v", err)
		return
	}

	// Open the Facebook logo image file
	facebookLogoFile, err := http.Dir(".").Open(FacebookLogoPath)
	if err != nil {
		http.Error(w, "Failed to open Facebook logo", http.StatusInternalServerError)
		log.Printf("generateFacebookQRCodeHandler: Failed to open Facebook logo - %v", err)
		return
	}
	defer facebookLogoFile.Close() // Close the file after processing

	// Decode the Facebook logo image
	facebookLogo, err := decodeImage(facebookLogoFile)
	if err != nil {
		http.Error(w, "Failed to decode Facebook logo", http.StatusInternalServerError)
		log.Printf("generateFacebookQRCodeHandler: Failed to decode Facebook logo - %v", err)
		return
	}

	// Overlay the Facebook logo onto the QR code with a specific logo size percentage
	qrCode, err = overlayImageOnQRCode(qrCode, facebookLogo, LogoPercent)
	if err != nil {
		http.Error(w, "Failed to overlay Facebook logo on QR code", http.StatusInternalServerError)
		log.Printf("generateFacebookQRCodeHandler: Failed to overlay Facebook logo on QR code - %v", err)
		return
	}

	// Set the content type header to indicate PNG image data
	w.Header().Set("Content-Type", "image/png")

	// Encode the QR code image as PNG format and write it to the HTTP response writer
	err = png.Encode(w, qrCode)
	if err != nil {
		log.Printf("generateFacebookQRCodeHandler: Failed to encode QR code as PNG - %v", err)
	}
}

func generateTikTokQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST, otherwise return an error
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("generateTikTokQRCodeHandler: Method not allowed")
		return
	}

	// Extract TikTok username from the request form
	username := r.FormValue("username")

	// Validate the presence of username
	if username == "" {
		http.Error(w, "Missing username", http.StatusBadRequest)
		log.Printf("generateTikTokQRCodeHandler: Missing username")
		return
	}

	// Extract size string from the request form
	sizeStr := r.FormValue("size")

	// Validate the presence of size parameter
	if sizeStr == "" {
		http.Error(w, "Missing size", http.StatusBadRequest)
		log.Printf("generateTikTokQRCodeHandler: Missing size")
		return
	}

	// Convert size string to integer and validate it against allowed sizes
	size, err := strconv.Atoi(sizeStr)
	if err != nil || !isValidQRCodeSize(size) {
		http.Error(w, "Invalid size", http.StatusBadRequest)
		log.Printf("generateTikTokQRCodeHandler: Invalid size - %v", err)
		return
	}

	// Construct the TikTok profile URL using the extracted username
	url := "https://www.tiktok.com/@" + username

	// Generate the QR code for the TikTok profile URL with the requested size
	qrCode, err := generateQRCode(url, size)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		log.Printf("generateTikTokQRCodeHandler: Failed to generate QR code - %v", err)
		return
	}

	// Open the TikTok logo image file
	tiktokLogoFile, err := http.Dir(".").Open(TikTokLogoPath)
	if err != nil {
		http.Error(w, "Failed to open TikTok logo", http.StatusInternalServerError)
		log.Printf("generateTikTokQRCodeHandler: Failed to open TikTok logo - %v", err)
		return
	}
	defer tiktokLogoFile.Close() // Close the file after processing

	// Decode the TikTok logo image
	tiktokLogo, err := decodeImage(tiktokLogoFile)
	if err != nil {
		http.Error(w, "Failed to decode TikTok logo", http.StatusInternalServerError)
		log.Printf("generateTikTokQRCodeHandler: Failed to decode TikTok logo - %v", err)
		return
	}

	// Overlay the TikTok logo onto the QR code with a specific logo size percentage
	qrCode, err = overlayImageOnQRCode(qrCode, tiktokLogo, LogoPercent)
	if err != nil {
		http.Error(w, "Failed to overlay TikTok logo on QR code", http.StatusInternalServerError)
		log.Printf("generateTikTokQRCodeHandler: Failed to overlay TikTok logo on QR code - %v", err)
		return
	}

	// Set the content type header to indicate PNG image data
	w.Header().Set("Content-Type", "image/png")

	// Encode the QR code image as PNG format and write it to the HTTP response writer
	err = png.Encode(w, qrCode)
	if err != nil {
		log.Printf("generateTikTokQRCodeHandler: Failed to encode QR code as PNG - %v", err)
	}
}

func generateInstagramQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST, otherwise return an error
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Extract Instagram username from the request form
	username := r.FormValue("username")

	// Validate the presence of username
	if username == "" {
		http.Error(w, "Missing username", http.StatusBadRequest)
		log.Printf("generateInstagramQRCodeHandler: Missing username")
		return
	}

	// Extract size string from the request form
	sizeStr := r.FormValue("size")

	// Validate the presence of size parameter
	if sizeStr == "" {
		http.Error(w, "Missing size", http.StatusBadRequest)
		log.Printf("generateInstagramQRCodeHandler: Missing size")
		return
	}

	// Convert size string to integer and validate it against allowed sizes
	size, err := strconv.Atoi(sizeStr)
	if err != nil || !isValidQRCodeSize(size) {
		http.Error(w, "Invalid size", http.StatusBadRequest)
		log.Printf("generateInstagramQRCodeHandler: Invalid size - %v", err)
		return
	}

	// Construct the Instagram profile URL using the extracted username
	url := "https://www.instagram.com/" + username

	// Generate the QR code for the Instagram profile URL with the requested size
	qrCode, err := generateQRCode(url, size)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		log.Printf("generateInstagramQRCodeHandler: Failed to generate QR code - %v", err)
		return
	}

	// Open the Instagram logo image file
	instagramLogoFile, err := http.Dir(".").Open(InstagramLogoPath)
	if err != nil {
		http.Error(w, "Failed to open Instagram logo", http.StatusInternalServerError)
		log.Printf("generateInstagramQRCodeHandler: Failed to open Instagram logo - %v", err)
		return
	}
	defer instagramLogoFile.Close() // Close the file after processing

	// Decode the Instagram logo image
	instagramLogo, err := decodeImage(instagramLogoFile)
	if err != nil {
		http.Error(w, "Failed to decode Instagram logo", http.StatusInternalServerError)
		log.Printf("generateInstagramQRCodeHandler: Failed to decode Instagram logo - %v", err)
		return
	}

	// Overlay the Instagram logo onto the QR code with a specific logo size percentage
	qrCode, err = overlayImageOnQRCode(qrCode, instagramLogo, LogoPercent)
	if err != nil {
		http.Error(w, "Failed to overlay Instagram logo on QR code", http.StatusInternalServerError)
		log.Printf("generateInstagramQRCodeHandler: Failed to overlay Instagram logo on QR code - %v", err)
		return
	}

	// Set the content type header to indicate PNG image data
	w.Header().Set("Content-Type", "image/png")

	// Encode the QR code image as PNG format and write it to the HTTP response writer
	err = png.Encode(w, qrCode)
	if err != nil {
		log.Printf("generateInstagramQRCodeHandler: Failed to encode QR code as PNG - %v", err)
	}
}

func generateVCardQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST, otherwise return an error
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract contact information from the request form
	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	title := r.FormValue("title")
	phone := r.FormValue("phone")
	mobile := r.FormValue("mobile")
	email := r.FormValue("email")
	address := r.FormValue("address")
	company := r.FormValue("company")
	url := r.FormValue("url")
	role := r.FormValue("role")
	lang := r.FormValue("lang")
	geo := r.FormValue("geo")

	// Generate a VCARD string representation of the contact information
	vCard := generateVCardString(firstName, lastName, title, phone, mobile, email, address, company, url, role, lang, geo)

	// Extract size string from the request form
	sizeStr := r.FormValue("size")

	// Validate the presence of size parameter
	if sizeStr == "" {
		http.Error(w, "Missing size", http.StatusBadRequest)
		log.Printf("generateVCardQRCodeHandler: Missing size")
		return
	}

	// Convert size string to integer and validate it against allowed sizes
	size, err := strconv.Atoi(sizeStr)
	if err != nil || !isValidQRCodeSize(size) {
		http.Error(w, "Invalid size", http.StatusBadRequest)
		log.Printf("generateVCardQRCodeHandler: Invalid size - %v", err)
		return
	}

	// Extract logo width percentage and opacity (optional) from the request form
	logoWidthPercentStr := r.FormValue("logoWidthPercent")
	logoOpacityStr := r.FormValue("logoOpacity")

	// Parse logo width percentage as float64
	logoWidthPercent, err := strconv.ParseFloat(logoWidthPercentStr, 64)
	if err != nil {
		http.Error(w, "Invalid logo width percent", http.StatusBadRequest)
		log.Printf("generateVCardQRCodeHandler: Invalid logo width percent - %v", err)
		return
	}

	// Parse logo opacity as float64 (handle potential parsing error with default value)
	logoOpacity, err := strconv.ParseFloat(logoOpacityStr, 64)
	if err != nil {
		logoOpacity = 1 // Use default opacity of 1 if parsing fails
	}

	// Generate the QR code for the VCARD data with the requested size
	qrCode, err := generateQRCode(vCard, size)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		log.Printf("generateVCardQRCodeHandler: Failed to generate QR code - %v", err)
		return
	}

	// Extract uploaded image file (optional)
	file, _, err := r.FormFile("image")

	// Handle errors except for missing file (handled separately)
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "Error reading image", http.StatusInternalServerError)
		log.Printf("generateVCardQRCodeHandler: Error reading image - %v", err)
		return
	}

	// If an image file was uploaded, process it
	if file != nil {
		// Decode the uploaded image
		overlayImage, err := decodeImage(file)
		if err != nil {
			http.Error(w, "Failed to decode image", http.StatusInternalServerError)
			log.Printf("generateVCardQRCodeHandler: Failed to decode image - %v", err)
			return
		}

		// Overlay the uploaded image onto the QR code with specified width percentage and opacity
		qrCode, err = overlayImageOnQRCodeWithOpacity(qrCode, overlayImage, logoWidthPercent, logoOpacity)
		if err != nil {
			http.Error(w, "Failed to overlay image on QR code", http.StatusInternalServerError)
			log.Printf("generateVCardQRCodeHandler: Failed to overlay image on QR code - %v", err)
			return
		}
	}

	// Set the content type header to indicate PNG image data
	w.Header().Set("Content-Type", "image/png")

	// Encode the QR code image as PNG format and write it to the HTTP response writer
	err = png.Encode(w, qrCode)
	if err != nil {
		log.Printf("generateVCardQRCodeHandler: Failed to encode QR code as PNG - %v", err)
	}
}

// Generate a QR code image from the given data string, with a specified size.
func generateQRCode(data string, size int) (image.Image, error) {
	// Create a new QR code instance with the given data and high error correction level.
	qr, err := qrcode.New(data, qrcode.High)
	if err != nil {
		// If there's an error creating the QR code, return it immediately.
		return nil, err
	}
	// Return the generated QR code image with the specified size.
	return qr.Image(size), nil
}

// Decode an image from a file reader, returning the image and any error.
func decodeImage(file io.Reader) (image.Image, error) {
	// Read the entire file into memory.
	imgData, err := io.ReadAll(file)
	if err != nil {
		// If there's an error reading the file, return it immediately.
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	// Attempt to decode the image using the standard image.Decode function.
	img, format, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		// If decoding fails, try again using format-specific decoders.
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Perform format-specific decoding if necessary.
	switch format {
	case "jpeg":
		// Decode JPEG images using the jpeg package.
		img, err = jpeg.Decode(bytes.NewReader(imgData))
	case "png":
		// Decode PNG images using the png package.
		img, err = png.Decode(bytes.NewReader(imgData))
	}

	if err != nil {
		// If decoding still fails, return the error.
		return nil, fmt.Errorf("failed to decode image after format detection: %w", err)
	}
	// Return the successfully decoded image.
	return img, nil
}

// Overlay an image on top of a QR code, returning the resulting image.
func overlayImageOnQRCode(qrCode image.Image, overlay image.Image, overlayPercent float64) (image.Image, error) {
	// Get the bounds of the QR code image.
	qrBounds := qrCode.Bounds()
	qrWidth := qrBounds.Dx()
	qrHeight := qrBounds.Dy()

	// Calculate the maximum size for the overlay image, based on the QR code size and the specified percentage.
	overlayMaxWidth := int(float64(qrWidth) * overlayPercent)
	overlayMaxHeight := int(float64(qrHeight) * overlayPercent)

	// Resize the overlay image to fit within the calculated maximum size, maintaining its aspect ratio.
	overlay = resize.Thumbnail(uint(overlayMaxWidth), uint(overlayMaxHeight), overlay, resize.Lanczos3)

	// Calculate the offset to center the overlay image on top of the QR code.
	offset := image.Pt((qrWidth-overlay.Bounds().Dx())/2, (qrHeight-overlay.Bounds().Dy())/2)

	// Create a new image with the same bounds as the QR code.
	b := qrBounds
	m := image.NewRGBA(b)

	// Draw the QR code onto the new image.
	draw.Draw(m, qrBounds, qrCode, image.Point{}, draw.Src)

	// Draw the overlay image on top of the QR code, centered and resized.
	draw.Draw(m, overlay.Bounds().Add(offset), overlay, image.Point{}, draw.Over)

	// Return the resulting image with the overlay.
	return m, nil
}

// Generate a vCard string from the given information.
func generateVCardString(firstName, lastName, title, phone, mobile, email, address, company, url, role, lang, geo string) string {
	// Create a string builder to efficiently build the vCard string.
	var sb strings.Builder
	sb.WriteString("BEGIN:VCARD\n")
	sb.WriteString("VERSION:3.0\n")

	// Add the formatted name (Last Name, First Name).
	sb.WriteString(fmt.Sprintf("N:%s;%s;;;\n", lastName, firstName))

	// Add the full name (First Name Last Name).
	sb.WriteString(fmt.Sprintf("FN:%s %s\n", firstName, lastName))

	// Add the company name if provided.
	if company != "" {
		sb.WriteString(fmt.Sprintf("ORG:%s\n", company))
	}

	// Add the title.
	sb.WriteString(fmt.Sprintf("TITLE:%s\n", title))

	// Add the work phone number.
	sb.WriteString(fmt.Sprintf("TEL;TYPE=WORK,VOICE:%s\n", phone))

	// Add the mobile phone number if provided.
	if mobile != "" {
		sb.WriteString(fmt.Sprintf("TEL;TYPE=CELL,VOICE:%s\n", mobile))
	}

	// Add the email address.
	sb.WriteString(fmt.Sprintf("EMAIL:%s\n", email))

	// Add the address.
	sb.WriteString(fmt.Sprintf("ADR:%s\n", address))

	// Add the URL if provided.
	if url != "" {
		sb.WriteString(fmt.Sprintf("URL:%s\n", url))
	}

	// Add the role if provided.
	if role != "" {
		sb.WriteString(fmt.Sprintf("ROLE:%s\n", role))
	}

	// Add the language if provided.
	if lang != "" {
		sb.WriteString(fmt.Sprintf("LANG:%s\n", lang))
	}

	// Add the geographical position if provided.
	if geo != "" {
		sb.WriteString(fmt.Sprintf("GEO:%s\n", geo))
	}

	// End the vCard.
	sb.WriteString("END:VCARD")
	return sb.String()
}

func generateEventQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	// Check for allowed method (POST only)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("generateEventQRCodeHandler: Method not allowed")
		return
	}

	// Extract event details from request form
	eventName := r.FormValue("eventName")
	startDateTime := r.FormValue("startDateTime")
	endDateTime := r.FormValue("endDateTime")
	location := r.FormValue("location")
	description := r.FormValue("description")

	// Validate presence of required event details
	if eventName == "" || startDateTime == "" || endDateTime == "" {
		http.Error(w, "Missing event details", http.StatusBadRequest)
		log.Printf("generateEventQRCodeHandler: Missing event details")
		return
	}

	// Extract and validate QR code size
	sizeStr := r.FormValue("size")
	if sizeStr == "" {
		http.Error(w, "Missing size", http.StatusBadRequest)
		log.Printf("generateEventQRCodeHandler: Missing size")
		return
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil || !isValidQRCodeSize(size) {
		http.Error(w, "Invalid size", http.StatusBadRequest)
		log.Printf("generateEventQRCodeHandler: Invalid size - %v", err)
		return
	}

	// Generate ICS string for event data
	icsString := fmt.Sprintf("BEGIN:VEVENT\nSUMMARY:%s\nDTSTART:%s\nDTEND:%s\nLOCATION:%s\nDESCRIPTION:%s\nEND:VEVENT",
		eventName, startDateTime, endDateTime, location, description)

	// Generate QR code from ICS string
	qrCode, err := generateQRCode(icsString, size)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		log.Printf("generateEventQRCodeHandler: Failed to generate QR code - %v", err)
		return
	}

	// Open event logo file
	eventLogoFile, err := http.Dir(".").Open(EventLogoPath)
	if err != nil {
		http.Error(w, "Failed to open event logo", http.StatusInternalServerError)
		log.Printf("generateEventQRCodeHandler: Failed to open event logo - %v", err)
		return
	}
	defer eventLogoFile.Close()

	// Decode event logo image
	eventLogo, err := decodeImage(eventLogoFile)
	if err != nil {
		http.Error(w, "Failed to decode event logo", http.StatusInternalServerError)
		log.Printf("generateEventQRCodeHandler: Failed to decode event logo - %v", err)
		return
	}

	// Overlay event logo on QR code
	qrCode, err = overlayImageOnQRCode(qrCode, eventLogo, LogoPercent)
	if err != nil {
		http.Error(w, "Failed to overlay event logo on QR code", http.StatusInternalServerError)
		log.Printf("generateEventQRCodeHandler: Failed to overlay event logo on QR code - %v", err)
		return
	}

	// Set content type for QR code image
	w.Header().Set("Content-Type", "image/png")

	// Encode QR code as PNG and write to response
	err = png.Encode(w, qrCode)
	if err != nil {
		log.Printf("generateEventQRCodeHandler: Failed to encode QR code as PNG - %v", err)
	}
}

// Generate a PayPal QrCode from the given information.
func generatePayPalQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	// Check for allowed method (POST only)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("generatePayPalQRCodeHandler: Method not allowed")
		return
	}

	// Extract payment details from request form
	email := r.FormValue("email")
	amount := r.FormValue("amount")
	currency := r.FormValue("currency")
	description := r.FormValue("description")

	// Validate presence of required payment details
	if email == "" || amount == "" || currency == "" {
		http.Error(w, "Missing payment details", http.StatusBadRequest)
		log.Printf("generatePayPalQRCodeHandler: Missing payment details")
		return
	}

	// Extract and validate QR code size
	sizeStr := r.FormValue("size")
	if sizeStr == "" {
		http.Error(w, "Missing size", http.StatusBadRequest)
		log.Printf("generatePayPalQRCodeHandler: Missing size")
		return
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil || !isValidQRCodeSize(size) {
		http.Error(w, "Invalid size", http.StatusBadRequest)
		log.Printf("generatePayPalQRCodeHandler: Invalid size - %v", err)
		return
	}

	// Generate PayPal payment URL
	paypalURL := fmt.Sprintf("https://www.paypal.com/cgi-bin/webscr?cmd=_xclick&business=%s&amount=%s&currency_code=%s&item_name=%s",
		email, amount, currency, description)

	// Generate QR code from PayPal URL
	qrCode, err := generateQRCode(paypalURL, size)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		log.Printf("generatePayPalQRCodeHandler: Failed to generate QR code - %v", err)
		return
	}

	// Open PayPal logo file
	paypalLogoFile, err := http.Dir(".").Open(PayPalLogoPath)
	if err != nil {
		http.Error(w, "Failed to open PayPal logo", http.StatusInternalServerError)
		log.Printf("generatePayPalQRCodeHandler: Failed to open PayPal logo - %v", err)
		return
	}
	defer paypalLogoFile.Close()

	// Decode PayPal logo image
	paypalLogo, err := decodeImage(paypalLogoFile)
	if err != nil {
		http.Error(w, "Failed to decode PayPal logo", http.StatusInternalServerError)
		log.Printf("generatePayPalQRCodeHandler: Failed to decode PayPal logo - %v", err)
		return
	}

	// Overlay PayPal logo on QR code
	qrCode, err = overlayImageOnQRCode(qrCode, paypalLogo, LogoPercent)
	if err != nil {
		http.Error(w, "Failed to overlay PayPal logo on QR code", http.StatusInternalServerError)
		log.Printf("generatePayPalQRCodeHandler: Failed to overlay PayPal logo on QR code - %v", err)
		return
	}

	// Set content type for QR code image
	w.Header().Set("Content-Type", "image/png")

	// Encode QR code as PNG and write to response
	err = png.Encode(w, qrCode)
	if err != nil {
		log.Printf("generatePayPalQRCodeHandler: Failed to encode QR code as PNG - %v", err)
	}
}

// Generates a QR code for opening a WhatsApp chat with a phone number and optional message.

func generateWhatsAppQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	// Check for allowed method (POST only)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("generateWhatsAppQRCodeHandler: Method not allowed")
		return
	}

	// Extract phone number and message from request form
	phone := r.FormValue("phone")
	message := r.FormValue("message")

	// Validate presence of phone number
	if phone == "" {
		http.Error(w, "Missing phone number", http.StatusBadRequest)
		log.Printf("generateWhatsAppQRCodeHandler: Missing phone number")
		return
	}

	// Extract and validate QR code size
	sizeStr := r.FormValue("size")
	if sizeStr == "" {
		http.Error(w, "Missing size", http.StatusBadRequest)
		log.Printf("generateWhatsAppQRCodeHandler: Missing size")
		return
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil || !isValidQRCodeSize(size) {
		http.Error(w, "Invalid size", http.StatusBadRequest)
		log.Printf("generateWhatsAppQRCodeHandler: Invalid size - %v", err)
		return
	}

	// Generate WhatsApp URL with phone number and message
	whatsappURL := fmt.Sprintf("https://wa.me/%s?text=%s", phone, message)

	// Generate QR code from WhatsApp URL
	qrCode, err := generateQRCode(whatsappURL, size)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		log.Printf("generateWhatsAppQRCodeHandler: Failed to generate QR code - %v", err)
		return
	}

	// Open WhatsApp logo file
	whatsappLogoFile, err := http.Dir(".").Open(WhatsAppLogoPath)
	if err != nil {
		http.Error(w, "Failed to open WhatsApp logo", http.StatusInternalServerError)
		log.Printf("generateWhatsAppQRCodeHandler: Failed to open WhatsApp logo - %v", err)
		return
	}
	defer whatsappLogoFile.Close()

	// Decode WhatsApp logo image
	whatsappLogo, err := decodeImage(whatsappLogoFile)
	if err != nil {
		http.Error(w, "Failed to decode WhatsApp logo", http.StatusInternalServerError)
		log.Printf("generateWhatsAppQRCodeHandler: Failed to decode WhatsApp logo - %v", err)
		return
	}

	// Overlay WhatsApp logo on QR code
	qrCode, err = overlayImageOnQRCode(qrCode, whatsappLogo, LogoPercent)
	if err != nil {
		http.Error(w, "Failed to overlay WhatsApp logo on QR code", http.StatusInternalServerError)
		log.Printf("generateWhatsAppQRCodeHandler: Failed to overlay WhatsApp logo on QR code - %v", err)
		return
	}

	// Set content type for QR code image
	w.Header().Set("Content-Type", "image/png")

	// Encode QR code as PNG and write to response
	err = png.Encode(w, qrCode)
	if err != nil {
		log.Printf("generateWhatsAppQRCodeHandler: Failed to encode QR code as PNG - %v", err)
	}
}

// Generates a QR code for X based on a username.
func generateXQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	// Check for allowed method (POST only)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("generateXQRCodeHandler: Method not allowed")
		return
	}

	// Extract username from request form
	username := r.FormValue("username")

	// Validate presence of username
	if username == "" {
		http.Error(w, "Missing username", http.StatusBadRequest)
		log.Printf("generateXQRCodeHandler: Missing username")
		return
	}

	// Extract and validate QR code size
	sizeStr := r.FormValue("size")
	if sizeStr == "" {
		http.Error(w, "Missing size", http.StatusBadRequest)
		log.Printf("generateXQRCodeHandler: Missing size")
		return
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil || !isValidQRCodeSize(size) {
		http.Error(w, "Invalid size", http.StatusBadRequest)
		log.Printf("generateXQRCodeHandler: Invalid size - %v", err)
		return
	}

	// Generate platform URL with username
	url := "https://www.twitter.com/" + username

	// Generate QR code from platform URL
	qrCode, err := generateQRCode(url, size)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		log.Printf("generateXQRCodeHandler: Failed to generate QR code - %v", err)
		return
	}

	// Open platform logo file
	xLogoFile, err := http.Dir(".").Open(XLogoPath)
	if err != nil {
		http.Error(w, "Failed to open X logo", http.StatusInternalServerError)
		log.Printf("generateXQRCodeHandler: Failed to open X logo - %v", err)
		return
	}
	defer xLogoFile.Close()

	// Decode platform logo image
	xLogo, err := decodeImage(xLogoFile)
	if err != nil {
		http.Error(w, "Failed to decode X logo", http.StatusInternalServerError)
		log.Printf("generateXQRCodeHandler: Failed to decode X logo - %v", err)
		return
	}

	// Overlay platform logo on QR code
	qrCode, err = overlayImageOnQRCode(qrCode, xLogo, LogoPercent)
	if err != nil {
		http.Error(w, "Failed to overlay X logo on QR code", http.StatusInternalServerError)
		log.Printf("generateXQRCodeHandler: Failed to overlay X logo on QR code - %v", err)
		return
	}

	// Set content type for QR code image
	w.Header().Set("Content-Type", "image/png")

	// Encode QR code as PNG and write to response
	err = png.Encode(w, qrCode)
	if err != nil {
		log.Printf("generateXQRCodeHandler: Failed to encode QR code as PNG - %v", err)
	}
}

// Generates a QR code for composing an email with a specific email address, subject, and body.
func generateEmailQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	// Check for allowed method (POST only)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("generateEmailQRCodeHandler: Method not allowed")
		return
	}

	// Extract email address, subject, and body from request form
	email := r.FormValue("email")
	subject := r.FormValue("subject")
	body := r.FormValue("body")

	// Validate presence of email address
	if email == "" {
		http.Error(w, "Missing email", http.StatusBadRequest)
		log.Printf("generateEmailQRCodeHandler: Missing email")
		return
	}

	// Extract and validate QR code size
	sizeStr := r.FormValue("size")
	if sizeStr == "" {
		http.Error(w, "Missing size", http.StatusBadRequest)
		log.Printf("generateEmailQRCodeHandler: Missing size")
		return
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil || !isValidQRCodeSize(size) {
		http.Error(w, "Invalid size", http.StatusBadRequest)
		log.Printf("generateEmailQRCodeHandler: Invalid size - %v", err)
		return
	}

	// Generatemailto URL with email address, subject, and body
	mailtoURL := fmt.Sprintf("mailto:%s?subject=%s&body=%s", email, subject, body)

	// Generate QR code from mailto URL
	qrCode, err := generateQRCode(mailtoURL, size)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		log.Printf("generateEmailQRCodeHandler: Failed to generate QR code - %v", err)
		return
	}

	// Open email logo file
	emailLogoFile, err := http.Dir(".").Open(EmailLogoPath)
	if err != nil {
		http.Error(w, "Failed to open email logo", http.StatusInternalServerError)
		log.Printf("generateEmailQRCodeHandler: Failed to open email logo - %v", err)
		return
	}
	defer emailLogoFile.Close()

	// Decode email logo image
	emailLogo, err := decodeImage(emailLogoFile)
	if err != nil {
		http.Error(w, "Failed to decode email logo", http.StatusInternalServerError)
		log.Printf("generateEmailQRCodeHandler: Failed to decode email logo - %v", err)
		return
	}

	// Overlay email logo on QR code
	qrCode, err = overlayImageOnQRCode(qrCode, emailLogo, LogoPercent)
	if err != nil {
		http.Error(w, "Failed to overlay email logo on QR code", http.StatusInternalServerError)
		log.Printf("generateEmailQRCodeHandler: Failed to overlay email logo on QR code - %v", err)
		return
	}

	// Set content type for QR code image
	w.Header().Set("Content-Type", "image/png")

	// Encode QR code as PNG and write to response
	err = png.Encode(w, qrCode)
	if err != nil {
		log.Printf("generateEmailQRCodeHandler: Failed to encode QR code as PNG - %v", err)
	}
}

// Generates a QR code for sending an SMS message to a phone number with an optional message.
func generateSMSQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	// Check for allowed method (POST only)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("generateSMSQRCodeHandler: Method not allowed")
		return
	}

	// Extract phone number and message from request form
	phoneNumber := r.FormValue("phoneNumber")
	message := r.FormValue("message")

	// Validate presence of phone number
	if phoneNumber == "" {
		http.Error(w, "Missing phone number", http.StatusBadRequest)
		log.Printf("generateSMSQRCodeHandler: Missing phone number")
		return
	}

	// Extract and validate QR code size
	sizeStr := r.FormValue("size")
	if sizeStr == "" {
		http.Error(w, "Missing size", http.StatusBadRequest)
		log.Printf("generateSMSQRCodeHandler: Missing size")
		return
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil || !isValidQRCodeSize(size) {
		http.Error(w, "Invalid size", http.StatusBadRequest)
		log.Printf("generateSMSQRCodeHandler: Invalid size - %v", err)
		return
	}

	// Generate SMS URL with phone number and message
	smsURL := fmt.Sprintf("sms:%s?body=%s", phoneNumber, message)

	// Generate QR code from SMS URL
	qrCode, err := generateQRCode(smsURL, size)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		log.Printf("generateSMSQRCodeHandler: Failed to generate QR code - %v", err)
		return
	}

	// Open SMS logo file
	smsLogoFile, err := http.Dir(".").Open(SMSLogoPath)
	if err != nil {
		http.Error(w, "Failed to open SMS logo", http.StatusInternalServerError)
		log.Printf("generateSMSQRCodeHandler: Failed to open SMS logo - %v", err)
		return
	}
	defer smsLogoFile.Close()

	// Decode SMS logo image
	smsLogo, err := decodeImage(smsLogoFile)
	if err != nil {
		http.Error(w, "Failed to decode SMS logo", http.StatusInternalServerError)
		log.Printf("generateSMSQRCodeHandler: Failed to decode SMS logo - %v", err)
		return
	}

	// Overlay SMS logo on QR code
	qrCode, err = overlayImageOnQRCode(qrCode, smsLogo, LogoPercent)
	if err != nil {
		http.Error(w, "Failed to overlay SMS logo on QR code", http.StatusInternalServerError)
		log.Printf("generateSMSQRCodeHandler: Failed to overlay SMS logo on QR code - %v", err)
		return
	}

	// Set content type for QR code image
	w.Header().Set("Content-Type", "image/png")

	// Encode QR code as PNG and write to response
	err = png.Encode(w, qrCode)
	if err != nil {
		log.Printf("generateSMSQRCodeHandler: Failed to encode QR code as PNG - %v", err)
	}
}

// Generates a QR code for calling a phone number.
func generatePhoneQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	// Check for allowed method (POST only)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("generatePhoneQRCodeHandler: Method not allowed")
		return
	}

	// Extract phone number from request form
	phoneNumber := r.FormValue("phoneNumber")

	// Validate presence of phone number
	if phoneNumber == "" {
		http.Error(w, "Missing phone number", http.StatusBadRequest)
		log.Printf("generatePhoneQRCodeHandler: Missing phone number")
		return
	}

	// Extract and validate QR code size
	sizeStr := r.FormValue("size")
	if sizeStr == "" {
		http.Error(w, "Missing size", http.StatusBadRequest)
		log.Printf("generatePhoneQRCodeHandler: Missing size")
		return
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil || !isValidQRCodeSize(size) {
		http.Error(w, "Invalid size", http.StatusBadRequest)
		log.Printf("generatePhoneQRCodeHandler: Invalid size - %v", err)
		return
	}

	// Generate phone URL with phone number
	phoneURL := fmt.Sprintf("tel:%s", phoneNumber)

	// Generate QR code from phone URL
	qrCode, err := generateQRCode(phoneURL, size)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		log.Printf("generatePhoneQRCodeHandler: Failed to generate QR code - %v", err)
		return
	}

	// Open phone logo file
	phoneLogoFile, err := http.Dir(".").Open(PhoneLogoPath)
	if err != nil {
		http.Error(w, "Failed to open phone logo", http.StatusInternalServerError)
		log.Printf("generatePhoneQRCodeHandler: Failed to open phone logo - %v", err)
		return
	}
	defer phoneLogoFile.Close()

	// Decode phone logo image
	phoneLogo, err := decodeImage(phoneLogoFile)
	if err != nil {
		http.Error(w, "Failed to decode phone logo", http.StatusInternalServerError)
		log.Printf("generatePhoneQRCodeHandler: Failed to decode phone logo - %v", err)
		return
	}

	// Overlay phone logo on QR code
	qrCode, err = overlayImageOnQRCode(qrCode, phoneLogo, LogoPercent)
	if err != nil {
		http.Error(w, "Failed to overlay phone logo on QR code", http.StatusInternalServerError)
		log.Printf("generatePhoneQRCodeHandler: Failed to overlay phone logo on QR code - %v", err)
		return
	}

	// Set content type for QR code image
	w.Header().Set("Content-Type", "image/png")

	// Encode QR code as PNG and write to response
	err = png.Encode(w, qrCode)
	if err != nil {
		log.Printf("generatePhoneQRCodeHandler: Failed to encode QR code as PNG - %v", err)
	}
}

// Generates a QR code for a Spotify URL.
func generateSpotifyQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	// Check for allowed method (POST only)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("generateSpotifyQRCodeHandler: Method not allowed")
		return
	}

	// Extract Spotify URL from request form
	spotifyURL := r.FormValue("spotifyURL")

	// Validate presence of Spotify URL
	if spotifyURL == "" {
		http.Error(w, "Missing Spotify URL", http.StatusBadRequest)
		log.Printf("generateSpotifyQRCodeHandler: Missing Spotify URL")
		return
	}

	// Extract and validate QR code size
	sizeStr := r.FormValue("size")
	if sizeStr == "" {
		http.Error(w, "Missing size", http.StatusBadRequest)
		log.Printf("generateSpotifyQRCodeHandler: Missing size")
		return
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil || !isValidQRCodeSize(size) {
		http.Error(w, "Invalid size", http.StatusBadRequest)
		log.Printf("generateSpotifyQRCodeHandler: Invalid size - %v", err)
		return
	}

	// Generate QR code from Spotify URL
	qrCode, err := generateQRCode(spotifyURL, size)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		log.Printf("generateSpotifyQRCodeHandler: Failed to generate QR code - %v", err)
		return
	}

	// Open Spotify logo file
	spotifyLogoFile, err := http.Dir(".").Open(SpotifyLogoPath)
	if err != nil {
		http.Error(w, "Failed to open Spotify logo", http.StatusInternalServerError)
		log.Printf("generateSpotifyQRCodeHandler: Failed to open Spotify logo - %v", err)
		return
	}
	defer spotifyLogoFile.Close()

	// Decode Spotify logo image
	spotifyLogo, err := decodeImage(spotifyLogoFile)
	if err != nil {
		http.Error(w, "Failed to decode Spotify logo", http.StatusInternalServerError)
		log.Printf("generateSpotifyQRCodeHandler: Failed to decode Spotify logo - %v", err)
		return
	}

	// Overlay Spotify logo on QR code
	qrCode, err = overlayImageOnQRCode(qrCode, spotifyLogo, LogoPercent)
	if err != nil {
		http.Error(w, "Failed to overlay Spotify logo on QR code", http.StatusInternalServerError)
		log.Printf("generateSpotifyQRCodeHandler: Failed to overlay Spotify logo on QR code - %v", err)
		return
	}

	// Set content type for QR code image
	w.Header().Set("Content-Type", "image/png")

	// Encode QR code as PNG and write to response
	err = png.Encode(w, qrCode)
	if err != nil {
		log.Printf("generateSpotifyQRCodeHandler: Failed to encode QR code as PNG - %v", err)
	}
}

// Generates a QR code for a Telegram.
func generateTelegramQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	// Check for allowed method (POST only)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("generateTelegramQRCodeHandler: Method not allowed")
		return
	}

	// Extract Telegram username or group name from request form
	telegramName := r.FormValue("telegramName")

	// Validate presence of Telegram name
	if telegramName == "" {
		http.Error(w, "Missing Telegram username or group name", http.StatusBadRequest)
		log.Printf("generateTelegramQRCodeHandler: Missing Telegram username or group name")
		return
	}

	// Extract and validate QR code size
	sizeStr := r.FormValue("size")
	if sizeStr == "" {
		http.Error(w, "Missing size", http.StatusBadRequest)
		log.Printf("generateTelegramQRCodeHandler: Missing size")
		return
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil || !isValidQRCodeSize(size) {
		http.Error(w, "Invalid size", http.StatusBadRequest)
		log.Printf("generateTelegramQRCodeHandler: Invalid size - %v", err)
		return
	}

	// Generate Telegram URL with username or group name
	telegramURL := fmt.Sprintf("https://t.me/%s", telegramName)

	// Generate QR code from Telegram URL
	qrCode, err := generateQRCode(telegramURL, size)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		log.Printf("generateTelegramQRCodeHandler: Failed to generate QR code - %v", err)
		return
	}

	// Open Telegram logo file
	telegramLogoFile, err := http.Dir(".").Open(TelegramLogoPath)
	if err != nil {
		http.Error(w, "Failed to open Telegram logo", http.StatusInternalServerError)
		log.Printf("generateTelegramQRCodeHandler: Failed to open Telegram logo - %v", err)
		return
	}
	defer telegramLogoFile.Close()

	// Decode Telegram logo image
	telegramLogo, err := decodeImage(telegramLogoFile)
	if err != nil {
		http.Error(w, "Failed to decode Telegram logo", http.StatusInternalServerError)
		log.Printf("generateTelegramQRCodeHandler: Failed to decode Telegram logo - %v", err)
		return
	}

	// Overlay Telegram logo on QR code
	qrCode, err = overlayImageOnQRCode(qrCode, telegramLogo, LogoPercent)
	if err != nil {
		http.Error(w, "Failed to overlay Telegram logo on QR code", http.StatusInternalServerError)
		log.Printf("generateTelegramQRCodeHandler: Failed to overlay Telegram logo on QR code - %v", err)
		return
	}

	// Set content type for QR code image
	w.Header().Set("Content-Type", "image/png")

	// Encode QR code as PNG and write to response
	err = png.Encode(w, qrCode)
	if err != nil {
		log.Printf("generateTelegramQRCodeHandler: Failed to encode QR code as PNG - %v", err)
	}
}

// Generates a QR code for joining a Zoom meeting.
func generateZoomQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	// Check for allowed method (POST only)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("generateZoomQRCodeHandler: Method not allowed")
		return
	}

	// Extract meeting ID and password from request form
	meetingID := r.FormValue("meetingID")
	password := r.FormValue("password")

	// Validate presence of meeting ID
	if meetingID == "" {
		http.Error(w, "Missing meeting ID", http.StatusBadRequest)
		log.Printf("generateZoomQRCodeHandler: Missing meeting ID")
		return
	}

	// Extract and validate QR code size
	sizeStr := r.FormValue("size")
	if sizeStr == "" {
		http.Error(w, "Missing size", http.StatusBadRequest)
		log.Printf("generateZoomQRCodeHandler: Missing size")
		return
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil || !isValidQRCodeSize(size) {
		http.Error(w, "Invalid size", http.StatusBadRequest)
		log.Printf("generateZoomQRCodeHandler: Invalid size - %v", err)
		return
	}

	// Generate Zoom meeting URL with meeting ID and password (optional)
	zoomURL := fmt.Sprintf("https://zoom.us/j/%s?pwd=%s", meetingID, password)

	// Generate QR code from Zoom meeting URL
	qrCode, err := generateQRCode(zoomURL, size)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		log.Printf("generateZoomQRCodeHandler: Failed to generate QR code - %v", err)
		return
	}

	// Open Zoom logo file
	zoomLogoFile, err := http.Dir(".").Open(ZoomLogoPath)
	if err != nil {
		http.Error(w, "Failed to open Zoom logo", http.StatusInternalServerError)
		log.Printf("generateZoomQRCodeHandler: Failed to open Zoom logo - %v", err)
		return
	}
	defer zoomLogoFile.Close()

	// Decode Zoom logo image
	zoomLogo, err := decodeImage(zoomLogoFile)
	if err != nil {
		http.Error(w, "Failed to decode Zoom logo", http.StatusInternalServerError)
		log.Printf("generateZoomQRCodeHandler: Failed to decode Zoom logo - %v", err)
		return
	}

	// Overlay Zoom logo on QR code
	qrCode, err = overlayImageOnQRCode(qrCode, zoomLogo, LogoPercent)
	if err != nil {
		http.Error(w, "Failed to overlay Zoom logo on QR code", http.StatusInternalServerError)
		log.Printf("generateZoomQRCodeHandler: Failed to overlay Zoom logo on QR code - %v", err)
		return
	}

	// Set content type for QR code image
	w.Header().Set("Content-Type", "image/png")

	// Encode QR code as PNG and write to response
	err = png.Encode(w, qrCode)
	if err != nil {
		log.Printf("generateZoomQRCodeHandler: Failed to encode QR code as PNG - %v", err)
	}
}

// overlayImageOnQRCodeWithOpacity overlays an image onto a QR code with a specified size and opacity.
func overlayImageOnQRCodeWithOpacity(qrCode image.Image, overlay image.Image, overlayPercent, overlayOpacity float64) (image.Image, error) {
	// Get the boundaries (width and height) of the QR code image
	qrBounds := qrCode.Bounds()
	qrWidth := qrBounds.Dx()
	qrHeight := qrBounds.Dy()

	// Calculate the maximum allowed width and height for the overlay image
	// based on the QR code dimensions and the provided percentage
	overlayMaxWidth := int(float64(qrWidth) * overlayPercent)
	overlayMaxHeight := int(float64(qrHeight) * overlayPercent)

	// Resize the overlay image to fit within the calculated maximum dimensions
	// while maintaining the aspect ratio using Lanczos resampling filter for better quality
	overlay = resize.Thumbnail(uint(overlayMaxWidth), uint(overlayMaxHeight), overlay, resize.Lanczos3)

	// Calculate the offset to center the overlay image on the QR code
	offset := image.Pt((qrWidth-overlay.Bounds().Dx())/2, (qrHeight-overlay.Bounds().Dy())/2)

	// Create a new image with the same bounds as the QR code
	b := qrBounds
	m := image.NewRGBA(b)

	// Draw the QR code onto the new image
	draw.Draw(m, qrBounds, qrCode, image.Point{}, draw.Src)

	// Apply the specified opacity to the overlay image
	overlay = applyOpacity(overlay, overlayOpacity)

	// Draw the overlaid image onto the new image with the calculated offset and "Over" compositing mode
	// which combines the overlay with the underlying QR code based on their alpha channels
	draw.Draw(m, overlay.Bounds().Add(offset), overlay, image.Point{}, draw.Over)

	// Return the new image with the overlaid image and any errors encountered
	return m, nil
}

// applyOpacity applies the specified opacity to an image.
func applyOpacity(img image.Image, opacity float64) image.Image {
	// Get the boundaries (width and height) of the image
	bounds := img.Bounds()

	// Create a new RGBA image with the same bounds as the original image
	newImg := image.NewRGBA(bounds)

	// Loop through each pixel of the image
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Get the original color of the pixel
			originalColor := img.At(x, y)

			// Extract the red, green, blue, and alpha (transparency) components
			r, g, b, a := originalColor.RGBA()

			// Apply the opacity by multiplying the original alpha value with the desired opacity
			a = uint32(float64(a) * opacity)

			// Create a new RGBA color with the original red, green, blue, and adjusted alpha values
			newColor := color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}

			// Set the pixel color in the new image
			newImg.Set(x, y, newColor)
		}
	}

	// Return the new image with the applied opacity
	return newImg
}
