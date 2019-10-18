package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"html/template"
	"log"
	"net/http"
	"time"
)

const baseTemplate = "layout.tmpl"

type LogEntryGet struct {
	Tags []string `json:"tags"`
	Text string `json:"text"`
	Time string `json:"time"`
}

type LogEntryPost struct {
	Tags []string `json:"tags"`
	Text string `json:"text"`
}

type LogEntry struct {
	Tags []string `firestore:"tags"`
	Text string `firestore:"text"`
	Time time.Time `firestore:"time"`
}

func main() {
	ctx := context.Background()

	opt := option.WithCredentialsFile("farzaparty-firebase-adminsdk-3trv0-c45c05f0cb.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatal(err.Error())
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	router := gin.Default()
	router.LoadHTMLGlob("./templates/**/*")
	router.Static("/static", "./static")
	logsCollection := client.Collection("logs")

	router.GET("/", func(c *gin.Context) {
		t, _ := template.ParseFiles("templates/layout/layout.tmpl", "templates/home/home.tmpl")
		t.ExecuteTemplate(c.Writer, "layout",  gin.H{
			"title": "Main website",
		})
	})

	router.GET("/past", func(c *gin.Context) {
		t, _ := template.ParseFiles("templates/layout/layout.tmpl", "templates/past/past.tmpl")
		t.ExecuteTemplate(c.Writer, "layout",  "")
	})

	router.GET("/log", func(c *gin.Context) {
		t, _ := template.ParseFiles("templates/layout/layout.tmpl", "templates/log/log.tmpl")
		q := logsCollection.OrderBy("time", firestore.Desc).Limit(100)
		iter := q.Documents(ctx)

		logs := make([]*LogEntryGet, 0)

		defer iter.Stop()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
				return
			}
			logEntry := &LogEntry{}
			if err := doc.DataTo(logEntry); err != nil {
				c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
				return
			}

			t := logEntry.Time
			loc, err := time.LoadLocation("America/New_York")
			if err := doc.DataTo(logEntry); err != nil {
				c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
				return
			}
			t = t.In(loc)
			logs = append(logs, &LogEntryGet{
				Text: logEntry.Text,
				Tags: logEntry.Tags,
				Time: t.Format(("Mon 3:04PM Jan _2 2006")),
			})
		}
		fmt.Println(logs)
		t.ExecuteTemplate(c.Writer, "layout", logs)
	})

	router.POST("/log", func(c *gin.Context) {
		entryPost := &LogEntryPost{}
		if err := c.ShouldBindJSON(entryPost); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		entry := &LogEntry{
			Time: time.Now(),
			Text: entryPost.Text,
			Tags: entryPost.Tags,
		}
		_, _, err := logsCollection.Add(ctx, entry)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		//dr.Update(ctx, []firestore.Update{{Path: "capital", Value: "Sacramento"}})
		c.JSON(http.StatusOK, entryPost)
	})

	router.Run(":80")
}