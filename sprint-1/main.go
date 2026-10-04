
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const page = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Sprint 1 | Docker Project</title>
	<style>
		* { box-sizing: border-box; }
		body {
			margin: 0;
			font-family: Arial, sans-serif;
			background: #f1f5f9;
			color: #1e293b;
		}
		header {
			background: #1e40af;
			color: white;
			padding: 20px;
			text-align: center;
		}
		main {
			max-width: 650px;
			margin: 60px auto;
			padding: 20px;
		}
		.card {
			background: white;
			padding: 30px;
			border-radius: 12px;
			text-align: center;
			box-shadow: 0 4px 12px #00000012;
		}
		h1 { color: #1e40af; }
		.info {
			margin-top: 25px;
			padding: 15px;
			background: #eff6ff;
			border-radius: 8px;
			line-height: 1.8;
		}
		.status {
			color: #15803d;
			font-weight: bold;
		}
		footer {
			text-align: center;
			color: #64748b;
			padding: 20px;
		}
	</style>
</head>
<body>
	<header>
		<h2>Sprint 1 - Docker Project</h2>
	</header>
	<main>
		<div class="card">
			<h1>Welcome!</h1>
			<p>My first containerized Go application.</p>
			<div class="info">
				<p><strong>Application:</strong> Go</p>
				<p><strong>Container:</strong> Docker</p>
				<p><strong>Port:</strong> 3000</p>
				<p><strong>Status:</strong>
					<span class="status">Running</span>
				</p>
			</div>
		</div>
	</main>
	<footer>DevOps Internship 2026</footer>
</body>
</html>`

func main() {
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		resp, err := http.Get("http://127.0.0.1:3000/health")
		if err != nil {
			os.Exit(1)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			os.Exit(1)
		}
		os.Exit(0)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, page)
	})


http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.WriteHeader(http.StatusOK)

    fmt.Fprint(w, `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Application Health</title>
    <style>
        * { box-sizing: border-box; }
        body {
            margin: 0;
            font-family: Arial, sans-serif;
            background: #f1f5f9;
            color: #1e293b;
        }
        header {
            background: #1e40af;
            color: white;
            padding: 20px;
            text-align: center;
        }
        main {
            max-width: 600px;
            margin: 60px auto;
            padding: 20px;
        }
        .card {
            background: white;
            padding: 30px;
            border-radius: 12px;
            text-align: center;
            box-shadow: 0 4px 12px #00000012;
        }
        .status {
            display: inline-block;
            padding: 10px 20px;
            margin: 15px 0;
            background: #dcfce7;
            color: #15803d;
            border-radius: 20px;
            font-weight: bold;
        }
        .dot {
            display: inline-block;
            width: 10px;
            height: 10px;
            background: #16a34a;
            border-radius: 50%;
            margin-right: 8px;
        }
        .info {
            margin-top: 20px;
            padding: 15px;
            background: #eff6ff;
            border-radius: 8px;
            text-align: left;
            line-height: 2;
        }
        footer {
            text-align: center;
            color: #64748b;
            padding: 20px;
        }
    </style>
</head>
<body>
    <header>
        <h2>Sprint 1 - Application Health</h2>
    </header>
    <main>
        <div class="card">
            <h1>Health Check</h1>
            <div class="status">
                <span class="dot"></span>Healthy
            </div>
            <p>Your Go application is running successfully.</p>
            <div class="info">
                <strong>Application:</strong> Go<br>
                <strong>Container:</strong> Docker<br>
                <strong>Port:</strong> 3000<br>
                <strong>HTTP Status:</strong> 200 OK
            </div>
        </div>
    </main>
    <footer>DevOps Internship 2026</footer>
</body>
</html>`)
})

http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    w.WriteHeader(http.StatusOK)
    fmt.Fprint(w, "READY")
})

	log.Println("Server running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}