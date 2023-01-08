package nasa

import (
	"log"
	"net/http"
	"time"
	"url-collector/internal/collector"
	"url-collector/internal/timeutils"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type GetPicturesQueryParams struct {
	From time.Time
	To   time.Time
}

func (qp *GetPicturesQueryParams) Validate() error {
	now := time.Now()
	if qp.From.After(now) || qp.To.After(now) {
		return errors.New("the given time cannot be future")
	}

	if qp.To.Before(qp.From) {
		return errors.New("'from' must be before 'to' param")
	}

	if !qp.To.IsZero() && qp.From.IsZero() {
		return errors.New("'from' must be defined if 'to' is set")
	}

	return nil
}

func (ctrl *Controller) GetPictures(c *gin.Context) {
	var (
		qp          GetPicturesQueryParams
		metadataSet collector.MetadataSet
		err         error
	)

	if qp, err = parseQueryParams(c); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("couldn't parse query params properly. They must be in 2006-01-02 format").Error()})
		return
	}

	if err = qp.Validate(); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.Wrap(err, "validation failed").Error()})
		return
	}

	if metadataSet, err = ctrl.nasa.Collect(c, qp.From, qp.To); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.New("failed to collect urls").Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"urls": metadataSet.ToURLList()})
}

func parseQueryParams(c *gin.Context) (qp GetPicturesQueryParams, err error) {
	if qp.From, err = parseQueryParam(c.Query("from")); err != nil {
		return
	}

	if qp.To, err = parseQueryParam(c.Query("to")); err != nil {
		return
	}

	if qp.To.IsZero() {
		now := time.Now()
		qp.To = now
		if qp.From.IsZero() {
			qp.From = now
		}
	}

	return
}

func parseQueryParam(v string) (t time.Time, err error) {
	if v == "" {
		return
	}
	return time.Parse(timeutils.DayFormat, v)
}
