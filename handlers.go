package main

import (
	"errors"
	"image/png"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func isValidQRCodeSize(size int) bool {
	return size == QRSmall || size == QRMedium || size == QRLarge || size == QRExtraLarge
}

var (
	errInvalidLatitude    = errors.New("invalid latitude")
	errInvalidLongitude   = errors.New("invalid longitude")
	errInvalidSecurity    = errors.New("invalid security type")
	errPasswordRequiredWPA = errors.New("password is required for WPA/WPA2/WPA3 security")
	errPasswordLengthWPA  = errors.New("password for WPA/WPA2/WPA3 must be between 8 and 63 characters")
	errPasswordLengthWEP  = errors.New("password for WEP must be exactly 5 or 13 characters")
)

type qrHandlerConfig struct {
	name         string
	required     []string
	validates    func(r *http.Request) error
	buildContent func(r *http.Request) (string, error)
	logoPath     string
}

func handleQR(w http.ResponseWriter, r *http.Request, cfg qrHandlerConfig) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("%s: Method not allowed", cfg.name)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			log.Printf("%s: Failed to parse form - %v", cfg.name, err)
			return
		}
	}

	formValue := func(key string) string {
		if r.MultipartForm != nil {
			if vs := r.MultipartForm.Value[key]; len(vs) > 0 {
				return vs[0]
			}
		}
		return r.FormValue(key)
	}

	for _, field := range cfg.required {
		if formValue(field) == "" {
			http.Error(w, "Missing "+field, http.StatusBadRequest)
			log.Printf("%s: Missing %s (form: %v)", cfg.name, field, r.Form)
			return
		}
	}

	sizeStr := formValue("size")
	log.Printf("%s: form fields: %v, size=%q", cfg.name, r.Form, sizeStr)
	size, err := strconv.Atoi(sizeStr)
	if err != nil || !isValidQRCodeSize(size) {
		http.Error(w, "Invalid size", http.StatusBadRequest)
		log.Printf("%s: Invalid size - %v", cfg.name, err)
		return
	}

	if cfg.validates != nil {
		if err := cfg.validates(r); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("%s: Validation error - %v", cfg.name, err)
			return
		}
	}

	content, err := cfg.buildContent(r)
	if err != nil {
		http.Error(w, "Failed to build content", http.StatusInternalServerError)
		log.Printf("%s: Failed to build content - %v", cfg.name, err)
		return
	}

	qrCode, err := generateQRCode(content, size)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		log.Printf("%s: Failed to generate QR code - %v", cfg.name, err)
		return
	}

	// Check for custom image upload
	file, _, err := r.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "Error reading image", http.StatusInternalServerError)
		log.Printf("%s: Error reading image - %v", cfg.name, err)
		return
	}

	if file != nil {
		overlayImage, err := decodeImage(file)
		if err != nil {
			http.Error(w, "Failed to decode image", http.StatusInternalServerError)
			log.Printf("%s: Failed to decode image - %v", cfg.name, err)
			return
		}

		logoWidthPercent := LogoPercent
		if v := r.FormValue("logoWidthPercent"); v != "" {
			parsed, err := strconv.ParseFloat(v, 64)
			if err != nil {
				http.Error(w, "Invalid logo width percent", http.StatusBadRequest)
				log.Printf("%s: Invalid logo width percent - %v", cfg.name, err)
				return
			}
			logoWidthPercent = parsed
		}

		logoOpacity := 1.0
		if v := r.FormValue("logoOpacity"); v != "" {
			parsed, err := strconv.ParseFloat(v, 64)
			if err == nil {
				logoOpacity = parsed
			}
		}

		qrCode, err = overlayImageOnQRCodeWithOpacity(qrCode, overlayImage, logoWidthPercent, logoOpacity)
		if err != nil {
			http.Error(w, "Failed to overlay image on QR code", http.StatusInternalServerError)
			log.Printf("%s: Failed to overlay image on QR code - %v", cfg.name, err)
			return
		}
	} else if cfg.logoPath != "" {
		logoFile, err := http.Dir(".").Open(cfg.logoPath)
		if err != nil {
			http.Error(w, "Failed to open logo", http.StatusInternalServerError)
			log.Printf("%s: Failed to open logo - %v", cfg.name, err)
			return
		}
		defer logoFile.Close()

		logo, err := decodeImage(logoFile)
		if err != nil {
			http.Error(w, "Failed to decode logo", http.StatusInternalServerError)
			log.Printf("%s: Failed to decode logo - %v", cfg.name, err)
			return
		}

		qrCode, err = overlayImageOnQRCode(qrCode, logo, LogoPercent)
		if err != nil {
			http.Error(w, "Failed to overlay logo on QR code", http.StatusInternalServerError)
			log.Printf("%s: Failed to overlay logo on QR code - %v", cfg.name, err)
			return
		}
	}

	w.Header().Set("Content-Type", "image/png")
	if err := png.Encode(w, qrCode); err != nil {
		log.Printf("%s: Failed to encode QR code as PNG - %v", cfg.name, err)
	}
}

