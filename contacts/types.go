package contacts

// Contact represents a contact in the Nylas API.
type Contact struct {
	// ID is the unique identifier for this contact.
	ID string `json:"id"`
	// GrantID is the ID of the grant (connected account) this contact belongs to.
	GrantID string `json:"grant_id"`
	// Object is the object type, always "contact".
	Object string `json:"object,omitempty"`
	// Birthday is the contact's birthday in YYYY-MM-DD format.
	Birthday string `json:"birthday,omitempty"`
	// CompanyName is the contact's employer or organization.
	CompanyName string `json:"company_name,omitempty"`
	// DisplayName is the full display name.
	DisplayName string `json:"display_name,omitempty"`
	// Emails contains the contact's email addresses.
	Emails []Email `json:"emails,omitempty"`
	// GivenName is the contact's first name.
	GivenName string `json:"given_name,omitempty"`
	// Groups contains the contact groups this contact belongs to.
	Groups []Group `json:"groups,omitempty"`
	// IMAddresses contains instant messaging addresses.
	IMAddresses []IMAddress `json:"im_addresses,omitempty"`
	// JobTitle is the contact's job title.
	JobTitle string `json:"job_title,omitempty"`
	// ManagerName is the contact's manager's name.
	ManagerName string `json:"manager_name,omitempty"`
	// MiddleName is the contact's middle name.
	MiddleName string `json:"middle_name,omitempty"`
	// Nickname is the contact's nickname.
	Nickname string `json:"nickname,omitempty"`
	// Notes contains free-form notes about the contact.
	Notes string `json:"notes,omitempty"`
	// OfficeLocation is the contact's office location.
	OfficeLocation string `json:"office_location,omitempty"`
	// PhoneNumbers contains the contact's phone numbers.
	PhoneNumbers []PhoneNumber `json:"phone_numbers,omitempty"`
	// PhysicalAddresses contains the contact's physical addresses.
	PhysicalAddresses []PhysicalAddress `json:"physical_addresses,omitempty"`
	// Picture is base64-encoded picture data.
	Picture string `json:"picture,omitempty"`
	// PictureURL is a URL to the contact's picture.
	PictureURL string `json:"picture_url,omitempty"`
	// Suffix is a name suffix (e.g., "Jr.", "III").
	Suffix string `json:"suffix,omitempty"`
	// Surname is the contact's last name.
	Surname string `json:"surname,omitempty"`
	// Source indicates the origin of the contact data.
	Source string `json:"source,omitempty"`
	// WebPages contains the contact's websites and social profiles.
	WebPages []WebPage `json:"web_pages,omitempty"`
}

// Email represents an email address for a contact.
type Email struct {
	// Email is the email address.
	Email string `json:"email"`
	// Type is the address type (e.g., "work", "home", "other").
	Type string `json:"type,omitempty"`
}

// PhoneNumber represents a phone number for a contact.
type PhoneNumber struct {
	// Number is the phone number.
	Number string `json:"number"`
	// Type is the number type (e.g., "mobile", "work", "home").
	Type string `json:"type,omitempty"`
}

// PhysicalAddress represents a physical/mailing address for a contact.
type PhysicalAddress struct {
	// Format is the address format type.
	Format string `json:"format,omitempty"`
	// StreetAddress is the street address.
	StreetAddress string `json:"street_address,omitempty"`
	// City is the city name.
	City string `json:"city,omitempty"`
	// PostalCode is the postal/ZIP code.
	PostalCode string `json:"postal_code,omitempty"`
	// State is the state or province.
	State string `json:"state,omitempty"`
	// Country is the country name.
	Country string `json:"country,omitempty"`
	// Type is the address type (e.g., "work", "home").
	Type string `json:"type,omitempty"`
}

// WebPage represents a website or social profile URL for a contact.
type WebPage struct {
	// URL is the web page URL.
	URL string `json:"url"`
	// Type is the page type (e.g., "profile", "blog", "homepage").
	Type string `json:"type,omitempty"`
}

// IMAddress represents an instant messaging address.
type IMAddress struct {
	// IMAddress is the instant messaging handle/username.
	IMAddress string `json:"im_address"`
	// Type is the IM service type (e.g., "aim", "gtalk", "skype").
	Type string `json:"type,omitempty"`
}

// Group represents a contact group/label.
type Group struct {
	// ID is the group's unique identifier.
	ID string `json:"id"`
	// Name is the group's display name.
	Name string `json:"name,omitempty"`
}

// ListOptions specifies options for listing contacts.
// All fields are optional; nil values are not included in the request.
type ListOptions struct {
	// Limit is the maximum number of contacts to return.
	Limit *int `json:"limit,omitempty"`
	// PageToken is the cursor for pagination.
	PageToken string `json:"page_token,omitempty"`
	// Email filters contacts by email address.
	Email *string `json:"email,omitempty"`
	// PhoneNumber filters contacts by phone number.
	PhoneNumber *string `json:"phone_number,omitempty"`
	// Source filters contacts by source (e.g., "address_book", "domain_contact").
	Source *string `json:"source,omitempty"`
	// Group filters contacts by group ID.
	Group *string `json:"group,omitempty"`
	// Recurse includes contacts from sub-groups.
	Recurse *bool `json:"recurse,omitempty"`
}

