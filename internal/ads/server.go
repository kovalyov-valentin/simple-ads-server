package ads

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/ferluci/fast-realip"
	"github.com/mssola/user_agent"
	"github.com/oschwald/geoip2-golang"
	"github.com/valyala/fasthttp"
)

type Server struct {
	geoip *geoip2.Reader
}

func NewServer(geoip *geoip2.Reader) *Server {
	return &Server{geoip: geoip}
}

func (s *Server) Listen() error {
	return fasthttp.ListenAndServe(":8080", s.handleHttp)
}

func (s *Server) handleHttp(ctx *fasthttp.RequestCtx) {
	ua := string(ctx.Request.Header.UserAgent())
	parsed := user_agent.New(ua)

	ip := realip.FromRequest(ctx)

	browserName, browserVersion := parsed.Browser()
	country, err := s.geoip.Country(net.ParseIP(ip))
	if err != nil {
		log.Println(err)
		return
	}

	u := &User{Country: country.Country.IsoCode, Browser: browserName}
	campaigns := GetCampaigns() 

	winner := MakeAuction(campaigns, u)
	if winner == nil {
		ctx.Redirect("https://example.com", http.StatusSeeOther)
		return 
	}

	ctx.Redirect(winner.ClickUrl, http.StatusSeeOther)

	ctx.WriteString(fmt.Sprintf("User-Agent: %s\n", ua))
	ctx.WriteString(fmt.Sprintf("Browser: %s %s\n", browserName, browserVersion))
	ctx.WriteString(fmt.Sprintf("IP: %s\n", ip))
	ctx.WriteString(fmt.Sprintf("Country: %s\n", country.Country.IsoCode))
}