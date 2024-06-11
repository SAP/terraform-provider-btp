package xsuaa_api

type ApiCredentialGetBody struct {
	CredentialType  	string	`json:"credential-type"`
	ClientId			string 	`json:"clientid"`
	Certificate			string	`json:"certificate,omitempty"`
	TokenUrl			string 	`json:"tokenurl"`
	ClientSecret 		string 	`json:"clientsecret,omitempty"`
	ApiUrl				string	`json:"apiurl"`
	Name				string	`json:"name"`
	ReadOnly			bool	`json:"read-only"`
}