package pl

import "encoding/xml"

type faDocument struct {
	XMLName xml.Name `xml:"http://crd.gov.pl/wzor/2025/06/25/13775/ Faktura"`
	Header  faHeader `xml:"Naglowek"`
	Seller  faParty  `xml:"Podmiot1"`
	Buyer   faParty  `xml:"Podmiot2"`
	Body    faBody   `xml:"Fa"`
}

type formCode struct {
	System string `xml:"kodSystemowy,attr"`
	Schema string `xml:"wersjaSchemy,attr"`
	Value  string `xml:",chardata"`
}

type faHeader struct {
	Code     formCode `xml:"KodFormularza"`
	Variant  int      `xml:"WariantFormularza"`
	Created  string   `xml:"DataWytworzeniaFa"`
	Software string   `xml:"SystemInfo"`
}

type identity struct {
	NIP       string `xml:"NIP,omitempty"`
	Country   string `xml:"KodKraju,omitempty"`
	ForeignID string `xml:"NrID,omitempty"`
	NoID      int    `xml:"BrakID,omitempty"`
	Name      string `xml:"Nazwa"`
}

type faAddress struct {
	Country string `xml:"KodKraju"`
	Street  string `xml:"AdresL1"`
	City    string `xml:"AdresL2"`
}
type faParty struct {
	Identity identity  `xml:"DaneIdentyfikacyjne"`
	Address  faAddress `xml:"Adres"`
	JST      int       `xml:"JST,omitempty"`
	GV       int       `xml:"GV,omitempty"`
}
type annotation struct {
	Cash       int `xml:"P_16"`
	Self       int `xml:"P_17"`
	Reverse    int `xml:"P_18"`
	Split      int `xml:"P_18A"`
	Exemption  int `xml:"Zwolnienie>P_19N"`
	Transport  int `xml:"NoweSrodkiTransportu>P_22N"`
	Simplified int `xml:"P_23"`
	Margin     int `xml:"PMarzy>P_PMarzyN"`
}
type faBody struct {
	Currency string     `xml:"KodWaluty"`
	Date     string     `xml:"P_1"`
	Number   string     `xml:"P_2"`
	Sums     []taxValue `xml:",any"`
	Total    string     `xml:"P_15"`
	Notes    annotation `xml:"Adnotacje"`
	Kind     string     `xml:"RodzajFaktury"`
	Rows     []faLine   `xml:"FaWiersz"`
}
type taxValue struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}
type faLine struct {
	Position int    `xml:"NrWierszaFa"`
	Name     string `xml:"P_7"`
	SKU      string `xml:"Indeks,omitempty"`
	Unit     string `xml:"P_8A,omitempty"`
	Quantity string `xml:"P_8B"`
	Price    string `xml:"P_9B"`
	Gross    string `xml:"P_11A"`
	Tax      string `xml:"P_11Vat"`
	Rate     string `xml:"P_12"`
}
