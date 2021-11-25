package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/justteddy/wallet/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type reportRequest struct {
	FromDate      string              `json:"from_date"`
	ToDate        string              `json:"to_date"`
	OperationType types.OperationType `json:"operation_type"`
}

func (h *Handler) HandleReport(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	format := params.ByName("format")
	walletID := params.ByName("wallet")

	var reportReq reportRequest
	if err := json.NewDecoder(r.Body).Decode(&reportReq); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, errors.Wrap(err, "decode request"))
		return
	}

	fromDate, toDate, err := h.validateReportRequest(format, walletID, reportReq)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	ops, err := h.s.Operations(r.Context(), types.WalletID(walletID), reportReq.OperationType, fromDate, toDate)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrap(err, "fetch operations"))
		return
	}

	data, err := h.e.Export(types.ExportFormat(format), types.TransformDBToExportOperation(ops))
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, errors.Wrap(err, "export operations"))
		return
	}

	if format == string(types.ExportFormatJSON) {
		w.Header().Add("Content-Type", "application/json")
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		log.WithError(err).Error("failed to write successful response")
	}
}

func (h *Handler) validateReportRequest(format, walletID string, reportReq reportRequest) (time.Time, time.Time, error) {
	var fromDate, toDate time.Time

	if walletID == "" {
		return fromDate, toDate, errors.New("empty wallet id")
	}
	if format == "" {
		return fromDate, toDate, errors.New("empty format")
	}

	if _, ok := types.AllExportFormats[types.ExportFormat(format)]; !ok {
		return fromDate, toDate, errors.New("unexpected export format")
	}

	if reportReq.OperationType != "" {
		if _, ok := types.AllOperationTypes[reportReq.OperationType]; !ok {
			return fromDate, toDate, errors.New("unexpected operation type")
		}
	}

	var err error
	if reportReq.FromDate != "" {
		if fromDate, err = time.Parse(types.DateLayout, reportReq.FromDate); err != nil {
			return fromDate, toDate, errors.Wrap(err, "invalid date format in from_date, should be YYYY-MM-DD")
		}
	}

	if reportReq.ToDate != "" {
		if toDate, err = time.Parse(types.DateLayout, reportReq.ToDate); err != nil {
			return fromDate, toDate, errors.Wrap(err, "invalid date format in to_date, should be YYYY-MM-DD")
		}
	}

	if fromDate.After(toDate) {
		return fromDate, toDate, errors.New("from_date is greater than to_date")
	}

	return fromDate, toDate, nil
}
