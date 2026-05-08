package main

import (
	"fmt"
	"strings"
)

func generateVCardString(firstName, lastName, title, phone, mobile, email, address, company, url, role, lang, geo string) string {
	var sb strings.Builder
	sb.WriteString("BEGIN:VCARD\n")
	sb.WriteString("VERSION:3.0\n")
	sb.WriteString(fmt.Sprintf("N:%s;%s;;;\n", lastName, firstName))
	sb.WriteString(fmt.Sprintf("FN:%s %s\n", firstName, lastName))

	if company != "" {
		sb.WriteString(fmt.Sprintf("ORG:%s\n", company))
	}

	sb.WriteString(fmt.Sprintf("TITLE:%s\n", title))
	sb.WriteString(fmt.Sprintf("TEL;TYPE=WORK,VOICE:%s\n", phone))

	if mobile != "" {
		sb.WriteString(fmt.Sprintf("TEL;TYPE=CELL,VOICE:%s\n", mobile))
	}

	sb.WriteString(fmt.Sprintf("EMAIL:%s\n", email))
	sb.WriteString(fmt.Sprintf("ADR:%s\n", address))

	if url != "" {
		sb.WriteString(fmt.Sprintf("URL:%s\n", url))
	}
	if role != "" {
		sb.WriteString(fmt.Sprintf("ROLE:%s\n", role))
	}
	if lang != "" {
		sb.WriteString(fmt.Sprintf("LANG:%s\n", lang))
	}
	if geo != "" {
		sb.WriteString(fmt.Sprintf("GEO:%s\n", geo))
	}

	sb.WriteString("END:VCARD")
	return sb.String()
}
