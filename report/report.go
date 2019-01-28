package report

import (
	"html/template"
	"log"
	"os"
	"talisman/detector"
)

// GenerateReport generates a talisman scan report in html format
func GenerateReport(r *detector.DetectionResults) {
	reportHTML := getReportHTML()
	reportTemplate := template.New("report")
	reportTemplate, _ = reportTemplate.Parse(reportHTML)

	file, err := os.Create("report.html")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	reportTemplate.ExecuteTemplate(file, "report", r)
	file.Close()
}

func getReportHTML() string {
	return `
	<html>
		<head>
			<style>
				body {
					background-image: radial-gradient(at center center, rgb(69, 72, 77) 0%, rgb(17, 17, 17) 100%);
					font-family: "Helvetica Neue",Helvetica,Arial,sans-serif;
				}
				table {
					border-collapse: collapse;
				}
				.commits_table, .details {
					width: 100%;
				}
				.report-table {
					width: 60%;
					margin: auto;
					background-color: #DCDCDC;
				}
				td, th {
					border: 1px solid #dddddd;
					text-align: left;
					padding: 8px;
					text-align: center;
				}
				.report-table th {
					background-color: maroon;
				}
				.report-table > tbody > tr > td {
					border: 1px solid grey;
				}
				.failure-message {
					width: 500px;
					word-break: break-word;
				}
				.heading {
					height: 100px;
					align-content: center;
				}
				#message {
					margin-right: 40%;
				}
				.details > tbody > tr > td {
					border: 2px solid #A9A9A9;
					width: 60%;
				}
				#file-path {
					width: 25%;
				}
				#heading {
					font-size: 40px;
					font-weight: 300;
					color: lightgrey;
					text-align: center;
					margin-top: 3%;
				}
			</style>
			<title>Talisman Report</title>
		</head>
		<body>
			<h1 id="heading">Talisman Scan Report</h1>
			<div>
				<table class="report-table">
					<tr>
						<th id="file-path">File Path</th>
						<th>
							<span id="message">Message</span>
							<span>Commits</span>
						</th>
					</tr>
					{{range $filePath, $FailureData := .Failures}}
						<tr>
							<td>{{$filePath}}</td>
							<td>
								<table class="details">
									{{range $failure := $FailureData}}
										<tr>
											<td class="failure-message">{{$failure.Message}}</td>
											<td>
												<table>
													{{range $commit := $failure.Commits}}
														<tr><td>{{$commit}}</td></tr>
													{{end}}
												</table>
											</td>
										</tr>
									{{end}}
								</table>
							</td>
						</tr>
					{{end}}
				</table>
			</div>
		</body>
	</html>`
}
