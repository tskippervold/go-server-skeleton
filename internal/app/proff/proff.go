package proff

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	baseURL   = "https://api.proff.no"
	authToken = "OUxJYPlz1OiW8eauiJfIZMH9g"
)

type Company struct {
	Name               string        `json:"name"`
	OrganisationNumber string        `json:"organisationNumber"`
	PostalAddress      interface{}   `json:"postalAddress"`
	CompanyType        string        `json:"companyType"`
	RegisteredForVat   bool          `json:"registeredForVat"`
	RegistrationDate   string        `json:"registrationDate"`
	NaceCategories     []string      `json:"naceCategories"`
	Shareholders       []interface{} `json:"shareholders"`
}

type CompanyOwners struct {
	People                   []Person `json:"people"`
	TotalSharePercentage     float64  `json:"totalSharePercentage"`
	MaxSingleSharePercentage float64  `json:"maxSingleSharePercentage"`
}

type Person struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	SharePercentage float64 `json:"sharePercentage"`
	EntityType      string  `json:"entityType"`
}

func LookupCompany(orgnum string) (Company, error) {
	url := urlWithPath(fmt.Sprintf("%s/%s", "api/companies/register/NO", orgnum))
	method := "GET"

	req, err := http.NewRequest(method, url.String(), nil)
	if err != nil {
		return Company{}, err
	}

	res, err := executeRequest(req)
	if err != nil {
		return Company{}, err
	}

	var result Company
	err = json.Unmarshal(res, &result)
	if err != nil {
		return Company{}, err
	}

	return result, nil
}

func GetCompanyOwners(orgnum string) (CompanyOwners, error) {
	directOwners, err := getCompanyOwners(orgnum, true)
	if err != nil {
		return CompanyOwners{}, err
	}

	indirectOwners, err := getCompanyOwners(orgnum, false)
	if err != nil {
		return CompanyOwners{}, err
	}

	all := append(directOwners, indirectOwners...)
	var people []Person
	var sharePercentage float64
	var maxSingleSharePercentage float64

	for _, owner := range all {
		if strings.ToUpper(owner.EntityType) == "PERSON" {
			owner.SharePercentage = math.Round(owner.SharePercentage*100) / 100

			if owner.SharePercentage > 0 {
				people = append(people, owner)
				sharePercentage = sharePercentage + owner.SharePercentage
			}

			maxSingleSharePercentage = math.Max(maxSingleSharePercentage, owner.SharePercentage)
		}
	}

	sharePercentage = math.Min(math.Round(sharePercentage*100)/100, 100)
	return CompanyOwners{people, sharePercentage, maxSingleSharePercentage}, nil
}

func getCompanyOwners(orgnum string, directOwnership bool) ([]Person, error) {
	type response struct {
		Relations []Person `json:"relations"`
	}

	url := urlWithPath(fmt.Sprintf("%s/%s", "api/shareholders/eniropro/NO/owners", orgnum))
	query := url.Query()
	query.Set("direct", strconv.FormatBool(directOwnership))
	query.Set("pageSize", "1000")
	url.RawQuery = query.Encode()

	method := "GET"
	req, err := http.NewRequest(method, url.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := executeRequest(req)
	if err != nil {
		return nil, err
	}

	var result response
	err = json.Unmarshal(res, &result)
	if err != nil {
		return nil, err
	}

	return result.Relations, nil
}

func urlWithPath(path string) *url.URL {
	u, _ := url.Parse(fmt.Sprintf("%s/%s", baseURL, path))
	return u
}

func executeRequest(req *http.Request) ([]byte, error) {
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", authToken))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}
