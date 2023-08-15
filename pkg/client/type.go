package client

import "time"

type ReplicationInfo struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Creator     string `json:"creator"`
	SrcRegistry struct {
		ID              int    `json:"id"`
		Name            string `json:"name"`
		Description     string `json:"description"`
		Type            string `json:"type"`
		URL             string `json:"url"`
		TokenServiceURL string `json:"token_service_url"`
		Credential      struct {
			Type         string `json:"type"`
			AccessKey    string `json:"access_key"`
			AccessSecret string `json:"access_secret"`
		} `json:"credential"`
		Insecure     bool      `json:"insecure"`
		Status       string    `json:"status"`
		CreationTime time.Time `json:"creation_time"`
		UpdateTime   time.Time `json:"update_time"`
	} `json:"src_registry"`
	DestRegistry struct {
		ID              int    `json:"id"`
		Name            string `json:"name"`
		Description     string `json:"description"`
		Type            string `json:"type"`
		URL             string `json:"url"`
		TokenServiceURL string `json:"token_service_url"`
		Credential      struct {
			Type         string `json:"type"`
			AccessKey    string `json:"access_key"`
			AccessSecret string `json:"access_secret"`
		} `json:"credential"`
		Insecure     bool      `json:"insecure"`
		Status       string    `json:"status"`
		CreationTime time.Time `json:"creation_time"`
		UpdateTime   time.Time `json:"update_time"`
	} `json:"dest_registry"`
	DestNamespace string `json:"dest_namespace"`
	Filters       []struct {
		Type  interface{} `json:"type"`
		Value interface{} `json:"value"`
	} `json:"filters"`
	Trigger struct {
		Type            string `json:"type"`
		TriggerSettings struct {
			Cron string `json:"cron"`
		} `json:"trigger_settings"`
	} `json:"trigger"`
	Deletion     bool      `json:"deletion"`
	Override     bool      `json:"override"`
	Enabled      bool      `json:"enabled"`
	CreationTime time.Time `json:"creation_time"`
	UpdateTime   time.Time `json:"update_time"`
}
