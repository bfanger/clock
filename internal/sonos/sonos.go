/**
 *
 * Logic ported from node-sonos, scope: Listen for volume changes
 */
package sonos

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type Speaker struct {
	IP   net.IP
	Name string `xml:"device>displayName"`
	Room string `xml:"device>roomName"`
}

// Find Sonos speakers on the local network
func FindRoom(room string) (*Speaker, error) {
	search := []string{"M-SEARCH * HTTP/1.1", "HOST: 239.255.255.250:1900", "MAN: \"ssdp:discover\"", "MX: 1", "ST: urn:schemas-upnp-org:device:ZonePlayer:1"}
	addr, err := net.ResolveUDPAddr("udp", ":1900")
	if err != nil {
		return nil, err
	}
	connection, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}
	err = connection.SetReadDeadline(time.Now().Add(time.Second * 5))
	if err != nil {
		return nil, err
	}
	defer connection.Close()
	multicast, err := net.ResolveUDPAddr("udp", "239.255.255.250:1900")
	if err != nil {
		return nil, err
	}
	message := new(bytes.Buffer)
	message.WriteString(strings.Join(search, "\r\n"))
	_, err = connection.WriteTo(message.Bytes(), multicast)
	if err != nil {
		return nil, err
	}
	for {
		buf := make([]byte, 1024)
		_, addr, err := connection.ReadFromUDP(buf)
		if err != nil {
			return nil, err
		}
		res, err := http.Get("http://" + addr.IP.String() + ":1400/xml/device_description.xml")
		if err != nil {
			continue
		}
		defer res.Body.Close()
		speaker := &Speaker{
			IP: addr.IP,
		}
		err = xml.NewDecoder(res.Body).Decode(speaker)
		if err != nil {
			return nil, err
		}
		if speaker.Room == room {
			return speaker, nil
		}
	}
}

// Get the current volume of the speaker
func (s *Speaker) GetVolume() (int, error) {
	res, err := s.request("RenderingControl", "GetVolume", map[string]string{
		"InstanceID": "0",
		"Channel":    "Master",
	})
	if err != nil {
		return -1, err
	}
	response := &struct {
		CurrentVolume int `xml:"CurrentVolume"`
	}{
		CurrentVolume: -1,
	}
	err = xml.Unmarshal([]byte(res), response)
	return response.CurrentVolume, err
}

// Start http server and register callback for (volume) events
func (s *Speaker) HandleVolumeEvents(fn func(int)) error {
	ip, err := localIP()
	if err != nil {
		return err
	}
	request, err := http.NewRequest("SUBSCRIBE", "http://"+s.IP.String()+":1400/MediaRenderer/RenderingControl/Event", bytes.NewBufferString(""))
	if err != nil {
		return err
	}
	request.Header.Set("callback", "<http://"+ip.String()+":4444/notify>")
	request.Header.Set("NT", "upnp:event")
	request.Header.Set("Timeout", "Second-1800")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	res.Body.Close()

	return http.ListenAndServe(":4444", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if r.Method != "NOTIFY" || r.URL.Path != "/notify" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
			return
		}

		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
			return
		}
		volume, err := parseVolumeEvent(string(bytes))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
		}
		w.Write([]byte("OK"))
		fn(volume)
	}))
}

func (s *Speaker) request(control string, action string, variables map[string]string) (string, error) {
	url := "http://" + s.IP.String() + ":1400/MediaRenderer/" + control + "/Control"
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(requestBody(control, action, variables)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "text/xml; charset=utf8")
	req.Header.Set("SOAPAction", "urn:schemas-upnp-org:service:"+control+":1#"+action)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	result := string(bytes)
	tag := "u:" + action + "Response"
	start := strings.Index(result, "<"+tag)
	end := strings.LastIndex(result, "</"+tag+">") + len(tag) + 3
	return result[start:end], nil
}

func requestBody(control string, action string, variables map[string]string) string {
	request := "<?xml version=\"1.0\" ?>\n<s:Envelope s:encodingStyle=\"http://schemas.xmlsoap.org/soap/encoding/\" xmlns:s=\"http://schemas.xmlsoap.org/soap/envelope/\">\n  <s:Body>\n"
	request += "    <u:" + action + " xmlns:u=\"urn:schemas-upnp-org:service:" + control + ":1\">\n"
	for k, v := range variables {
		request += "      <" + k + ">" + v + "</" + k + ">\n"
	}
	request += "    </u:" + action + ">\n"
	request += "  </s:Body>\n</s:Envelope>"
	return request

}

func localIP() (net.IP, error) {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addresses {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP, nil
			}
		}
	}
	return nil, errors.New("no IP address found")

}

func parseVolumeEvent(request string) (int, error) {

	start := strings.Index(request, "<LastChange>") + 12
	end := strings.LastIndex(request, "</LastChange>")
	if end == -1 {
		return -1, errors.New("could not find <LastChange> tag")
	}
	decoded := html.UnescapeString(request[start:end])

	start = strings.Index(decoded, `<Volume channel="Master" val="`)
	if start == -1 {
		return 0, errors.New("could not find <Volume> tag")
	}
	decoded = decoded[start+30:]
	end = strings.Index(decoded, `"`)
	if end == -1 {
		return 0, errors.New("could not find closing quote")
	}
	return strconv.Atoi(decoded[:end])
}
