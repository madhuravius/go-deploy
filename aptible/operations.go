package aptible

import (
	"strings"

	"github.com/aptible/go-deploy/client/operations"
	"github.com/aptible/go-deploy/models"
	"github.com/go-openapi/swag"
)

type Operation struct {
	ID            int64
	Type          string
	Handle        string
	Status        string
	EnvironmentID int64
	StackID       int64
	Certificate   string
	SSHUser       string
	SSHPty        bool
	CreatedAt     string
}

func (c *Client) CreateSSHPortalConnectionOperation(environmentID int64, publicKey string) (Operation, error) {
	environment, err := c.GetEnvironment(environmentID)
	if err != nil {
		return Operation{}, err
	}

	params := operations.NewPostOperationsOperationIDSSHPortalConnectionsParams().WithAppRequest(&models.AppRequest33{
		SSHPublicKey: &publicKey,
	})
	result, err := c.Client.Operations.PostOperationsOperationIDSSHPortalConnections(params, c.Token)
	if err != nil {
		return Operation{}, err
	}

	payload := result.GetPayload()
	operationID, err := GetIDFromHref(payload.Links.Operation.Href.String())
	if err != nil {
		return Operation{}, err
	}

	return Operation{
		ID:          operationID,
		StackID:     environment.StackID,
		CreatedAt:   swag.StringValue(payload.CreatedAt),
		Certificate: swag.StringValue(payload.SSHCertificateBody),
		SSHUser:     swag.StringValue(payload.SSHUser),
		SSHPty:      swag.BoolValue(payload.SSHPty),
	}, nil
}

func (c *Client) GetOperation(operationID int64) (Operation, error) {
	params := operations.NewGetOperationsIDParams().WithID(operationID)
	result, err := c.Client.Operations.GetOperationsID(params, c.Token)
	if err != nil {
		return Operation{}, err
	}
	payload := result.GetPayload()

	var handle string
	resourceLink := payload.Links.Resource.Href.String()
	resourceID, err := GetIDFromHref(resourceLink)
	if err != nil {
		return Operation{}, err
	}
	if strings.Contains(resourceLink, "apps") {
		app, err := c.GetApp(resourceID)
		if err != nil {
			return Operation{}, nil
		}
		handle = app.Handle
	} else if strings.Contains(resourceLink, "databases") {
		db, err := c.GetDatabase(resourceID)
		if err != nil {
			return Operation{}, nil
		}
		handle = db.Handle
	}

	var environmentID int64
	if payload.Links.Account != nil {
		environmentID, err = GetIDFromHref(payload.Links.Account.Href.String())
		if err != nil {
			return Operation{}, err
		}
	}

	return Operation{
		ID:            swag.Int64Value(payload.ID),
		Type:          swag.StringValue(payload.Type),
		Handle:        handle,
		EnvironmentID: environmentID,
		Status:        swag.StringValue(payload.Status),
	}, nil
}
