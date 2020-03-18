package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/sethvargo/go-diceware/diceware"
)

type phoneNumber struct {
	Country string
	Number  string
}

type jitsiMeeting struct {
	Name  string
	ID    string
	Phone []phoneNumber
}

func (j *jitsiMeeting) getJitsiName() error {
	list, err := diceware.Generate(4)
	if err != nil {
		return err
	}
	j.Name = strings.Join(list, "-")
	return nil
}

func (j *jitsiMeeting) getURL() string {
	return fmt.Sprintf("https://meet.jit.si/%s", j.Name)
}

func (j *jitsiMeeting) getPIN() string {
	if len(j.ID) == 10 {
		return fmt.Sprintf("%s %s %s#", j.ID[0:4], j.ID[4:8], j.ID[8:10])
	}
	return fmt.Sprintf("%s#", j.ID)
}

func (j *jitsiMeeting) getJitsiID() error {
	type jitsiMeetingHTTPResponse struct {
		Message    string `json:"message,omitempty"`
		ID         int64  `json:"id,omitempty"`
		Conference string `json:"conference,omitempty"`
	}

	queryURL := fmt.Sprintf("https://api.jitsi.net/conferenceMapper?conference=%s@conference.meet.jit.si", url.QueryEscape(j.Name))
	resp, err := http.Get(queryURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var jR jitsiMeetingHTTPResponse
	if err := json.Unmarshal(respBody, &jR); err != nil {
		return err
	}
	j.ID = strconv.FormatInt(jR.ID, 10)
	return nil
}

func (j *jitsiMeeting) getJitsiNumbers() error {
	type jitsiPhoneHTTPResponse struct {
		Message string              `json:"message,omitempty"`
		Numbers map[string][]string `json:"numbers,omitempty"`
		Enabled bool                `json:"numbersEnabled,omitempty"`
	}
	queryURL := fmt.Sprintf("https://api.jitsi.net/phoneNumberList?conference=%s@conference.meet.jit.si", url.QueryEscape(j.Name))
	resp, err := http.Get(queryURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var jR jitsiPhoneHTTPResponse
	if err := json.Unmarshal(respBody, &jR); err != nil {
		return err
	}
	for key, value := range jR.Numbers {
		j.Phone = append(j.Phone, phoneNumber{
			Country: key,
			Number:  value[0],
		})
	}
	// then sort them alphabetically
	sort.Slice(j.Phone, func(a, b int) bool { return j.Phone[a].Country < j.Phone[b].Country })
	return nil
}

func newJitsiMeeting() (jitsiMeeting, error) {
	result := jitsiMeeting{}
	result.getJitsiName()
	if err := result.getJitsiID(); err != nil {
		return result, err
	}
	if err := result.getJitsiNumbers(); err != nil {
		return result, err
	}
	return result, nil
}
