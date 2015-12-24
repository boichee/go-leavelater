package main

import (
    "log"
    "time"
    // "fmt"
    "os"

    "googlemaps.github.io/maps"
    "github.com/kr/pretty"
    "golang.org/x/net/context"
)

func parseArgs() (string, string) {
    if len(os.Args) != 3 {
        log.Fatalf("fatal error: %s", "You forgot to pass the origin and destination arguments")
    }

    return os.Args[1], os.Args[2]
}

func checkTraffic(c maps.Client,r *maps.DirectionsRequest) bool {
    // Google must have updated their Go Lang package without updating the online documentation.
    // c.Directions now returns a 3rd value, the geocoded waypoint, in 2nd postiion. 
    resp, geopoint, err := c.Directions(context.Background(), r)
    if err != nil {
        log.Fatalf("fatal error: %s", err)
    }

    for idx, val := range geopoint {
        if val.PartialMatch == true {
            log.Fatalf("fatal error: %d %s", idx, "The origin/destination could not be precisely found. Try adding more information.")
        }
    }

    l := resp[0].Legs[0]

    pretty.Println("Normal amount of time:")
    pretty.Println(l.Duration.String())

    pretty.Println("Time in traffic:")
    pretty.Println(l.DurationInTraffic.String())

    return l.Duration.Seconds()*1.1 < l.DurationInTraffic.Seconds()
}

func main() {

    // First parse the arguments, because if there are no arguments, there's really no point in going further


    c, err := maps.NewClient(maps.WithAPIKey("AIzaSyC8tAh9LPl_KfWtuDWNgpiFbwrLnV0Rn6U"))
    if err != nil {
        log.Fatalf("fatal error: %s", err)
    }

    o, d := parseArgs()

    r := &maps.DirectionsRequest{
        Mode:           "driving",
        Origin:         o,
        Destination:    d,
        DepartureTime:  "now",
    }

    


    // Ok so now I'm at the point where I can retrieve a duration for a particular route, and then print out the no traffic estimate, and current estimate with traffic. Now if I can make the application continue to loop until the traffic time matches the normal time, we're there.

    // Just discovered time.Sleep, which pauses the goroutine for a duration of time.

    for checkTraffic(*c, r) {
        pretty.Println("Traffic still sucks. We'll check again in 5.")
        time.Sleep(5 * time.Minute)
    }

    // This method is working, but its hardly the most elegant thing I've ever seen. I think the more appropriate action would be to use channels so I could still allow user input which would allow the user to quit at any time.

    pretty.Println("All is well. You can now leave for your destination.")
}