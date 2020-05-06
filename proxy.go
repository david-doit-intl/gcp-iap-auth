package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"gcp-iap-auth/jwt"

	"github.com/go-ldap/ldap"
)

type proxy struct {
	backend     *url.URL
	emailHeader string
	proxy       *httputil.ReverseProxy
}

func newProxy(backendURL, emailHeader string) (*proxy, error) {
	backend, err := url.Parse(backendURL)
	if err != nil {
		return nil, fmt.Errorf("Could not parse URL '%s': %s", backendURL, err)
	}
	return &proxy{
		backend:     backend,
		emailHeader: emailHeader,
		proxy:       httputil.NewSingleHostReverseProxy(backend),
	}, nil
}

func setEmailHeader(emailheader string, email string, req *http.Request) {
	if *emailHeader != "" {
		req.Header.Set(*emailHeader, email)
	}
}

func checkValidEmail(email string) bool {
	l, err := ldap.DialURL(*ldapURL)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer l.Close()

	searchRequest := ldap.NewSearchRequest(
		*baseDomain, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=organizationalPerson))", // The filter to apply
		[]string{"dn", "cn"},                    // A list attributes to retrieve
		nil,
	)

	return l.Search(searchRequest)
}

func (p *proxy) handler(res http.ResponseWriter, req *http.Request) {
	claims, err := jwt.RequestClaims(req, cfg)
	email := claims.Email
	if err != nil {
		if claims == nil || len(email) == 0 {
			log.Printf("Failed to authenticate %q (%v)\n", email, err)
		}
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}

	setEmailHeader(p.emailHeader, email, req)
	// Check if email is in ldap server
	if checkValidEmail(email) {
		p.proxy.ServeHTTP(res, req)
		return
	}

	http.Error(res, "Email not found in LDAP server", http.StatusUnauthorized)
}
