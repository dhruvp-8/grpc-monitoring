package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/pusher/pusher-http-go"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const defaultSleepTime = time.Second * 2
const bucket = "grpc-monitoring-config-data"
const item = "env"

func main() {
	httpPort := flag.Int("http.port", 4000, "HTTP Port to run server on")
	mongoDSN := flag.String("mongo.dsn", "localhost:27017", "DSN for mongoDB server")

	flag.Parse()

	loadEnvFileFromS3(bucket, item)

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	appID := os.Getenv("PUSHER_APP_ID")
	appKey := os.Getenv("PUSHER_APP_KEY")
	appSecret := os.Getenv("PUSHER_APP_SECRET")
	appCluster := os.Getenv("PUSHER_APP_CLUSTER")
	appIsSecure := os.Getenv("PUSHER_APP_SECURE")

	var isSecure bool
	if appIsSecure == "1" {
		isSecure = true
	}

	client := &pusher.Client{
		AppID:   appID,
		Key:     appKey,
		Secret:  appSecret,
		Cluster: appCluster,
		Secure:  isSecure,
		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}

	mux := chi.NewRouter()

	log.Println("Connecting to MongoDB")
	m, err := newMongo(*mongoDSN)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully connected to MongoDB")

	mux.Use(analyticsMiddleware(m, client))

	var once sync.Once
	var t *template.Template

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "public")
	fileServer(mux, "/public", http.Dir(filesDir))

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {

		once.Do(func() {
			tem, err := template.ParseFiles("public/index.html")
			if err != nil {
				log.Fatal(err)
			}

			t = tem.Lookup("index.html")
		})

		t.Execute(w, nil)
	})

	mux.Get("/api/analytics", analyticsAPI(m))
	mux.Get("/wait/{seconds}", waitHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), mux))
}

func loadEnvFileFromS3(bucket string, item string) {

	file, err := os.Create("." + item)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", item, err)
	}

	defer file.Close()

	// Initialize a session in us-west-1 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1")},
	)

	downloader := s3manager.NewDownloader(sess)

	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(item),
		})
	if err != nil {
		exitErrorf("Unable to download item %q, %v", item, err)
	}

	fmt.Println("Successfully Downloaded", file.Name(), numBytes, "bytes")
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}

	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func analyticsAPI(m mongo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := m.getAggregatedAnalytics()
		if err != nil {
			log.Println(err)

			json.NewEncoder(w).Encode(&struct {
				Message   string `json:"message"`
				TimeStamp int64  `json:"timestamp"`
			}{
				Message:   "An error occurred while fetching analytics data",
				TimeStamp: time.Now().Unix(),
			})

			return
		}

		// Handling Cors
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}

func analyticsMiddleware(m mongo, client *pusher.Client) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			startTime := time.Now()

			defer func() {

				if strings.HasPrefix(r.URL.String(), "/wait") {

					data := requestAnalytics{
						URL:         r.URL.String(),
						Method:      r.Method,
						RequestTime: time.Now().Unix() - startTime.Unix(),
						Day:         startTime.Weekday().String(),
						Hour:        startTime.Hour(),
					}

					if err := m.Write(data); err != nil {
						log.Println(err)
					}

					aggregatedData, err := m.getAggregatedAnalytics()
					if err == nil {
						client.Trigger("grpc-monitoring", "data", aggregatedData)
					}
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func waitHandler(w http.ResponseWriter, r *http.Request) {
	var sleepTime = defaultSleepTime

	secondsToSleep := chi.URLParam(r, "seconds")
	n, err := strconv.Atoi(secondsToSleep)
	if err == nil && n >= 2 {
		sleepTime = time.Duration(n) * time.Second
	} else {
		n = 2
	}

	log.Printf("Sleeping for %d seconds", n)
	time.Sleep(sleepTime)
	w.Write([]byte(`Done`))
}
