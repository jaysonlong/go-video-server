package dbops

import (
	"log"
)

func AddVideoDeletionRecord(videoId string) error {
	stmtIns, err := dbConn.Prepare("INSERT INTO video_del_rec (video_id) VALUES (?)")
	if err != nil {
		return err
	}

	_, err = stmtIns.Exec(videoId)
	defer stmtIns.Close()
	if err != nil {
		log.Printf("AddVideoDeletionRecord error: %v", err)
		return err
	}

	return nil
}
