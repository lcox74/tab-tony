package zerotier

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var ztControllerToken = os.Getenv("ZEROTIER_ACCESS_KEY")
var ztControllerNodeId = os.Getenv("ZEROTIER_NODE_ID")

func makeRequest(method string, endpoint string, body io.Reader) ([]byte, int, error) {
	token := ztControllerToken

	// Create Request
	client := &http.Client{}
	req, err := http.NewRequest(method, fmt.Sprintf("http://localhost:9993/%s", endpoint), body)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Add("X-ZT1-AUTH", token)
	req.Header.Add("Content-Type", "application/json")

	// Send Request
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, 0, err
	}
	defer res.Body.Close()

	// Extract Response
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, 0, err
	}

	return resBody, res.StatusCode, nil
}

func GetStatus() ([]byte, int, error) {
	return makeRequest("GET", "status", nil)
}

func CreateNetwork() (*NetworkAPI, error) {

	// Send Request to Create new Network
	body, status, err := makeRequest("POST", fmt.Sprintf("controller/network/%s______", ztControllerNodeId), strings.NewReader("{ }"))
	if err != nil {
		return nil, err
	}

	// Check if valid response
	if status != 200 {
		return nil, fmt.Errorf("failed to create new network (%d): %s", status, string(body))
	}

	// Unmarshal response into ZeroTier Network Structure
	zt_net := &NetworkAPI{}
	json.Unmarshal(body, zt_net)
	return zt_net, nil
}

func SetNetworkIPRange(nwid string, start string, end string, subnet string) error {
	// Send Request to Create new Network
	body, status, err := makeRequest("POST", fmt.Sprintf("controller/network/%s", nwid), strings.NewReader(fmt.Sprintf(`{"ipAssignmentPools": [{"ipRangeStart": "%s", "ipRangeEnd": "%s"}], "routes": [{"target": "%s", "via": null}], "v4AssignMode": "zt", "private": true }`, start, end, subnet)))
	if err != nil {
		return err
	}

	// Check if valid response
	if status != 200 {
		return fmt.Errorf("failed to set network ip range (%d): %s", status, string(body))
	}

	return nil
}

func DeleteNetwork(nwid string) error {
	// Send Request to Delete Network
	body, status, err := makeRequest("DELETE", fmt.Sprintf("controller/network/%s", nwid), nil)
	if err != nil {
		return err
	}

	// Check if valid response
	if status != 200 {
		return fmt.Errorf("failed to delete network (%d): %s", status, string(body))
	}

	return nil
}

func GetNetworks() ([]NetworkAPI, error) {
	// Send Request to Get Network
	body, status, err := makeRequest("GET", "controller/network", strings.NewReader(""))
	if err != nil {
		return nil, err
	}

	// Check if valid response
	if status != 200 {
		return nil, fmt.Errorf("failed to get networks (%d): %s", status, string(body))
	}

	// Unmarshal response into ZeroTier Network Structure
	zt_nets := []NetworkAPI{}
	json.Unmarshal(body, &zt_nets)
	return zt_nets, nil
}

func GetNetwork(nwid string) (*NetworkAPI, error) {
	// Send Request to Get Network
	body, status, err := makeRequest("GET", fmt.Sprintf("controller/network/%s", nwid), nil)
	if err != nil {
		return nil, err
	}

	// Check if valid response
	if status != 200 {
		return nil, fmt.Errorf("failed to get network (%d): %s", status, string(body))
	}

	// Unmarshal response into ZeroTier Network Structure
	zt_net := &NetworkAPI{}
	json.Unmarshal(body, zt_net)
	return zt_net, nil
}

func GetNetworkMembers(nwid string) ([]NetworkMemberAPI, error) {
	// Send Request to Get Network
	body, status, err := makeRequest("GET", fmt.Sprintf("controller/network/%s/member", nwid), nil)
	if err != nil {
		return nil, err
	}

	// Check if valid response
	if status != 200 {
		return nil, fmt.Errorf("failed to get network members (%d): %s", status, string(body))
	}

	members := make(map[string]int)
	json.Unmarshal(body, &members)

	zt_net_members := []NetworkMemberAPI{}
	for k := range members {
		zt_net_member, err := GetNetworkMemberInfo(nwid, k)
		if err != nil {
			return nil, err
		}
		zt_net_members = append(zt_net_members, *zt_net_member)
	}

	return zt_net_members, nil
}

func GetNetworkMemberInfo(nwid string, memberid string) (*NetworkMemberAPI, error) {
	// Send Request to Get Network
	body, status, err := makeRequest("GET", fmt.Sprintf("controller/network/%s/member/%s", nwid, memberid), nil)
	if err != nil {
		return nil, err
	}

	// Check if valid response
	if status != 200 {
		return nil, fmt.Errorf("failed to get network member (%d): %s", status, string(body))
	}

	// Unmarshal response into ZeroTier Member Structure
	zt_net_member := &NetworkMemberAPI{}
	json.Unmarshal(body, zt_net_member)
	return zt_net_member, nil
}

func AuthoriseMember(nwid string, memberid string) error {
	// Send Request to Authorise Member
	body, status, err := makeRequest("POST", fmt.Sprintf("controller/network/%s/member/%s", nwid, memberid), strings.NewReader("{ \"authorized\": true }"))
	if err != nil {
		return err
	}

	// Check if valid response
	if status != 200 {
		return fmt.Errorf("failed to authorise member (%d): %s", status, string(body))
	}

	return nil
}

func DeauthoriseMember(nwid string, memberid string) error {
	// Send Request to Authorise Member
	body, status, err := makeRequest("POST", fmt.Sprintf("controller/network/%s/member/%s", nwid, memberid), strings.NewReader("{ \"authorized\": false }"))
	if err != nil {
		return err
	}

	// Check if valid response
	if status != 200 {
		return fmt.Errorf("failed to deauthorise member (%d): %s", status, string(body))
	}

	return nil
}

func DeleteMember(nwid string, memberid string) error {
	// Send Request to Delete Member
	body, status, err := makeRequest("DELETE", fmt.Sprintf("controller/network/%s/member/%s", nwid, memberid), nil)
	if err != nil {
		return err
	}

	// Check if valid response
	if status != 200 {
		return fmt.Errorf("failed to delete member (%d): %s", status, string(body))
	}

	return nil
}
