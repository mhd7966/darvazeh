package models
 
type Zone struct {
	Id               string   `json:"id"`
	Name             string   `json:"name"`
	Type             string   `json:"type" default:"Zone"`
	URL              string   `json:"url"`
	Kind             string   `json:"kine"`
	RRSets           []RRSet  `json:"rrsets"`
	Serial           int      `json:"serial"`
	NotifiedSerial   int      `json:"notified_serial"`
	EditedSerial     int      `json:"edited_serial"`
	Masters          []string `json:"masters"`
	DNSSec           bool     `json:"dnsdec"`
	NSEC3Param       string   `json:"nsec3param"`
	NSEC3Narrow      bool     `json:"nsec3narrow"`
	Presigned        bool     `json:"presigned"`
	SOAEdit          string   `json:"soa_edit"`
	SOAEditApi       string   `json:"soa_edit_api"`
	ApiRectify       bool     `json:"api_rectify"`
	Zone             string   `json:"zone"`
	Account          string   `json:"account" binding:"required"`
	NameServers      []string `json:"nameservers"`
	MasterTSIGKeyIds []string `json:"master_tsig_key_ids"`
	SlaveTSIGKeyIds  []string `json:"slave_tsig_key_ids"`
}

/**
Zone{

description:	
This represents an authoritative DNS Zone.

id	string
Opaque zone id (string), assigned by the server, should not be interpreted by the application. Guaranteed to be safe for embedding in URLs.

name	string
Name of the zone (e.g. “example.com.”) MUST have a trailing dot

type	string
Set to “Zone”

url	string
API endpoint for this zone

kind	string
Zone kind, one of “Native”, “Master”, “Slave”

Enum:

Array [ 3 ]

rrsets	[...]

serial	integer
The SOA serial number

notified_serial	integer
The SOA serial notifications have been sent out for

edited_serial	integer
The SOA serial as seen in query responses. Calculated using the SOA-EDIT metadata, default-soa-edit and default-soa-edit-signed settings

masters	[...]

dnssec	boolean
Whether or not this zone is DNSSEC signed (inferred from presigned being true XOR presence of at least one cryptokey with active being true)

nsec3param	string
The NSEC3PARAM record

nsec3narrow	boolean
Whether or not the zone uses NSEC3 narrow

presigned	boolean
Whether or not the zone is pre-signed

soa_edit	string
The SOA-EDIT metadata item

soa_edit_api	string
The SOA-EDIT-API metadata item

api_rectify	boolean
Whether or not the zone will be rectified on data changes via the API

zone	string
MAY contain a BIND-style zone file when creating a zone

account	string
MAY be set. Its value is defined by local policy

nameservers	[...]

master_tsig_key_ids	[...]

slave_tsig_key_ids	[...]
 
}
**/