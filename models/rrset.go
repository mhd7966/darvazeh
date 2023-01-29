package models

import "database/sql"

type RecordBody struct{
	ID int `json:"id"`
	Domain  string `json:"domain" binding:"required"`
	Name string `json:"name" binding:"required"`
	Type string `json:"type" binding:"required"`
	TTL int `json:"ttl" binding:"required" default:"3600"`
	Value string `json:"value" binding:"required"`
	Priority int `json:"priority" default:"7000"`
}

type RecordModel struct{
	Name string `json:"name" binding:"required"`
	RRSet []RRSet `json:"rrsets" binding:"required"`
}

type RRSet struct {
	Name       string    `json:"name" binding:"required"`
	Type       string    `json:"type" binding:"required"`
	TTL        int       `json:"ttl" binding:"required" default:"360"`
	ChangeType string    `json:"changetype" default:"REPLACE"`

	Records    []Record  `json:"records" binding:"required"`
	Comments   []Comment `json:"comments"`
}

/**
RRSet{

description:	
This represents a Resource Record Set (all records with the same name and type).

name*	string
Name for record set (e.g. “www.powerdns.com.”)

type*	string
Type of this record (e.g. “A”, “PTR”, “MX”)

ttl*	integer
DNS TTL of the records, in seconds. MUST NOT be included when changetype is set to “DELETE”.

changetype*	string
MUST be added when updating the RRSet. Must be REPLACE or DELETE. With DELETE, all existing RRs matching name and type will be deleted, including all comments. With REPLACE: when records is present, all existing RRs matching name and type will be deleted, and then new records given in records will be created. If no records are left, any existing comments will be deleted as well. When comments is present, all existing comments for the RRs matching name and type will be deleted, and then new comments given in comments will be created. 
}
**/

type Record struct {
	Content  string `json:"content" binding:"required"`
	Disabled bool   `json:"disabled"`
}

/**
Record{

description:	
The RREntry object represents a single record.

content*	string
The content of this record

disabled	boolean
Whether or not this record is disabled. When unset, the record is not disabled
 
}
**/
type Comment struct {
	Content    string `json:"content"`
	Account    string `json:"account"`
	ModifiedAt int    `json:"modified_at"`
}

/**
Comment{

description:	
A comment about an RRSet.

content	string
The actual comment

account	string
Name of an account that added the comment

modified_at	integer
Timestamp of the last change to the comment
 
}
**/

type RecordDB struct{
	ID int `json:"id"`
	DomainID int `json:"domain_id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Content string `json:"content"`
	TTL int `json:"ttl"`
	Priority int `json:"prio"`
	Disabled int `json:"disabled"`
	OrderName sql.NullString `json:"ordername"`
	Auth int `json:"auth"`
}