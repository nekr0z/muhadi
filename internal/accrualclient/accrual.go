package accrualclient

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/nekr0z/muhadi/internal/ctxlog"
	"github.com/nekr0z/muhadi/internal/reconciler"
	"go.uber.org/zap"
)

const defaultRetries = 3

const endpoint = "/api/orders/"

var _ reconciler.Accrual = &Client{}

type Client struct {
	c        *http.Client
	endpoint string
}

func New(address string) *Client {
	return &Client{
		c: resty.New().
			SetRetryCount(defaultRetries).
			GetClient(),
		endpoint: address + endpoint,
	}
}

func (c *Client) Status(ctx context.Context, orderID int) (float64, error) {
	ep := c.endpoint + strconv.Itoa(orderID)

	ctxlog.Debug(ctx, "Requesting accrual status", zap.Int("order_id", orderID), zap.String("endpoint", ep))
	resp, err := c.c.Get(ep)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		ctxlog.Warn(ctx, "Accrual overload")
		return parseOverload(ctx, resp.Header)
	}

	if resp.StatusCode == http.StatusOK {
		return parseAccrual(resp.Body)
	}

	return 0, reconciler.ErrAccrualNotReady
}

func parseOverload(ctx context.Context, h http.Header) (float64, error) {
	backoff, err := strconv.Atoi(h.Get("Retry-After"))
	if err != nil {
		ctxlog.Warn(ctx, "Can't parse Retry-After header", zap.String("header", h.Get("Retry-After")))
		backoff = 0
	}

	return 0, &reconciler.ErrAccrualOverload{Duration: time.Second * time.Duration(backoff)}
}

func parseAccrual(r io.Reader) (float64, error) {
	decoder := json.NewDecoder(r)

	var resp accrualResponse

	if err := decoder.Decode(&resp); err != nil {
		return 0, err
	}

	switch resp.Status {
	case "PROCESSED":
		return resp.Accrual, nil
	case "INVALID":
		return 0, reconciler.ErrAccrualRejected
	default:
		return 0, reconciler.ErrAccrualNotReady
	}
}

type accrualResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}
