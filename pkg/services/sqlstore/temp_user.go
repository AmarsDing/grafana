package sqlstore

import (
	"time"

	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/models"
)

func init() {
	bus.AddHandler("sql", CreateTempUser)
	bus.AddHandler("sql", GetTempUsersQuery)
	bus.AddHandler("sql", UpdateTempUserStatus)
	bus.AddHandler("sql", GetTempUserByCode)
	bus.AddHandler("sql", UpdateTempUserWithEmailSent)
	bus.AddHandler("sql", ExpireOldUserInvites)
}

func UpdateTempUserStatus(cmd *models.UpdateTempUserStatusCommand) error {
	return inTransaction(func(sess *DBSession) error {
		var rawSql = "UPDATE temp_user SET status=? WHERE code=?"
		_, err := sess.Exec(rawSql, string(cmd.Status), cmd.Code)
		return err
	})
}

func CreateTempUser(cmd *models.CreateTempUserCommand) error {
	return inTransaction(func(sess *DBSession) error {
		// create user
		user := &models.TempUser{
			Email:           cmd.Email,
			Name:            cmd.Name,
			OrgId:           cmd.OrgId,
			Code:            cmd.Code,
			Role:            cmd.Role,
			Status:          cmd.Status,
			RemoteAddr:      cmd.RemoteAddr,
			InvitedByUserId: cmd.InvitedByUserId,
			EmailSentOn:     time.Now(),
			Created:         time.Now(),
			Updated:         time.Now(),
		}

		if _, err := sess.Insert(user); err != nil {
			return err
		}

		cmd.Result = user
		return nil
	})
}

func UpdateTempUserWithEmailSent(cmd *models.UpdateTempUserWithEmailSentCommand) error {
	return inTransaction(func(sess *DBSession) error {
		user := &models.TempUser{
			EmailSent:   true,
			EmailSentOn: time.Now(),
		}

		_, err := sess.Where("code = ?", cmd.Code).Cols("email_sent", "email_sent_on").Update(user)

		return err
	})
}

func GetTempUsersQuery(query *models.GetTempUsersQuery) error {
	rawSql := `SELECT
	                tu.id             as id,
	                tu.org_id         as org_id,
	                tu.email          as email,
									tu.name           as name,
									tu.role           as role,
									tu.code           as code,
									tu.status         as status,
									tu.email_sent     as email_sent,
									tu.email_sent_on  as email_sent_on,
									tu.created				as created,
									u.login						as invited_by_login,
									u.name						as invited_by_name,
									u.email						as invited_by_email
	                FROM ` + dialect.Quote("temp_user") + ` as tu
									LEFT OUTER JOIN ` + dialect.Quote("user") + ` as u on u.id = tu.invited_by_user_id
									WHERE tu.status=?`
	params := []interface{}{string(query.Status)}

	if query.OrgId > 0 {
		rawSql += ` AND tu.org_id=?`
		params = append(params, query.OrgId)
	}

	if query.Email != "" {
		rawSql += ` AND tu.email=?`
		params = append(params, query.Email)
	}

	rawSql += " ORDER BY tu.created desc"

	query.Result = make([]*models.TempUserDTO, 0)
	sess := x.SQL(rawSql, params...)
	err := sess.Find(&query.Result)
	return err
}

func GetTempUserByCode(query *models.GetTempUserByCodeQuery) error {
	var rawSql = `SELECT
	                tu.id             as id,
	                tu.org_id         as org_id,
	                tu.email          as email,
									tu.name           as name,
									tu.role           as role,
									tu.code           as code,
									tu.status         as status,
									tu.email_sent     as email_sent,
									tu.email_sent_on  as email_sent_on,
									tu.created				as created,
									u.login						as invited_by_login,
									u.name						as invited_by_name,
									u.email						as invited_by_email
	                FROM ` + dialect.Quote("temp_user") + ` as tu
									LEFT OUTER JOIN ` + dialect.Quote("user") + ` as u on u.id = tu.invited_by_user_id
	                WHERE tu.code=?`

	var tempUser models.TempUserDTO
	sess := x.SQL(rawSql, query.Code)
	has, err := sess.Get(&tempUser)

	if err != nil {
		return err
	} else if !has {
		return models.ErrTempUserNotFound
	}

	query.Result = &tempUser
	return err
}

func ExpireOldUserInvites(cmd *models.ExpireTempUsersCommand) error {
	return inTransaction(func(sess *DBSession) error {
		var rawSql = "UPDATE temp_user SET status = ? WHERE created <= ? AND status != ?"
		expiredStatus := string(models.TmpUserExpired)

		if result, err := sess.Exec(rawSql, expiredStatus, cmd.OlderThan, expiredStatus); err != nil {
			return err
		} else if cmd.ExpiredInvites, err = result.RowsAffected(); err != nil {
			return err
		}
		return nil
	})
}
