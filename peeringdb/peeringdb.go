package peeringdb

import (
	"encoding/json"
	"fmt"
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

func (p *peeringdb) callAPI(uri string, i interface{}) error {
	client := &http.Client{}

	req, err := http.NewRequest("GET", p.apiURL+uri, nil)
	if err != nil {
		return fmt.Errorf("can not generate new http request: %s", err)
	}

	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("User-Agent", "ixgen/golang")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request not successful: %s", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP Api server responded with %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&i)
	if err != nil {
		return fmt.Errorf("Problems decoding http api output: %s", err)
	}

	return nil
}

func (p *peeringdb) GetPeersOnIXByIxLanID(ixLanID int64) (apiResult Netixlan, err error) {
	v := url.Values{}
	v.Set("ixlan_id", strconv.FormatInt(ixLanID, 10))
	if err = p.callAPI("/netixlan?"+v.Encode(), &apiResult); err != nil {
		return Netixlan{}, err
	}
	return apiResult, nil
}

func (p *peeringdb) GetPeersOnIX(ixName string, ixId string, byId bool) (apiResult Netixlan, err error) {
	var ixlanid string
	var iX Ix
	var nameAndNet []string

	if byId {
		iX, err = p.SearchIXByIxId(ixId)
		if err != nil {
			return Netixlan{}, err
		}
	} else {
		nameAndNet = strings.Split(ixName, "||")
		iX, err = p.SearchIXByIxName(nameAndNet[0])
		if err != nil {
			return Netixlan{}, err
		}
	}

	v := url.Values{}

	if len(nameAndNet) > 1 && nameAndNet[1] != "" {
		for _, value := range iX.Data[0].IxlanSet {
			if value.Name == nameAndNet[1] {
				ixlanid = strconv.FormatInt(value.ID, 10)
				goto end
			}
		}
		return Netixlan{}, fmt.Errorf("attention, Net: %s given for %s, but not found", nameAndNet[1], nameAndNet[0])
	} else if len(iX.Data[0].IxlanSet) > 1 {
		return Netixlan{}, fmt.Errorf("there a multiple Nets to choose for %s, please specify in the ini-file\n", nameAndNet[0])
	} else {
		ixlanid = strconv.FormatInt(iX.Data[0].IxlanSet[0].ID, 10)
	}

end:
	v.Set("ixlan_id", ixlanid)
	if err = p.callAPI("/netixlan?"+v.Encode(), &apiResult); err != nil {
		return Netixlan{}, err
	}
	sort.Sort(netsIxLanDataSortedByASN(apiResult.Data))
	return apiResult, nil
}

func (p *peeringdb) GetIxLANByIxLanID(ixLanID int64) (apiResult IxLAN, err error) {
	if err = p.callAPI("/ixlan/"+strconv.FormatInt(ixLanID, 10), &apiResult); err != nil {
		return IxLAN{}, err
	}
	return apiResult, nil
}

func (p *peeringdb) ListIX() (apiResult Ix, err error) {
	if err = p.callAPI("/ix", &apiResult); err != nil {
		return Ix{}, err
	}
	return apiResult, nil
}

func (p *peeringdb) SearchIXByIxName(ixName string) (apiResult Ix, err error) {
	v := url.Values{}
	v.Set("name", ixName)
	if err = p.callAPI("/ix?"+v.Encode(), &apiResult); err != nil {
		return Ix{}, err
	}

	if len(apiResult.Data) == 0 {
		return Ix{}, fmt.Errorf("%s is not a valid ixName or was not found on peeringdb", ixName)
	}

	if err = p.callAPI("/ix/"+strconv.FormatInt(apiResult.Data[0].ID, 10), &apiResult); err != nil {
		return Ix{}, err

	}
	return apiResult, nil
}

func (p *peeringdb) SearchIXByIxId(ixId string) (apiResult Ix, err error) {
	v := url.Values{}
	v.Set("id", ixId)
	if err = p.callAPI("/ix?"+v.Encode(), &apiResult); err != nil {
		return Ix{}, err
	}

	if len(apiResult.Data) == 0 {
		return Ix{}, fmt.Errorf("%s is not a valid ixID or was not found on peeringdb", ixId)
	}

	if err = p.callAPI("/ix/"+strconv.FormatInt(apiResult.Data[0].ID, 10), &apiResult); err != nil {
		return Ix{}, err

	}
	return apiResult, nil
}

func (p *peeringdb) ListFaculty() (apiResult Ix, err error) {
	if err = p.callAPI("/fac", &apiResult); err != nil {
		return Ix{}, err
	}
	return apiResult, nil
}

func (p *peeringdb) SearchFacultyByFacName(facultyName string) (apiResult Ix, err error) {
	v := url.Values{}
	v.Set("name", facultyName)
	if err = p.callAPI("/fac?"+v.Encode(), &apiResult); err != nil {
		return Ix{}, err
	}
	return apiResult, nil
}

func (p *peeringdb) GetNetworkByAsN(asn int64) (apiResult Net, err error) {
	v := url.Values{}
	v.Set("asn", strconv.FormatInt(asn, 10))
	if err = p.callAPI("/net?"+v.Encode(), &apiResult); err != nil {
		return Net{}, err
	}
	return apiResult, nil
}
