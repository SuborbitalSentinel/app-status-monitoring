package main

import (
    "fmt"
    "net/http"
    "strings"
    "time"
)

var websiteTemplate = `
<!doctype html>
<html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-1BmE4kWBq78iYhFldvKuhfTAU6auU8tT94WrHftjDbrCEXSU1oBoqyl2QvZ6jIW3" crossorigin="anonymous">
        <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.7.2/font/bootstrap-icons.css">
        <link rel="stylesheet" href="https://unpkg.com/bootstrap-table@1.20.2/dist/bootstrap-table.min.css">
    </head>
    <body style="background-color:#222224">
        <div class="d-flex flex-column">
            <div class="p-2">
                <form action="/alive" method="POST">
                    <div class="d-flex flex-row">
                        <input name="site" id="site" value="">
                        <button type="submit" class="btn btn-primary">Add Site</button>
                    </div>
                </form>
            </div>
            <div class="p-2">
                <table data-toggle="table" data-sort-name="site" data-sort-order="asc"  data-classes="table table-bordered table-hover table-striped table-dark">
                    <thead>
                        <tr>
                            <th data-field="site">Site Name</th>
                            <th>Status</th>
                            <th>Last Heartbeat</th>
                        </tr>
                    </thead>
                    <tbody>
                        {{site_table}}
                    </tbody>
                </table>
            </div>
        </div>
    </body>
    <script src="https://cdn.jsdelivr.net/npm/jquery/dist/jquery.min.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-ka7Sk0Gln4gmtz2MlQnikT1wXgYsOg+OMhuP+IlRH9sENBO0LRn5q+8nbTov4+1p" crossorigin="anonymous"></script>
    <script src="https://unpkg.com/bootstrap-table@1.20.2/dist/bootstrap-table.min.js"></script>
</html>`

var sites = make(map[string]*SiteStatus)

func buildSiteStatusTable() string {
    const rowTemplate string = "<tr> <td>%v</td> <td>%v</td> <td>%v</td> </tr>"

    builder := strings.Builder{}
    now := time.Now()

    for siteName, ss := range sites {
        statusSymbol := ""
        if ss.wasCheckinMissed(&now) {
            statusSymbol = `<i class="bi bi-x"></i>`
        } else {
            statusSymbol = `<i class="bi bi-check-lg"></i>`
        }
        formattedDate := ss.lastHeartbeat.Format("Jan 01 15:04:05")
        builder.WriteString(fmt.Sprintf(rowTemplate, siteName, statusSymbol, formattedDate))
    }

    return builder.String()
}

type SiteStatus struct {
    lastHeartbeat time.Time
    hasAlerted    bool
}

func (ss SiteStatus) wasCheckinMissed(currentTime *time.Time) bool {
    return ss.lastHeartbeat.Add(1 * time.Minute).Before(*currentTime)
}

func monitor() {
    for {
        fmt.Printf("Number of checked Site(s): %v\n", len(sites))
        now := time.Now()
        for siteName, ss := range sites {
            checkinMissed := ss.wasCheckinMissed(&now)
            if checkinMissed && !ss.hasAlerted {
                fmt.Printf("Site: %s, appears to be down!\n", siteName)
                ss.hasAlerted = true
            } else if !checkinMissed && ss.hasAlerted {
                fmt.Printf("Site: %s, is back online!\n", siteName)
                ss.hasAlerted = false
            }
        }
        time.Sleep(10 * time.Second)
    }
}

func handleHome(w http.ResponseWriter, _ *http.Request) {
    formattedSite := strings.Replace(websiteTemplate, "{{site_table}}", buildSiteStatusTable(), 1)
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    fmt.Fprint(w, formattedSite)
}

func updateHeartbeat(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    siteName := r.FormValue("site")
    if site, ok := sites[siteName]; ok {
        site.lastHeartbeat = time.Now()
    } else {
        sites[siteName] = &SiteStatus{time.Now(), false}
    }

    http.Redirect(w, r, "/", http.StatusMovedPermanently)
    // w.WriteHeader(200)
}

func main() {
    go monitor()

    http.HandleFunc("/", handleHome)
    http.HandleFunc("/alive", updateHeartbeat)
    http.ListenAndServe(":1911", nil)
}
