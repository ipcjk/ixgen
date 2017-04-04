package peeringdb

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"
)

var matchIxRegex = `\/api\/ix\/(\d+)$`
var matchIxLanRegex = `\/api\/ixlan\/(\d+)$`

type handler struct {
	CacheDir string
	Data     []byte
	mutex    sync.Mutex
}

type netHandler struct {
	handler
	NetData Net
}

type Apiserver struct {
	AddrPort string
	CacheDir string
}

type getIX struct {
	handler
	match *regexp.Regexp
}

type getIXLan struct {
	handler
	match *regexp.Regexp
}

type getIXLans handler
type getIXes handler
type getNetIXLan handler
type getFac handler
type getNet netHandler
type getAll handler

// NewAPIServer returns a new Apiserver object, than can be
// started to answer to peeringdb-style api questions.
// It will take the ListenAddr and Port and also a source directory where
// to serve the object files from as arguments.
func NewAPIServer(Addrport, CacheDir string) *Apiserver {
	return &Apiserver{Addrport, CacheDir}
}

// RunAPIServer starts the created Apiserver
func (a *Apiserver) RunAPIServer() {
	r := http.NewServeMux()
	matchIx, _ := regexp.Compile(matchIxRegex)
	matchIxLan, _ := regexp.Compile(matchIxLanRegex)

	r.Handle("/api/ix", &getIXes{a.CacheDir, nil, sync.Mutex{}})
	r.Handle("/api/ix/", &getIX{handler{a.CacheDir, nil, sync.Mutex{}}, matchIx})
	r.Handle("/api/netixlan", &getNetIXLan{a.CacheDir, nil, sync.Mutex{}})
	r.Handle("/api/net", &getNet{handler{a.CacheDir, nil, sync.Mutex{}}, Net{}})
	r.Handle("/api/ixlan", &getIXLans{a.CacheDir, nil, sync.Mutex{}})
	r.Handle("/api/ixlan/", &getIXLan{handler{a.CacheDir, nil, sync.Mutex{}}, matchIxLan})

	listener, err := net.Listen("tcp", a.AddrPort)
	if err != nil {
		log.Fatalf("Cant spin up local api-service: %s", err)
	}
	a.AddrPort = listener.Addr().String()
	go http.Serve(listener, r)
}

func getJSON(r io.Reader, i interface{}) {
	err := json.NewDecoder(r).Decode(&i)
	if err != nil {
		log.Fatal("Problems decoding from json")
	}
}

func writeJSON(w io.Writer, i interface{}) {
	err := json.NewEncoder(w).Encode(&i)
	if err != nil {
		log.Fatal("Problems encoding ix from json")
	}
}

func readFile(fileName string) []byte {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	buf := new(bytes.Buffer)

	_, err = buf.ReadFrom(file)
	if err != nil {
		log.Fatalf("Cant read from file :%s", fileName)
	}
	return buf.Bytes()
}

func (h *getNet) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var apiResult Net
	params := r.URL.Query()

	h.mutex.Lock()
	/* Only the first request, will load the file into our structure */
	if len(h.Data) == 0 {
		h.Data = readFile(h.CacheDir + "/net")
		getJSON(bytes.NewBuffer(h.Data), &h.NetData)
	}
	h.mutex.Unlock()

	/* No params? Then write out all */
	if len(params) == 0 {
		apiResult = h.NetData
		goto end
	}

	/* Search for the network, can be optimized by loading things into a HASH */
	for k := range h.NetData.Data {
		if params["asn"][0] == strconv.FormatInt(h.NetData.Data[k].Asn, 10) {
			apiResult.Data = append(apiResult.Data, h.NetData.Data[k])
			break
		}
	}
end:
	writeJSON(w, &apiResult)
}

