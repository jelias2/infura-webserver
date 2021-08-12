package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"time"

	"github.com/jelias2/infra-test/src/apis"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

type Handler struct {
	Log                        *zap.Logger
	Resty                      *resty.Client
	Mainnet_http_endpoint      string
	Mainnet_websocket_endpoint string
}

// Healthcheck will display test response to make sure the server is running
func (h *Handler) Healthcheck(w http.ResponseWriter, r *http.Request) {

	log.Info("Healthcheck ")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apis.Healthcheck{
		Status:   200,
		Message:  "Healthcheck response",
		Datetime: time.Now().String(),
	})
}

// Get ethblock number
func (h *Handler) GetBlockNumber(w http.ResponseWriter, r *http.Request) {
	h.Log.Info("Entered GetBlockNumber")
	body := `{"jsonrpc":"2.0","method":"eth_blockNumber","params": [],"id":1}`
	h.Log.Info("GetBlockNumber body", zap.String("Body", body))
	result := &apis.GetBlockNumberResponse{}
	resp, err := h.Resty.R().SetBody(body).
		SetResult(result).
		Post(h.Mainnet_http_endpoint)
	h.debugResponse("GetBlockNumber", resp, err)
	json.NewEncoder(w).Encode(result)
}

// Get GetGasPrice number
func (h *Handler) GetGasPrice(w http.ResponseWriter, r *http.Request) {

	h.Log.Info("Entered GetGasPrice")
	body := `{"jsonrpc":"2.0","method":"eth_gasPrice","params": [],"id":1}`
	result := &apis.GetBlockNumberResponse{}
	resp, err := h.Resty.R().SetBody(body).
		SetResult(result).
		Post(h.Mainnet_http_endpoint)
	h.debugResponse("GetBlockNumber", resp, err)
	json.NewEncoder(w).Encode(result)
}

// GetBlockByNumber
func (h *Handler) GetBlockByNumber(w http.ResponseWriter, r *http.Request) {

	var txdetails bool
	var err error

	w.Header().Set("Content-Type", "application/json")
	h.Log.Info("Entered GetBlockByNumber")
	reqBody, _ := ioutil.ReadAll(r.Body)
	var getBlockByNumberRequest apis.GetBlockByNumberRequest
	if err := json.Unmarshal(reqBody, &getBlockByNumberRequest); err != nil {
		h.Log.Error("Error unmarshalling GetBlockByNumberRequest", zap.Error(err))
	}

	txdetails, err = strconv.ParseBool(getBlockByNumberRequest.TxDetails)
	if getBlockByNumberRequest.Block == "" || err != nil {
		json.NewEncoder(w).Encode(apis.ErrorResponse{
			StatusCode: 400,
			Message:    "Malformed Request",
		})
	} else {
		body := fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["%s",%s],"id":1}`, getBlockByNumberRequest.Block, getBlockByNumberRequest.TxDetails)
		h.Log.Debug("GetBlockByNumber body", zap.String("Body", body))

		if txdetails {
			json.NewEncoder(w).Encode(h.TxDetailsResponse(body))
		} else {
			json.NewEncoder(w).Encode(h.NoTxDetailsResponse(body))
		}
	}
}

func (h *Handler) TxDetailsResponse(body string) *apis.GetBlockByNumberTxDetailsResponse {
	var err error
	var resp *resty.Response
	result := &apis.GetBlockByNumberTxDetailsResponse{}
	resp, err = h.Resty.R().SetBody(body).
		SetResult(result).
		Post(h.Mainnet_http_endpoint)
	if err != nil {
		h.Log.Error("Error", zap.Error(err))
	}
	h.debugResponse("GetBlockByNumber", resp, err)
	return result
}

func (h *Handler) NoTxDetailsResponse(body string) *apis.GetBlockByNumberNoTxDetailsResponse {
	var err error
	var resp *resty.Response
	result := &apis.GetBlockByNumberNoTxDetailsResponse{}
	resp, err = h.Resty.R().SetBody(body).
		SetResult(result).
		Post(h.Mainnet_http_endpoint)
	if err != nil {
		h.Log.Error("Error", zap.Error(err))
	}
	h.debugResponse("GetBlockByNumber", resp, err)
	return result
}

// GetBlockByNumber
func (h *Handler) GetTransactionByBlockNumberAndIndex(w http.ResponseWriter, r *http.Request) {

	var err error
	var resp *resty.Response

	w.Header().Set("Content-Type", "application/json")
	h.Log.Info("Entered GetTransactionByBlockNumberAndIndex")
	reqBody, _ := ioutil.ReadAll(r.Body)
	var getTxReq apis.GetTransactionByBlockNumberAndIndexRequest
	if err := json.Unmarshal(reqBody, &getTxReq); err != nil {
		h.Log.Error("Error unmarshalling GetBlockByNumberRequest", zap.Error(err))
	}

	// if getTxReq.Block == "" || getTxReq.Index == "" || !isHex(getTxReq.Block) || !isHex(getTxReq.Index) {
	if getTxReq.Block == "" || getTxReq.Index == "" {
		json.NewEncoder(w).Encode(apis.ErrorResponse{
			StatusCode: 400,
			Message:    "Malformed Request",
		})
	} else {
		body := fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getTransactionByBlockNumberAndIndex","params":["%s","%s"],"id":1}`, getTxReq.Block, getTxReq.Index)
		h.Log.Info("GetBlockByNumber body", zap.String("Body", body))

		result := &apis.GetTransactionByBlockNumberAndIndexResponse{}
		resp, err = h.Resty.R().SetBody(body).
			SetResult(result).
			Post(h.Mainnet_http_endpoint)
		if err != nil {
			h.Log.Error("Error", zap.Error(err))
		}
		h.debugResponse("GetTransactionByBlockNumberAndIndex", resp, err)

		json.NewEncoder(w).Encode(result)
	}
}

// func isHex(s string) bool {
// 	if _, err := hex.DecodeString(s); err == nil {
// 		return true
// 	}
// 	return false
// }

func (h *Handler) debugResponse(caller string, resp *resty.Response, err error) {
	h.Log.Info("Handling response from", zap.String("caller", caller))
	// Explore response object
	h.Log.Info("Response Info:",
		zap.Error(err),
		zap.Int("Status Code:", resp.StatusCode()),
		zap.String("Status     :", resp.Status()),
		zap.String("Proto      :", resp.Proto()),
		zap.Time("Received At:", resp.ReceivedAt()))
	// zap.String("Body :\n", string(resp.Body())))
}
