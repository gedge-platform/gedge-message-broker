package controller

import (
	"backend/model"
	"backend/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type MessageInfo struct {
	Msg       map[string]interface{} `json:"msg"`
	Protocol  string                 `json:"protocol"`
	Channel   string                 `json:"channel"`
	ClientId  string                 `json:"clientId"`
	Metadata  string                 `json:"metadata"`
	ClusterId int                    `json:"clusterId"`
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	db := model.DBConn()

	var messages []MessageInfo
	rows, err := db.Query("SELECT protocol, channel, clientId, metadata ,msg FROM Messages WHERE ClusterId = ? ORDER BY id DESC", id)
	if util.CheckHttpError(w, err, "Check DB") {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var m MessageInfo
		var msg []byte
		err = rows.Scan(&m.Protocol, &m.Channel, &m.ClientId, &m.Metadata, &msg)
		if util.CheckHttpError(w, err, "Check DB") {
			return
		}

		err = json.Unmarshal(msg, &m.Msg)
		if util.CheckHttpError(w, err, "JSON Unmarshal") {
			return
		}

		messages = append(messages, m)
	}
	err = rows.Err()
	if util.CheckHttpError(w, err, "Check DB") {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func SendMessage(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if util.CheckHttpError(w, err, "Check DB") {
		return
	}

	var mInfo MessageInfo
	err = json.Unmarshal(body, &mInfo)
	if util.CheckHttpError(w, err, "Check DB") {
		return
	}

	db := model.DBConn()
	// defer db.Close()
	rows, err := db.Query("SELECT id, name, ip, rpcPort, state, date FROM Clusters WHERE id = ?", mInfo.ClusterId)
	if util.CheckHttpError(w, err, "Check DB") {
		return
	}
	var c Cluster
	for rows.Next() {
		scanErr := rows.Scan(&c.ID, &c.Name, &c.IP, &c.RPCPort, &c.State, &c.Date)
		if util.CheckHttpError(w, scanErr, "Check DB attribute") {
			return
		}
		break
	}
	rows.Close()
	jsonData, err := json.Marshal(mInfo.Msg)
	if util.CheckHttpError(w, err, "Check JSON Message") {
		return
	}

	switch mInfo.Protocol {
	case "pub/sub":
		err := pubMsg(c, mInfo, jsonData)
		if util.CheckHttpError(w, err, "Check pub/sub request") {
			return
		}
	case "queue":
		err := queueMsg(c, mInfo, jsonData)
		if util.CheckHttpError(w, err, "Check queue request") {
			return
		}
	case "query":
		err := queryMsg(c, mInfo, jsonData)
		if util.CheckHttpError(w, err, "Check query request") {
			return
		}
	case "command":
		err := commandMsg(c, mInfo, jsonData)
		if util.CheckHttpError(w, err, "Check command request") {
			return
		}
	}

	insert, dbInsterErr := db.Query("INSERT INTO Messages (protocol, channel, clientId, metadata, msg, ClusterId) VALUES (?, ?, ?, ?, ?, ?)",
		mInfo.Protocol, mInfo.Channel, mInfo.ClientId, mInfo.Metadata, jsonData, mInfo.ClusterId)
	if util.CheckHttpError(w, dbInsterErr, "Check data") {
		return
	}
	defer insert.Close()
	date := time.Now().Format("2006-01-02 15:04:05")
	_, updateErr := db.Exec("UPDATE Clusters SET date = ? WHERE id = ?", date, c.ID)
	if util.CheckHttpError(w, updateErr, "DB Update Failed") {
		return
	}
	fmt.Fprintf(w, "Data Inserted Successfully")
}

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w)
	fmt.Println(r)
	fmt.Fprintln(w, "AddMultiClusterList!")
}

func pubMsg(c Cluster, mInfo MessageInfo, jsonData []byte) error {
	return nil
}

func queueMsg(c Cluster, mInfo MessageInfo, jsonData []byte) error {
	return nil
}

func queryMsg(c Cluster, mInfo MessageInfo, jsonData []byte) error {
	return nil
}

func commandMsg(c Cluster, mInfo MessageInfo, jsonData []byte) error {
	return nil
}