// Values converts ListOptions to URL query parameters.
func (o *ListOptions) Values() map[string]any {
	if o == nil {
		return nil
	}
	v := make(map[string]any)
	if o.Limit != nil {
		v["limit"] = *o.Limit
	}
	if o.PageToken != "" {
		v["page_token"] = o.PageToken
	}
	if o.Email != nil {
		v["email"] = *o.Email
	}
	if o.PhoneNumber != nil {
		v["phone_number"] = *o.PhoneNumber
	}
	if o.Source != nil {
		v["source"] = *o.Source
	}
	if o.Group != nil {
		v["group"] = *o.Group
	}
	if o.Recurse != nil {
		v["recurse"] = *o.Recurse
	}
	return v
}

// CreateRequest represents a request to create a contact.
type CreateRequest struct {
	// Birthday is the contact's birthday (YYYY-MM-DD format).
	Birthday string `json:"birthday,omitempty"`
	// CompanyName is the contact's employer.
	CompanyName string `json:"company_name,omitempty"`
	// DisplayName is the full display name.
	DisplayName string `json:"display_name,omitempty"`
	// Emails is the list of email addresses.
	Emails []Email `json:"emails,omitempty"`
	// GivenName is the first name.
	GivenName string `json:"given_name,omitempty"`
	// Groups is the list of contact groups to add to.
	Groups []Group `json:"groups,omitempty"`
	// IMAddresses is the list of IM addresses.
	IMAddresses []IMAddress `json:"im_addresses,omitempty"`
	// JobTitle is the job title.
	JobTitle string `json:"job_title,omitempty"`
	// ManagerName is the manager's name.
	ManagerName string `json:"manager_name,omitempty"`
	// MiddleName is the middle name.
	MiddleName string `json:"middle_name,omitempty"`
	// Nickname is the nickname.
	Nickname string `json:"nickname,omitempty"`
	// Notes is free-form notes.
	Notes string `json:"notes,omitempty"`
	// OfficeLocation is the office location.
	OfficeLocation string `json:"office_location,omitempty"`
	// PhoneNumbers is the list of phone numbers.
	PhoneNumbers []PhoneNumber `json:"phone_numbers,omitempty"`
	// PhysicalAddresses is the list of addresses.
	PhysicalAddresses []PhysicalAddress `json:"physical_addresses,omitempty"`
	// Suffix is the name suffix.
	Suffix string `json:"suffix,omitempty"`
	// Surname is the last name.
	Surname string `json:"surname,omitempty"`
	// WebPages is the list of websites.
	WebPages []WebPage `json:"web_pages,omitempty"`
}

// UpdateRequest represents a request to update a contact.
type UpdateRequest struct {
	// Birthday is the contact's birthday (YYYY-MM-DD format).
	Birthday *string `json:"birthday,omitempty"`
	// CompanyName is the contact's employer.
	CompanyName *string `json:"company_name,omitempty"`
	// DisplayName is the full display name.
	DisplayName *string `json:"display_name,omitempty"`
	// Emails replaces the email addresses.
	Emails []Email `json:"emails,omitempty"`
	// GivenName is the first name.
	GivenName *string `json:"given_name,omitempty"`
	// Groups replaces the contact groups.
	Groups []Group `json:"groups,omitempty"`
	// IMAddresses replaces the IM addresses.
	IMAddresses []IMAddress `json:"im_addresses,omitempty"`
	// JobTitle is the job title.
	JobTitle *string `json:"job_title,omitempty"`
	// ManagerName is the manager's name.
	ManagerName *string `json:"manager_name,omitempty"`
	// MiddleName is the middle name.
	MiddleName *string `json:"middle_name,omitempty"`
	// Nickname is the nickname.
	Nickname *string `json:"nickname,omitempty"`
	// Notes is free-form notes.
	Notes *string `json:"notes,omitempty"`
	// OfficeLocation is the office location.
	OfficeLocation *string `json:"office_location,omitempty"`
	// PhoneNumbers replaces the phone numbers.
	PhoneNumbers []PhoneNumber `json:"phone_numbers,omitempty"`
	// PhysicalAddresses replaces the addresses.
	PhysicalAddresses []PhysicalAddress `json:"physical_addresses,omitempty"`
	// Suffix is the name suffix.
	Suffix *string `json:"suffix,omitempty"`
	// Surname is the last name.
	Surname *string `json:"surname,omitempty"`
	// WebPages replaces the websites.
	WebPages []WebPage `json:"web_pages,omitempty"`
}