// --- Simple social handlers ---

func generateInstagramQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	handleQR(w, r, qrHandlerConfig{
		name:     "generateInstagram",
		required: []string{"username"},
		logoPath: "static/instagram_logo.png",
		buildContent: func(r *http.Request) (string, error) {
			return "https://www.instagram.com/" + r.FormValue("username"), nil
		},
	})
}

func generateFacebookQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	handleQR(w, r, qrHandlerConfig{
		name:     "generateFacebook",
		required: []string{"username"},
		logoPath: "static/facebook_logo.png",
		buildContent: func(r *http.Request) (string, error) {
			return "https://www.facebook.com/" + r.FormValue("username"), nil
		},
	})
}

func generateTikTokQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	handleQR(w, r, qrHandlerConfig{
		name:     "generateTikTok",
		required: []string{"username"},
		logoPath: "static/tiktok_logo.png",
		buildContent: func(r *http.Request) (string, error) {
			return "https://www.tiktok.com/@" + r.FormValue("username"), nil
		},
	})
}

func generateLinkedInQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	handleQR(w, r, qrHandlerConfig{
		name:     "generateLinkedIn",
		required: []string{"username"},
		logoPath: "static/linkedin_logo.png",
		buildContent: func(r *http.Request) (string, error) {
			return "https://www.linkedin.com/in/" + r.FormValue("username"), nil
		},
	})
}

func generateYouTubeQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	handleQR(w, r, qrHandlerConfig{
		name:     "generateYouTube",
		required: []string{"channel"},
		logoPath: "static/youtube_logo.png",
		buildContent: func(r *http.Request) (string, error) {
			return "https://www.youtube.com/channel/" + r.FormValue("channel"), nil
		},
	})
}

func generateXQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	handleQR(w, r, qrHandlerConfig{
		name:     "generateX",
		required: []string{"username"},
		logoPath: "static/x_logo.png",
		buildContent: func(r *http.Request) (string, error) {
			return "https://x.com/" + r.FormValue("username"), nil
		},
	})
}

func generateTelegramQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	handleQR(w, r, qrHandlerConfig{
		name:     "generateTelegram",
		required: []string{"telegramName"},
		logoPath: "static/telegram_logo.png",
		buildContent: func(r *http.Request) (string, error) {
			return "https://t.me/" + r.FormValue("telegramName"), nil
		},
	})
}

func generateSpotifyQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	handleQR(w, r, qrHandlerConfig{
		name:     "generateSpotify",
		required: []string{"spotifyURL"},
		logoPath: "static/spotify_logo.png",
		buildContent: func(r *http.Request) (string, error) {
			return r.FormValue("spotifyURL"), nil
		},
	})
}

// --- Handlers with custom validation ---

func generateMapQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	handleQR(w, r, qrHandlerConfig{
		name:     "generateMap",
		required: []string{"latitude", "longitude"},
		logoPath: "static/map_logo.png",
		validates: func(r *http.Request) error {
			lat, err := strconv.ParseFloat(r.FormValue("latitude"), 64)
			if err != nil || lat < -90 || lat > 90 {
				return errInvalidLatitude
			}
			lon, err := strconv.ParseFloat(r.FormValue("longitude"), 64)
			if err != nil || lon < -180 || lon > 180 {
				return errInvalidLongitude
			}
			return nil
		},
		buildContent: func(r *http.Request) (string, error) {
			return "geo:" + r.FormValue("latitude") + "," + r.FormValue("longitude"), nil
		},
	})
}

func generateWiFiQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	handleQR(w, r, qrHandlerConfig{
		name:     "generateWiFi",
		required: []string{"ssid"},
		logoPath: "static/wifi_logo.png",
		validates: func(r *http.Request) error {
			security := r.FormValue("security")
			password := r.FormValue("password")
			validSecurities := map[string]bool{"WPA": true, "WPA2": true, "WPA3": true, "WEP": true, "nopass": true}

			if !validSecurities[security] {
				return errInvalidSecurity
			}
			if security == "WPA" || security == "WPA2" || security == "WPA3" {
				if password == "" {
					return errPasswordRequiredWPA
				}
				if len(password) < 8 || len(password) > 63 {
					return errPasswordLengthWPA
				}
			}
			if security == "WEP" && len(password) != 5 && len(password) != 13 {
				return errPasswordLengthWEP
			}
			return nil
		},
		buildContent: func(r *http.Request) (string, error) {
			return "WIFI:T:" + r.FormValue("security") + ";S:" + r.FormValue("ssid") + ";P:" + r.FormValue("password") + ";;", nil
		},
	})
}

func generateEventQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	handleQR(w, r, qrHandlerConfig{
		name:     "generateEvent",
		required: []string{"eventName", "startDateTime", "endDateTime"},
		logoPath: "static/event_logo.png",
		buildContent: func(r *http.Request) (string, error) {
			return "BEGIN:VEVENT\nSUMMARY:" + r.FormValue("eventName") +
				"\nDTSTART:" + r.FormValue("startDateTime") +
				"\nDTEND:" + r.FormValue("endDateTime") +
				"\nLOCATION:" + r.FormValue("location") +
				"\nDESCRIPTION:" + r.FormValue("description") +
				"\nEND:VEVENT", nil
		},
	})
}

func generatePayPalQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	handleQR(w, r, qrHandlerConfig{
		name:     "generatePayPal",
		required: []string{"email", "amount", "currency"},
		logoPath: "static/paypal_logo.png",
		buildContent: func(r *http.Request) (string, error) {
			return "https://www.paypal.com/cgi-bin/webscr?cmd=_xclick&business=" +
				r.FormValue("email") + "&amount=" + r.FormValue("amount") +
				"&currency_code=" + r.FormValue("currency") +
				"&item_name=" + r.FormValue("description"), nil
		},
	})
}

func generateWhatsAppQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	handleQR(w, r, qrHandlerConfig{
		name:     "generateWhatsApp",
		required: []string{"phone"},
		logoPath: "static/whatsapp_logo.png",
		buildContent: func(r *http.Request) (string, error) {
			return "https://wa.me/" + r.FormValue("phone") + "?text=" + r.FormValue("message"), nil
		},
	})
}

func generatePhoneQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	handleQR(w, r, qrHandlerConfig{
		name:     "generatePhone",
		required: []string{"phoneNumber"},
		logoPath: "static/phone_logo.png",
		buildContent: func(r *http.Request) (string, error) {
			return "tel:" + r.FormValue("phoneNumber"), nil
		},
	})
}

func generateEmailQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	handleQR(w, r, qrHandlerConfig{
		name:     "generateEmail",
		required: []string{"email"},
		logoPath: "static/email_logo.png",
		buildContent: func(r *http.Request) (string, error) {
			return "mailto:" + r.FormValue("email") + "?subject=" +
				url.QueryEscape(r.FormValue("subject")) + "&body=" +
				url.QueryEscape(r.FormValue("body")), nil
		},
	})
}

func generateSMSQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	handleQR(w, r, qrHandlerConfig{
		name:     "generateSMS",
		required: []string{"phoneNumber"},
		logoPath: "static/sms_logo.png",
		buildContent: func(r *http.Request) (string, error) {
			return "sms:" + r.FormValue("phoneNumber") + "?body=" + r.FormValue("message"), nil
		},
	})
}

func generateZoomQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	handleQR(w, r, qrHandlerConfig{
		name:     "generateZoom",
		required: []string{"meetingID"},
		logoPath: "static/zoom_logo.png",
		buildContent: func(r *http.Request) (string, error) {
			return "https://zoom.us/j/" + r.FormValue("meetingID") + "?pwd=" + r.FormValue("password"), nil
		},
	})
}

// --- Handlers with custom image upload support ---

func generateQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	handleQR(w, r, qrHandlerConfig{
		name:     "generateQRCode",
		required: []string{"url"},
		buildContent: func(r *http.Request) (string, error) {
			return r.FormValue("url"), nil
		},
	})
}

func generateVCardQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	handleQR(w, r, qrHandlerConfig{
		name: "generateVCard",
		buildContent: func(r *http.Request) (string, error) {
			return generateVCardString(
				r.FormValue("firstName"), r.FormValue("lastName"),
				r.FormValue("title"), r.FormValue("phone"),
				r.FormValue("mobile"), r.FormValue("email"),
				r.FormValue("address"), r.FormValue("company"),
				r.FormValue("url"), r.FormValue("role"),
				r.FormValue("lang"), r.FormValue("geo"),
			), nil
		},
	})
}
