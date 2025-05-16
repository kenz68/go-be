package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Job struct {
	URL string
}
type Result struct {
	URL  string
	Data string
}

var (
	jobs    = make(chan Job, 10)
	results = make(chan Result, 10)
)

func main() {
	urls := []string{
		"https://go.dev/doc/effective_go#concurrency",
		"https://golangbot.com/buffered-channels-worker-pools/",
		"https://chatgpt.com/c/33681582-e751-42e5-a5f9-98a6f5cdd7f0",
	}

	startTime := time.Now()

	noOfworkers := 2

	go allocateJobs(urls)
	done := make(chan bool)
	go resultsContent(done)
	createWorkerPool(noOfworkers)
	<-done

	endTime := time.Now()
	diff := endTime.Sub(startTime)
	println("Total time taken is: ", diff.Seconds(), " seconds")
}

func createWorkerPool(noOfworkers int) {
	var wg sync.WaitGroup
	for i := range noOfworkers {
		wg.Add(1)
		go worker(i, &wg)
	}
	wg.Wait()
	close(results)
}

func allocateJobs(urls []string) {
	for _, url := range urls {
		job := Job{
			URL: url,
		}
		jobs <- job
	}
	close(jobs)
}

func worker(id int, wg *sync.WaitGroup) {
	fmt.Printf("Worker [%d] looking job\n", id)
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("Worker [%d] started fetchData %s\n", id, job.URL)
		data, err := fetchData(job.URL)
		if err != nil {
			println("[%d]Error fetching data from ", id, job.URL, " : ", err.Error())
			continue
		}
		res := Result{
			URL:  job.URL,
			Data: data,
		}
		results <- res
		fmt.Printf("Worker [%d] finished\n", id)
	}
}

func fetchData(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	//Extracting the title
	title := doc.Find("title").Text()

	//Extracting the first paragraph
	var paragraphs []string
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		paragraphs = append(paragraphs, s.Text())
	})
	//Combining all paragraphs into a single thing
	allParagraphs := strings.Join(paragraphs, "\n\n")

	//format the extracted data
	data := fmt.Sprintf("Title: %s\n\nFirst Paragraph:%s\n", title, allParagraphs)
	return data, nil
}

func resultsContent(done chan bool) {
	for resValue := range results {
		println("Result: ", resValue.URL, " : len\n", len(resValue.Data))
	}
	done <- true
}
