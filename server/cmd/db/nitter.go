package db

import (
	"errors"
	ps "teniditter-server/cmd/planetscale"
)

// Query Nittos by its name from DB; if the Nittos is not yet in the db this function will insert it.
func GetNittos(username string, depth ...int) (*NittosModel, error) {
	if len(depth) == 1 && depth[0] > 3 {
		return nil, errors.New("recursion emergency stop triggered")
	}

	db := ps.DBManager.Connect()
	if db == nil {
		return nil, ps.ErrDbNotFound
	}

	var result NittosModel
	err := db.QueryRow("SELECT * FROM Twittos WHERE username=?", username).Scan(&result.NittosID, &result.Username)
	if err != nil || result.Username != username {
		// Not in db --> insert it
		if _, err = db.Exec("INSERT INTO Twittos (username) VALUES (?);", username); err == nil {
			var depthVal int
			if len(depth) == 1 {
				depthVal = depth[0]
			}

			return GetNittos(username, depthVal+1) // refetch by controlled recursion
		}

		return nil, err
	}

	return &result, nil
}

func SearchNittos(username string) ([]NittosModel, error) {
	db := ps.DBManager.Connect()
	if db == nil {
		return nil, ps.ErrDbNotFound
	}

	rows, err := db.Query("SELECT * FROM Twittos WHERE username LIKE ?", "%"+username+"%")
	if err != nil {
		return nil, errors.New("error when fetching Subteddits")
	}
	defer rows.Close()

	var AllNittos []NittosModel
	for rows.Next() {
		var sub NittosModel
		if err := rows.Scan(&sub.NittosID, &sub.Username); err != nil {
			return AllNittos, nil
		}
		AllNittos = append(AllNittos, sub)
	}
	if err = rows.Err(); err != nil {
		return AllNittos, err
	}

	return AllNittos, nil
}