func (h *getNetIXLan) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var data Netixlan
	var apiResult Netixlan

	params := r.URL.Query()

	h.mutex.Lock()
	if len(h.Data) == 0 {
		h.Data = readFile(h.CacheDir + "/netixlan")
	}
	h.mutex.Unlock()

	getJSON(bytes.NewBuffer(h.Data), &data)

	if len(params) == 0 {
		apiResult = data
		goto end
	}

	for k := range data.Data {
		if params["ixlan_id"][0] == strconv.FormatInt(data.Data[k].IxlanID, 10) {
			apiResult.Data = append(apiResult.Data, data.Data[k])
		}
	}
end:
	writeJSON(w, &apiResult)
}

func (h *getIXLan) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var data IxLAN
	var apiResult IxLAN

	matches := h.match.FindStringSubmatch(r.RequestURI)
	h.mutex.Lock()
	if len(h.Data) == 0 {
		h.Data = readFile(h.CacheDir + "/ixlan")
	}
	h.mutex.Unlock()

	getJSON(bytes.NewBuffer(h.Data), &data)

	for k := range data.Data {
		if strconv.FormatInt(data.Data[k].ID, 10) == matches[1] {
			apiResult.Data = append(apiResult.Data, data.Data[k])
			break
		}
	}

	writeJSON(w, &apiResult)
}

func (h *getIXLans) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var data IxLAN
	var apiResult IxLAN

	params := r.URL.Query()

	h.mutex.Lock()
	if len(h.Data) == 0 {
		h.Data = readFile(h.CacheDir + "/ixlan")
	}
	h.mutex.Unlock()

	getJSON(bytes.NewBuffer(h.Data), &data)

	if len(params) == 0 {
		apiResult = data
		goto end
	}

	for k := range data.Data {
		for kp, kv := range params {
			if kp == "ix_id" && strconv.FormatInt(data.Data[k].IxID, 10) == kv[0] {
				apiResult.Data = append(apiResult.Data, data.Data[k])
			}
		}
	}

end:
	writeJSON(w, &apiResult)
}

func (h *getIX) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var data Ix
	var apiResult Ix

	matches := h.match.FindStringSubmatch(r.RequestURI)
	h.mutex.Lock()
	if len(h.Data) == 0 {
		h.Data = readFile(h.CacheDir + "/ix")
	}
	h.mutex.Unlock()

	getJSON(bytes.NewBuffer(h.Data), &data)

	for k := range data.Data {
		if strconv.FormatInt(data.Data[k].ID, 10) == matches[1] {
			apiResult.Data = append(apiResult.Data, data.Data[k])
			break
		}
	}

	/* populate ixlan_set */
	var ixLanData IxLAN
	ixData := readFile(h.CacheDir + "/ixlan")
	getJSON(bytes.NewBuffer(ixData), &ixLanData)

	for k := range ixLanData.Data {
		if strconv.FormatInt(ixLanData.Data[k].IxID, 10) == matches[1] {
			apiResult.Data[0].IxlanSet =
				append(apiResult.Data[0].IxlanSet, ixLanData.Data[k])
		}
	}

	writeJSON(w, &apiResult)
}

func (h *getIXes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var data Ix
	var apiResult Ix

	params := r.URL.Query()

	h.mutex.Lock()
	if len(h.Data) == 0 {
		h.Data = readFile(h.CacheDir + "/ix")
	}
	h.mutex.Unlock()

	getJSON(bytes.NewBuffer(h.Data), &data)

	if len(params) == 0 {
		apiResult = data
		goto end
	}

	for k := range data.Data {
		var hits = 0
		for kp, kv := range params {
			if kp == "name" && data.Data[k].Name == kv[0] {
				hits++
			} else if kp == "id" && strconv.FormatInt(data.Data[k].ID, 10) == kv[0] {
				hits++
			}
		}
		if hits == len(params) {
			apiResult.Data = append(apiResult.Data, data.Data[k])
			break
		}
	}

end:
	writeJSON(w, &apiResult)
}

func (h *getFac) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
