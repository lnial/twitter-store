package main

import (
	"encoding/json"
	"flag"
	"github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"regexp"
	"database/sql"
	_ "github.com/lib/pq"
)

type ApiConf struct {
	ConsumerKey       string `json:"consumer_key"`
	ConsumerSecret    string `json:"consumer_secret"`
	AccessToken       string `json:"access_token"`
	AccessTokenSecret string `json:"access_token_secret"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}


func main() {

	var apiConf ApiConf
	{
		apiConfPath := flag.String("conf", "config.json", "API Config File")
		flag.Parse()
		data, err_file := ioutil.ReadFile(*apiConfPath)
		check(err_file)
		err_json := json.Unmarshal(data, &apiConf)
		check(err_json)
	}

	anaconda.SetConsumerKey(apiConf.ConsumerKey)
	anaconda.SetConsumerSecret(apiConf.ConsumerSecret)
	api := anaconda.NewTwitterApi(apiConf.AccessToken, apiConf.AccessTokenSecret)

	db, err := sql.Open("postgres", "user=test password=test dbname=hoge sslmode=disable")
	check(err)

	twitterStream := api.PublicStreamSample(nil)
	for {
		// Channel
		x := <-twitterStream.C
		switch tweet := x.(type) {
		case anaconda.Tweet:
			check := match(tweet.Text)
			if check {
				// DBへ突っ込む
				if isRecord(db, tweet.Id, tweet.Text) {
					insertRecord(db, tweet.Id, tweet.Text)
				}

			}
		default:
		// pass
		}
	}
	db.Close()
}

func insertRecord(db *sql.DB, tweetId int64, tweet string) {
	stmt, err := db.Prepare("INSERT INTO tweet_media(tweet_id,url) VALUES($1,$2)")
	check(err)
	res, err := stmt.Exec(tweetId, "'" + tweet +"'")
	check(err)
	res.LastInsertId()
}

func isRecord(db *sql.DB, tweetId int64, tweet string) bool  {
	var cnt int
	err := db.QueryRow("SELECT COUNT(*) FROM tweet_media where tweet_id = $1 and url = $2;", string(tweetId), tweet).Scan(&cnt)
	check(err)
	return cnt == 0
}

func match(text string) bool {
	r := regexp.MustCompile(`#rezero`)
	return r.MatchString(text)
}
