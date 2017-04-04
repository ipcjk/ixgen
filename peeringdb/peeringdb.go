package peeringdb

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

type peeringdb struct {
	apiURL string
}

func Peeringdb(apiURL string) peeringdb {
	return peeringdb{apiURL}
}

func (p *peeringdb) callAPI(uri string, i interface{}) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", p.apiURL+uri, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("User-Agent", "joerg-golang")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("HTTP request not successful: %s", err)
	}

	if resp.StatusCode != 200 {
		log.Fatalf("HTTP Api server responded with %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&i)
	if err != nil {
		log.Fatalf("Problems decoding http api output: %s", err)
	}
}

func (p *peeringdb) GetPeersOnIXByIxLanID(ixLanID int64) (apiResult Netixlan) {
	v := url.Values{}
	v.Set("ixlan_id", strconv.FormatInt(ixLanID, 10))
	p.callAPI("/netixlan?"+v.Encode(), &apiResult)
	return
}

func (p *peeringdb) GetPeersOnIXByIxName(ixName string) (apiResult Netixlan) {
	var ixlanid string
	nameAndNet := strings.Split(ixName, "/")

	iX := p.SearchIXByIxName(nameAndNet[0])
	v := url.Values{}

	if len(nameAndNet) > 1 && nameAndNet[1] != "" {
		for _, value := range iX.Data[0].IxlanSet {
			if value.Name == nameAndNet[1] {
				ixlanid = strconv.FormatInt(value.ID, 10)
				goto end
			}
		}
		log.Fatalf("Attention, Net: %s given for %s, but not found", nameAndNet[1], nameAndNet[0])
	} else if len(iX.Data[0].IxlanSet) > 1 {
		log.Fatalf("There a multiple Nets to choose for %s, please specify in the ini-file\n", nameAndNet[0])
	} else {
		ixlanid = strconv.FormatInt(iX.Data[0].IxlanSet[0].ID, 10)
	}

end:
	v.Set("ixlan_id", ixlanid)
	p.callAPI("/netixlan?"+v.Encode(), &apiResult)
	sort.Sort(netsIxLanDataSortedByASN(apiResult.Data))
	return
}

func (p *peeringdb) GetIxLANByIxLanID(ixLanID int64) (apiResult IxLAN) {
	p.callAPI("/ixlan/"+strconv.FormatInt(ixLanID, 10), &apiResult)
	return
}

func (p *peeringdb) ListIX() (apiResult Ix) {
	p.callAPI("/ix", &apiResult)
	return
}

func (p *peeringdb) SearchIXByIxName(ixName string) (apiResult Ix) {
	v := url.Values{}
	v.Set("name", ixName)
	p.callAPI("/ix?"+v.Encode(), &apiResult)

	if len(apiResult.Data) == 0 {
		log.Fatalf("%s is not a valid ixName or was not found on peeringdb", ixName)
	}

	p.callAPI("/ix/"+strconv.FormatInt(apiResult.Data[0].ID, 10), &apiResult)
	return
}

func (p *peeringdb) ListFaculty() (apiResult Ix) {
	p.callAPI("/fac", &apiResult)
	return
}

func (p *peeringdb) SearchFacultyByFacName(facultyName string) (apiResult Ix) {
	v := url.Values{}
	v.Set("name", facultyName)
	p.callAPI("/fac?"+v.Encode(), &apiResult)
	return
}

func (p *peeringdb) GetNetworkByAsN(asn int64) (apiResult Net) {
	v := url.Values{}
	v.Set("asn", strconv.FormatInt(asn, 10))
	p.callAPI("/net?"+v.Encode(), &apiResult)
	return
}
