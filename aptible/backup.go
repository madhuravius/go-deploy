package aptible

import (
	"github.com/aptible/go-deploy/client/operations"
	"github.com/go-openapi/swag"
)

type Backup struct {
	ID        int64
	CreatedAt string
	Region    string
	Manual    bool
	Copy      interface{}
}

func (c *Client) GetBackups(databaseId int64) ([]Backup, error) {
	params := operations.NewGetDatabasesDatabaseIDBackupsParams().WithDatabaseID(databaseId)
	result, err := c.Client.Operations.GetDatabasesDatabaseIDBackups(params, c.Token)
	if err != nil {
		return nil, err
	}
	var backups []Backup
	for _, backup := range result.GetPayload().Embedded.Backups {
		backups = append(backups, Backup{
			ID:        backup.ID,
			Region:    backup.AwsRegion,
			CreatedAt: backup.CreatedAt,
			Manual:    swag.BoolValue(backup.Manual),
			Copy:      backup.Embedded.CopiedFrom,
		})
	}

	return backups, nil
}
